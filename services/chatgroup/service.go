package chatgroup

import (
	"context"
	"errors"
	"fmt"
	"im/models"
	"im/pkg/database"
	"im/pkg/resp"
	"im/pkg/utils"
	"im/services/chat"
	"im/services/system/config"
	"strconv"
	"strings"
	"time"

	redislib "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	log         *zap.SugaredLogger
	chatSvc     *chat.Service
	mysqlClient *database.Client
	redisClient *redislib.Client
	config      *config.Config
}

func NewService(
	chatSvc *chat.Service,
	mysqlClient *database.Client,
	redisClient *redislib.Client,
	config *config.Config,
) *Service {
	return &Service{
		log:         zap.S().With("module", "services.chat_group.service"),
		chatSvc:     chatSvc,
		mysqlClient: mysqlClient,
		redisClient: redisClient,
		config:      config,
	}
}

// 通过ID查找群聊
func (s *Service) FindChatGroupByID(
	ctx context.Context,
	chatGroupId uint,
) (*models.ChatGroup, error) {
	db := s.mysqlClient.Db()
	chatGroup := &models.ChatGroup{}
	err := db.Model(&models.ChatGroup{}).
		Where("id = ?", chatGroupId).
		First(chatGroup).Error
	if err != nil {
		return nil, err
	}
	if chatGroup.ID == 0 {
		return nil, errors.New(resp.CHAT_GROUP_NOT_EXISTS)
	}

	return chatGroup, nil
}

// 用户是否是群聊成员
func (s *Service) IsChatGroupMember(
	ctx context.Context,
	chatGroupId, accountId uint,
) (bool, bool, bool, error) {
	db := s.mysqlClient.Db()
	isMember, isManager, isOwner := false, false, false
	chatGroupMember := &models.ChatGroupMember{}
	err := db.Model(&models.ChatGroupMember{}).
		Where("chat_group_id = ?", chatGroupId).
		Where("account_id = ?", accountId).
		First(chatGroupMember).Error
	if err != nil {
		return false, false, false, err
	}
	if chatGroupMember.ID != 0 {
		isMember = true
	}
	if chatGroupMember.Role == models.ChatGroupMemberRoleManager || chatGroupMember.Role == models.ChatGroupMemberRoleOwner {
		isManager = true
	}
	if chatGroupMember.Role == models.ChatGroupMemberRoleOwner {
		isOwner = true
	}

	return isMember, isManager, isOwner, nil
}

// 查找加群申请
func (s *Service) FindChatGroupJoin(
	ctx context.Context,
	chatGroupId, accountId uint,
) (*models.ChatGroupJoin, error) {
	db := s.mysqlClient.Db()
	chatGroupApply := &models.ChatGroupJoin{}
	err := db.Model(&models.ChatGroupJoin{}).
		Where("chat_group_id = ?", chatGroupId).
		Where("account_id = ?", accountId).
		Where("status = ?", models.ApplyStatusWait).
		First(chatGroupApply).Error
	if err != nil {
		return nil, err
	}

	return chatGroupApply, nil
}

// 找出群聊所有的管理员ID列表
func (s *Service) FindChatGroupManagerIDList(
	ctx context.Context,
	chatGroupId uint,
) ([]uint, error) {
	db := s.mysqlClient.Db()
	chatGroupMembers := make([]*models.ChatGroupMember, 0)
	err := db.Model(&models.ChatGroupMember{}).
		Where("chat_group_id = ?", chatGroupId).
		Where("role = ? OR role = ?", models.ChatGroupMemberRoleManager, models.ChatGroupMemberRoleOwner).
		Find(&chatGroupMembers).Error
	if err != nil {
		return nil, err
	}

	idList := make([]uint, 0)
	for _, v := range chatGroupMembers {
		idList = append(idList, v.AccountId)
	}
	return idList, nil
}

// 搜索群聊
func (s *Service) SearchChatGroup(
	ctx context.Context,
	accountId uint,
	req *models.SearchChatGroupReq,
) ([]*models.ChatGroup, error) {
	db := s.mysqlClient.Db()
	chatGroups := make([]*models.ChatGroup, 0)
	err := db.Model(&models.ChatGroup{}).
		Where("name like ?", "%"+req.Name+"%").
		Find(&chatGroups).Error
	if err != nil {
		s.log.Errorf("SearchChatGroup select %v", err)
		return nil, err
	}

	for _, v := range chatGroups {
		v.CreatedAt = nil
		v.UpdatedAt = nil
	}

	return chatGroups, nil
}

// 创建群聊
func (s *Service) CreateChatGroup(
	ctx context.Context,
	accountId uint,
	req *models.CreateChatGroupReq,
) (*models.ChatGroup, error) {
	if s.config.GetString("create_chatgroup") != "true" {
		return nil, errors.New(resp.CHAT_GROUP_CREATE_OFF)
	}

	db := s.mysqlClient.Db()

	chatGroupCreateLimit := s.config.GetInt("chatgroup_create_limit")
	chatGroupMembersLimit := s.config.GetInt("chatgroup_members_limit")
	count, err := s.GetChatGorupNumberByOwner(ctx, accountId)
	if err != nil {
		return nil, err
	}
	if count > 0 && count >= int64(chatGroupCreateLimit) {
		return nil, errors.New(resp.CHAT_GROUP_CREATE_LIMIT)
	}

	chatGroup := &models.ChatGroup{
		Name:              req.Name,
		Intro:             req.Intro,
		Avatar:            req.Avatar,
		Members:           1,
		MembersLimit:      uint(chatGroupMembersLimit),
		Latitude:          req.Latitude,
		Longitude:         req.Longitude,
		Country:           req.Country,
		Province:          req.Province,
		City:              req.City,
		District:          req.District,
		Address:           req.Address,
		DisableAddMember:  req.DisableAddMember,
		DisableViewMember: req.DisableViewMember,
		DisbaleAddGroup:   req.DisbaleAddGroup,
		EnbaleBeforeMsg:   req.EnbaleBeforeMsg,
	}
	if chatGroup.Avatar == "" {
		chatGroup.Avatar = s.config.GetString("default_chatgroup_avatar")
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(chatGroup).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.ChatGroupMember{
			ChatGroupId: chatGroup.ID,
			AccountId:   accountId,
			Role:        models.ChatGroupMemberRoleOwner,
		}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		s.log.Errorf("CreateChatGroup Transaction %v", err)
		return nil, err
	}

	geoKey := "chatGroupLocation"
	if req.Latitude != 0 && req.Longitude != 0 {
		go func() {
			err := s.redisClient.GeoAdd(context.Background(), geoKey, &redislib.GeoLocation{
				Name:      fmt.Sprintf("%d", chatGroup.ID),
				Longitude: req.Longitude,
				Latitude:  req.Latitude,
			}).Err()
			if err != nil {
				s.log.Errorf("CreateChatGroup geoadd %v", err)
			}
		}()
	}

	return chatGroup, nil
}

// 修改群资料
func (s *Service) EditChatGroup(
	ctx context.Context,
	accountId uint,
	req *models.EditChatGroupReq,
) error {
	db := s.mysqlClient.Db()
	_, _, isOwner, err := s.IsChatGroupMember(ctx, req.GroupId, accountId)
	if err != nil {
		s.log.Errorf("EditChatGroup IsChatGroupMember %v", err)
		return err
	}
	if !isOwner {
		return errors.New(resp.CHAT_GROUP_NOT_OWNER)
	}
	err = db.Model(&models.ChatGroup{}).
		Where("id = ?", req.GroupId).
		Updates(&models.ChatGroup{
			Name:      req.Name,
			Intro:     req.Intro,
			Avatar:    req.Avatar,
			Longitude: req.Longitude,
			Latitude:  req.Latitude,
			Address:   req.Address,
			Country:   req.Country,
			Province:  req.Province,
			City:      req.City,
			District:  req.District,
		}).Error
	if err != nil {
		s.log.Errorf("EditChatGroup update-1 %v", err)
		return err
	}

	err = db.Model(&models.ChatGroup{}).
		Where("id = ?", req.GroupId).
		Select("disable_add_member", "disable_view_member", "disable_add_group", "enable_before_msg").
		Updates(&models.ChatGroup{
			DisableAddMember:  req.DisableAddMember,
			DisableViewMember: req.DisableViewMember,
			DisbaleAddGroup:   req.DisbaleAddGroup,
			EnbaleBeforeMsg:   req.EnbaleBeforeMsg,
		}).Error
	if err != nil {
		s.log.Errorf("EditChatGroup update-2 %v", err)
		return err
	}

	geoKey := "chatGroupLocation"
	if req.Latitude != 0 && req.Longitude != 0 {
		go func() {
			err := s.redisClient.GeoAdd(context.Background(), geoKey, &redislib.GeoLocation{
				Name:      fmt.Sprintf("%d", req.GroupId),
				Longitude: req.Longitude,
				Latitude:  req.Latitude,
			}).Err()
			if err != nil {
				s.log.Errorf("CreateChatGroup geoadd %v", err)
			}
		}()
	}

	return nil
}

// 加群申请
func (s *Service) JoinChatGroup(
	ctx context.Context,
	accountId uint,
	req *models.JoinChatGroupReq,
) error {
	db := s.mysqlClient.Db()
	chatGroup, err := s.FindChatGroupByID(ctx, req.GroupId)
	if err != nil {
		return err
	}
	if chatGroup.DisbaleAddGroup {
		return errors.New(resp.CHAT_GROUP_DISABLE_ADD_GROUP)
	}
	if chatGroup.Members >= chatGroup.MembersLimit {
		return errors.New(resp.CHAT_GROUP_MEMBERS_LIMIT)
	}

	isMember, _, _, err := s.IsChatGroupMember(ctx, req.GroupId, accountId)
	if err != nil {
		return err
	}
	if isMember {
		return errors.New(resp.CHAT_GROUP_IS_MEMBER)
	}

	chatGroupApply, err := s.FindChatGroupJoin(ctx, req.GroupId, accountId)
	if err != nil {
		return err
	}
	if chatGroupApply.ID != 0 {
		return errors.New(resp.CHAT_GROUP_APPLY_EXISTS)
	}

	managerAccountList := ""
	idList, err := s.FindChatGroupManagerIDList(ctx, req.GroupId)
	if err != nil {
		s.log.Errorf("JoinChatGroup FindChatGroupManagerIDList %v", err)
		return err
	}
	for _, v := range idList {
		managerAccountList += "," + strconv.Itoa(int(v))
	}
	managerAccountList = strings.Trim(managerAccountList, ",")
	err = db.Create(&models.ChatGroupJoin{
		AccountId:          accountId,
		ChatGroupId:        req.GroupId,
		ManagerAccountList: managerAccountList,
		ApplyReason:        req.Reason,
	}).Error
	if err != nil {
		s.log.Errorf("JoinChatGroup Create %v", err)
		return err
	}

	// 实时通知管理员有新的加群请求
	go func() {
		message := &models.Message{
			FromId: req.GroupId,
			Ope:    models.MessageOpeSystem,
			Type:   models.MessageTypeJoinGroup,
		}
		for _, v := range idList {
			message.ToId = v
			err := s.chatSvc.RPC.SendMessageCall(context.Background(), message)
			if err != nil {
				s.log.Errorf("JoinChatGroup send %v", err)
			}
		}
	}()

	return nil
}

// 加群审批
func (s *Service) JoinChatGroupReply(
	ctx context.Context,
	accountId uint,
	req *models.JoinChatGroupReplyReq,
) error {
	db := s.mysqlClient.Db()

	chatGroupApply, err := s.FindChatGroupJoin(ctx, req.GroupId, req.AccountId)
	if err != nil {
		s.log.Errorf("JoinChatGroupReply FindChatGroupJoin %v", err)
		return err
	}
	if chatGroupApply.ID == 0 {
		return errors.New(resp.CHAT_GROUP_APPLY_NOT_EXISTS)
	}
	if chatGroupApply.Status != models.ApplyStatusWait {
		return errors.New(resp.CHAT_GROUP_APPLY_NOT_WAIT)
	}

	_, isManager, _, err := s.IsChatGroupMember(ctx, req.GroupId, accountId)
	if err != nil {
		s.log.Errorf("JoinChatGroupReply IsChatGroupMember %v", err)
		return err
	}
	if !isManager {
		return errors.New(resp.CHAT_GROUP_NOT_MANAGER)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.ChatGroupJoin{}).
			Where("id = ?", 1).
			Updates(&models.ChatGroupJoin{
				Status:     req.Status,
				DenyReason: req.Reason,
			}).Error; err != nil {
			return err
		}

		if req.Status == models.ApplyStatusPass {
			if err := tx.Create(&models.ChatGroupMember{
				ChatGroupId: req.GroupId,
				AccountId:   req.AccountId,
				Role:        models.ChatGroupMemberRoleGeneral,
			}).Error; err != nil {
				return err
			}

			if err := tx.Model(&models.ChatGroup{}).
				Where("id = ?", req.GroupId).
				UpdateColumn("members", gorm.Expr("members + 1")).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		s.log.Errorf("JoinChatGroupReply Transaction %v", err)
		return err
	}

	return nil
}

// 加群申请列表
func (s *Service) JoinChatGroupList(
	ctx context.Context,
	accountId uint,
) ([]*models.ChatGroupJoin, error) {
	db := s.mysqlClient.Db()

	chatGroupJoins := make([]*models.ChatGroupJoin, 0)
	err := db.Model(&models.ChatGroupJoin{}).
		Where("FIND_IN_SET(?, manager_account_list)", accountId).
		Where("status = ?", models.ApplyStatusWait).
		Preload("Account").
		Preload("ChatGroup").
		Find(&chatGroupJoins).Error
	if err != nil {
		s.log.Errorf("JoinChatGroupList select %v", err)
		return nil, err
	}

	for _, v := range chatGroupJoins {
		v.ID = 0
		v.UpdatedAt = nil
		if v.Account != nil {
			v.Account.CreatedAt = nil
			v.Account.UpdatedAt = nil
		}
		if v.ChatGroup != nil {
			v.ChatGroup.CreatedAt = nil
			v.ChatGroup.UpdatedAt = nil
		}
	}

	return chatGroupJoins, nil
}

// 我的群聊列表
func (s *Service) ChatGroupList(
	ctx context.Context,
	accountId uint,
) ([]*models.ChatGroup, error) {
	db := s.mysqlClient.Db()
	chatGroupMembers := make([]*models.ChatGroupMember, 0)
	err := db.Model(&models.ChatGroupMember{}).
		Where("account_id = ?", accountId).
		Find(&chatGroupMembers).Error
	if err != nil {
		s.log.Errorf("ChatGroupList select chatGroupMembers %v", err)
		return nil, err
	}

	chatGroupIdList := make([]uint, 0)
	for _, v := range chatGroupMembers {
		chatGroupIdList = append(chatGroupIdList, v.ChatGroupId)
	}
	chatGroups := make([]*models.ChatGroup, 0)
	err = db.Model(&models.ChatGroup{}).
		Where("id in ?", chatGroupIdList).
		Order("name asc").
		Find(&chatGroups).Error
	if err != nil {
		s.log.Errorf("ChatGroupList select chatGroups %v", err)
		return nil, err
	}

	for _, v := range chatGroups {
		v.CreatedAt = nil
		v.UpdatedAt = nil
	}

	return chatGroups, nil
}

// 群聊信息(包括成员列表)
func (s *Service) ChatGroupInfo(
	ctx context.Context,
	accountId uint,
	req *models.GroupIdReq,
) (*models.ChatGroup, error) {
	db := s.mysqlClient.Db()
	chatGroup, err := s.FindChatGroupByID(ctx, req.GroupId)
	if err != nil {
		return nil, err
	}

	members := make([]*models.ChatGroupMember, 0)
	err = db.Model(&models.ChatGroupMember{}).
		Where("chat_group_id = ?", req.GroupId).
		Order("FIELD(role,'owner','manager','general')").
		Preload("Account").
		Find(&members).Error
	if err != nil {
		s.log.Errorf("ChatGroupInfo select1 %v", err)
		return nil, err
	}

	self := &models.ChatGroupMember{}
	err = db.Model(&models.ChatGroupMember{}).
		Where("chat_group_id = ?", req.GroupId).
		Where("account_id = ?", accountId).
		First(self).Error
	if err != nil {
		s.log.Errorf("ChatGroupInfo select2 %v", err)
		return nil, err
	}

	for _, v := range members {
		if v.BannedOverTime != nil && v.BannedOverTime.Unix() > time.Now().Unix() {
			v.IsBanned = true
		}
	}
	chatGroup.MembersList = members

	chatGroup.SelfInfo = &models.ChatGroupSelfInfo{
		Role: self.Role,
	}
	if self.BannedOverTime != nil && self.BannedOverTime.Unix() > time.Now().Unix() {
		chatGroup.SelfInfo.IsBanned = true
	}

	chatGroup.CreatedAt = nil
	chatGroup.UpdatedAt = nil
	for _, v := range chatGroup.MembersList {
		v.ID = 0
		v.CreatedAt = nil
		v.UpdatedAt = nil
		if v.Account != nil {
			v.Account.CreatedAt = nil
			v.Account.UpdatedAt = nil
		}
	}

	return chatGroup, err
}

// 退出群聊
func (s *Service) ExitChatGroup(
	ctx context.Context,
	accountId uint,
	req *models.GroupIdReq,
) error {
	db := s.mysqlClient.Db()
	isMemeber, _, isOwner, err := s.IsChatGroupMember(ctx, req.GroupId, accountId)
	if err != nil {
		s.log.Errorf("ExitChatGroup IsChatGroupMember %v", err)
		return err
	}
	if !isMemeber {
		return errors.New(resp.CHAT_GROUP_NOT_MEMBER)
	}
	if isOwner {
		return errors.New(resp.CHAT_GROUP_OWNER_EXIT)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("chat_group_id = ?", req.GroupId).
			Where("account_id = ?", accountId).
			Delete(&models.ChatGroupMember{}).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.ChatGroup{}).
			Where("id = ?", req.GroupId).
			UpdateColumn("members", gorm.Expr("members - 1")).
			Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		s.log.Errorf("ExitChatGroup Transaction %v", err)
		return err
	}

	return nil
}

// 转让群聊
func (s *Service) TransferChatGroup(
	ctx context.Context,
	accountId uint,
	req *models.ChatGroupToIDReq,
) error {
	db := s.mysqlClient.Db()
	_, _, isOwner, err := s.IsChatGroupMember(ctx, req.GroupId, accountId)
	if err != nil {
		s.log.Errorf("TransferChatGroup IsChatGroupMember %v", err)
		return err
	}
	if !isOwner {
		return errors.New(resp.CHAT_GROUP_NOT_OWNER)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.ChatGroupMember{}).
			Where("chat_group_id = ?", req.GroupId).
			Where("account_id = ?", req.ToID).
			UpdateColumn("role", models.ChatGroupMemberRoleOwner).Error; err != nil {
			return err
		}

		if err := tx.Where("chat_group_id = ?", req.GroupId).
			Where("account_id = ?", accountId).
			Delete(&models.ChatGroupMember{}).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.ChatGroup{}).
			Where("id = ?", req.GroupId).
			UpdateColumn("members", gorm.Expr("members - 1")).
			Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		s.log.Errorf("TransferChatGroup Transaction %v", err)
		return err
	}

	return nil
}

// 解散群聊
func (s *Service) DissolveChatGroup(
	ctx context.Context,
	accountId uint,
	req *models.GroupIdReq,
) error {
	db := s.mysqlClient.Db()
	_, _, isOwner, err := s.IsChatGroupMember(ctx, req.GroupId, accountId)
	if err != nil {
		s.log.Errorf("DissolveChatGroup IsChatGroupMember %v", err)
		return err
	}
	if !isOwner {
		return errors.New(resp.CHAT_GROUP_NOT_OWNER)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", req.GroupId).
			Delete(&models.ChatGroup{}).Error; err != nil {
			return err
		}

		if err := tx.Where("chat_group_id = ?", req.GroupId).
			Delete(&models.ChatGroupMember{}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		s.log.Errorf("DissolveChatGroup Transaction %v", err)
		return err
	}

	return nil
}

// 群聊踢出成员
func (s *Service) ChatGroupKickMember(
	ctx context.Context,
	accountId uint,
	req *models.ChatGroupToIDReq,
) error {
	db := s.mysqlClient.Db()
	_, isManager, isOwner, err := s.IsChatGroupMember(ctx, req.GroupId, accountId)
	if err != nil {
		s.log.Errorf("ChatGroupKickMember IsChatGroupMember1 %v", err)
		return err
	}
	if !isOwner && !isManager {
		return errors.New(resp.CHAT_GROUP_NOT_MANAGER)
	}

	isMember, isManager, isOwner, err := s.IsChatGroupMember(ctx, req.GroupId, req.ToID)
	if err != nil {
		s.log.Errorf("ChatGroupKickMember IsChatGroupMember2 %v", err)
		return err
	}
	if isManager || isOwner {
		return errors.New(resp.CHAT_GROUP_NOT_ALLOW)
	}
	if !isMember {
		return errors.New(resp.CHAT_GROUP_MEMBER_NOT_EXISTS)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("chat_group_id = ?", req.GroupId).
			Where("account_id = ?", req.ToID).
			Delete(&models.ChatGroupMember{}).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.ChatGroup{}).
			Where("id = ?", req.GroupId).
			UpdateColumn("members", gorm.Expr("members - 1")).
			Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		s.log.Errorf("ChatGroupKickMember Transaction %v", err)
		return err
	}

	return nil
}

// 设置群聊管理员
func (s *Service) ChatGroupSetManager(
	ctx context.Context,
	accountId uint,
	req *models.ChatGroupSetManagerReq,
) error {
	db := s.mysqlClient.Db()
	_, _, isOwner, err := s.IsChatGroupMember(ctx, req.GroupId, accountId)
	if err != nil {
		s.log.Errorf("ChatGroupSetManager IsChatGroupMember1 %v", err)
		return err
	}
	if !isOwner {
		return errors.New(resp.CHAT_GROUP_NOT_OWNER)
	}

	isMember, _, _, err := s.IsChatGroupMember(ctx, req.GroupId, req.ToID)
	if err != nil {
		s.log.Errorf("ChatGroupSetManager IsChatGroupMember2 %v", err)
		return err
	}
	if !isMember {
		return errors.New(resp.CHAT_GROUP_MEMBER_NOT_EXISTS)
	}

	role := models.ChatGroupMemberRoleGeneral
	if req.IsManager {
		role = models.ChatGroupMemberRoleManager
	}
	err = db.Model(&models.ChatGroupMember{}).
		Where("chat_group_id = ?", req.GroupId).
		Where("account_id = ?", req.ToID).
		UpdateColumn("role", role).Error
	if err != nil {
		s.log.Errorf("ChatGroupSetManager update %v", err)
		return err
	}

	return nil
}

// 群聊成员禁言
func (s *Service) ChatGroupBannedMember(
	ctx context.Context,
	accountId uint,
	req *models.ChatGroupBannedMemberReq,
) error {
	db := s.mysqlClient.Db()
	{
		_, isManager, isOwner, err := s.IsChatGroupMember(ctx, req.GroupId, accountId)
		if err != nil {
			s.log.Errorf("ChatGroupBannedMember IsChatGroupMember1 %v", err)
			return err
		}
		if !isOwner && !isManager {
			return errors.New(resp.CHAT_GROUP_NOT_MANAGER)
		}
	}
	{
		isMember, isManager, isOwner, err := s.IsChatGroupMember(ctx, req.GroupId, req.ToID)
		if err != nil {
			s.log.Errorf("ChatGroupBannedMember IsChatGroupMember2 %v", err)
			return err
		}
		if !isMember {
			return errors.New(resp.CHAT_GROUP_MEMBER_NOT_EXISTS)
		}
		if isOwner && isManager {
			return errors.New(resp.CHAT_GROUP_NOT_ALLOW)
		}
	}
	bannedOverTime := time.Now().Add(time.Minute * time.Duration(req.Minute))
	err := db.Model(&models.ChatGroupMember{}).
		Where("chat_group_id = ?", req.GroupId).
		Where("account_id = ?", req.ToID).
		UpdateColumn("banned_over_time", bannedOverTime).Error
	if err != nil {
		s.log.Errorf("ChatGroupBannedMember update %v", err)
		return err
	}

	return nil
}

// 获取好友共同群组
func (s *Service) GetFriendCommonChatGroups(
	ctx context.Context,
	accountId uint,
	req *models.ToIDReq,
) ([]*models.ChatGroup, error) {
	db := s.mysqlClient.Db()

	// 好友加的群
	friendChatGroupMembers := make([]*models.ChatGroupMember, 0)
	err := db.Model(&models.ChatGroupMember{}).
		Where("account_id = ?", req.ToID).
		Find(&friendChatGroupMembers).Error
	if err != nil {
		s.log.Errorf("GetFriendCommonChatGroups select friend %v", err)
		return nil, err
	}
	friendGropIdList := make([]uint, 0)
	for _, v := range friendChatGroupMembers {
		friendGropIdList = append(friendGropIdList, v.ChatGroupId)
	}

	// 我加的群
	selfChatGroupMembers := make([]*models.ChatGroupMember, 0)
	err = db.Model(&models.ChatGroupMember{}).
		Where("account_id = ?", accountId).
		Find(&selfChatGroupMembers).Error
	if err != nil {
		s.log.Errorf("GetFriendCommonChatGroups select self %v", err)
		return nil, err
	}
	selfGropIdList := make([]uint, 0)
	for _, v := range selfChatGroupMembers {
		selfGropIdList = append(selfGropIdList, v.ChatGroupId)
	}

	// 查询公共群聊
	commonChatGroupIdList := utils.IntersectUint(friendGropIdList, selfGropIdList)
	chatGroups := make([]*models.ChatGroup, 0)
	err = db.Model(&models.ChatGroup{}).
		Where("id in ?", commonChatGroupIdList).
		Find(&chatGroups).Error
	if err != nil {
		s.log.Errorf("GetFriendCommonChatGroups select common %v", err)
		return nil, err
	}
	for _, v := range chatGroups {
		v.CreatedAt = nil
		v.UpdatedAt = nil
	}

	return chatGroups, nil
}

// 附近的群
func (s *Service) NearChatGroups(
	ctx context.Context,
	accountId uint,
	req *models.NearChatGroupsReq,
) ([]*models.NearChatGroupsItem, error) {
	db := s.mysqlClient.Db()
	geoKey := "chatGroupLocation"

	largeDistance := s.config.GetFloat64("near_chatgroup_distance")
	if largeDistance == 0 {
		largeDistance = 1000
	}

	res, err := s.redisClient.GeoRadius(ctx, geoKey, req.Longitude, req.Latitude, &redislib.GeoRadiusQuery{
		Radius:    largeDistance,
		Unit:      "km",
		WithCoord: true,
		WithDist:  true,
		Count:     100,
		Sort:      "ASC",
	}).Result()
	if err != nil {
		s.log.Errorf("NearChatGroups cmd.GeoRadius %v", err)
		return nil, err
	}

	items := make([]*models.NearChatGroupsItem, 0)
	for _, v := range res {
		chatGroup := &models.ChatGroup{}
		err := db.Model(&models.ChatGroup{}).
			Where("id = ?", v.Name).
			First(chatGroup).Error
		if err != nil {
			s.log.Errorf("NearChatGroups chatGroup select %v", err)
			continue
		}
		chatGroup.Longitude = v.Longitude
		chatGroup.Latitude = v.Latitude
		items = append(items, &models.NearChatGroupsItem{
			ChatGroup: chatGroup,
			Distance:  fmt.Sprintf("%.2fkm", v.Dist),
		})
	}

	for _, v := range items {
		if v.ChatGroup != nil {
			v.ChatGroup.CreatedAt = nil
			v.ChatGroup.UpdatedAt = nil
		}
	}

	return items, nil
}

// 群聊邀请成员
func (s *Service) ChatGroupInviteMember(
	ctx context.Context,
	accountId uint,
	req *models.ChatGroupToIDListReq,
) error {
	db := s.mysqlClient.Db()
	_, isManager, isOwner, err := s.IsChatGroupMember(ctx, req.GroupId, accountId)
	if err != nil {
		s.log.Errorf("ChatGroupInviteMembers IsChatGroupMember %v", err)
		return err
	}
	if !isOwner && !isManager {
		return errors.New(resp.CHAT_GROUP_NOT_MANAGER)
	}

	invites := make([]*models.ChatGroupMember, 0)
	for _, toID := range req.ToIDList {
		isMember, _, _, err := s.IsChatGroupMember(ctx, req.GroupId, toID)
		if err != nil {
			s.log.Errorf("ChatGroupInviteMember IsChatGroupMember %v", err)
			return err
		}
		if !isMember {
			invites = append(invites, &models.ChatGroupMember{
				ChatGroupId: req.GroupId,
				AccountId:   toID,
				Role:        models.ChatGroupMemberRoleGeneral,
			})
		}
	}
	if len(invites) > 0 {
		err = db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&models.ChatGroupMember{}).
				Create(invites).Error; err != nil {
				return err
			}

			if err := tx.Model(&models.ChatGroup{}).
				Where("id = ?", req.GroupId).
				UpdateColumn("members", gorm.Expr("members + ?", len(invites))).
				Error; err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			s.log.Errorf("ChatGroupInviteMember Transaction %v", err)
			return err
		}
	}

	return nil
}

// 获取群主创建的群聊数量
func (s *Service) GetChatGorupNumberByOwner(
	ctx context.Context,
	accountId uint,
) (int64, error) {
	db := s.mysqlClient.Db()
	var count int64 = 0
	err := db.Model(&models.ChatGroupMember{}).
		Where("account_id = ?", accountId).
		Where("role = ?", models.ChatGroupMemberRoleOwner).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
