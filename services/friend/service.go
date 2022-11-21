package friend

import (
	"context"
	"errors"
	"im/models"
	"im/pkg/database"
	"im/pkg/resp"
	"im/services/system/config"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	log         *zap.SugaredLogger
	mysqlClient *database.Client
	config      *config.Config
}

func NewService(
	mysqlClient *database.Client,
	config *config.Config,
) *Service {
	return &Service{
		log:         zap.S().With("module", "services.friend.service"),
		mysqlClient: mysqlClient,
		config:      config,
	}
}

// 查找好友
func (s *Service) SearchFriend(
	ctx context.Context,
	accountId uint,
	req *models.SearchFriendReq,
) ([]*models.Account, error) {
	accounts := make([]*models.Account, 0)
	err := s.mysqlClient.Db().
		Model(&models.Account{}).
		Where("`id` <> ? AND (`username` = ? OR `nickname` = ? OR `mobile` = ? OR `email` = ? OR `id` = ?)",
			accountId, req.Account, req.Account, req.Account, req.Account, req.Account).
		Find(&accounts).
		Error
	if err != nil {
		return nil, err
	}

	for _, v := range accounts {
		v.CreatedAt = nil
		v.UpdatedAt = nil
	}

	return accounts, nil
}

// 添加好友
func (s *Service) AddFriend(
	ctx context.Context,
	accountId uint,
	req *models.AddFriendReq,
) error {
	db := s.mysqlClient.Db()
	// 查询是否已经是好友
	friend := &models.Friend{}
	err := db.Model(&models.Friend{}).
		Where("account_id = ?", accountId).
		Where("friend_id = ?", req.ToID).
		First(friend).Error
	if err != nil {
		return err
	}
	if friend.ID != 0 {
		return errors.New(resp.FRIEND_EXISTS)
	}
	// 查询是否请求过还没验证
	friendApply := &models.FriendApply{}
	err = db.Model(&models.FriendApply{}).
		Where("from_account_id = ?", accountId).
		Where("to_account_id = ?", req.ToID).
		Where("status = ?", models.FriendApplyStatusWait).
		First(friendApply).Error
	if err != nil {
		return err
	}
	if friendApply.ID != 0 {
		return errors.New(resp.FRIEND_APPLY_EXISTS)
	}
	// 创建好友请求
	err = db.Model(&models.FriendApply{}).
		Create(&models.FriendApply{
			FromAccountId: accountId,
			ToAccountId:   req.ToID,
			ApplyReason:   req.Reason,
			Status:        models.FriendApplyStatusWait,
		}).Error
	if err != nil {
		return err
	}

	return nil
}

// 同意/拒绝好友请求
func (s *Service) AddFriendReply(
	ctx context.Context,
	accountId uint,
	req *models.AddFriendReplyReq,
) error {
	db := s.mysqlClient.Db()
	// 查询是否已经是好友
	friend := &models.Friend{}
	err := db.Model(&models.Friend{}).
		Where("account_id = ?", accountId).
		Where("friend_id = ?", req.ToID).
		First(friend).Error
	if err != nil {
		return err
	}
	if friend.ID != 0 {
		return errors.New(resp.FRIEND_EXISTS)
	}
	// 查询请求是否存在
	friendApply := &models.FriendApply{}
	err = db.Model(&models.FriendApply{}).
		Where("from_account_id = ?", req.ToID).
		Where("to_account_id = ?", accountId).
		Where("status = ?", models.FriendApplyStatusWait).
		First(friendApply).Error
	if err != nil {
		return err
	}
	if friendApply.ID == 0 {
		return errors.New(resp.FRIEND_APPLY_NOT_EXISTS)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.FriendApply{}).
			Where("id = ?", friendApply.ID).
			Select("status", "deny_reason").
			Updates(&models.FriendApply{
				Status:     req.Status,
				DenyReason: req.Reason,
			}).Error; err != nil {
			return err
		}

		if req.Status == models.FriendApplyStatusPass {
			if err := db.Model(&models.Friend{}).
				Create(&models.Friend{
					AccountId: req.ToID,
					FriendId:  accountId,
				}).Error; err != nil {
				return err
			}
			if err := db.Model(&models.Friend{}).
				Create(&models.Friend{
					AccountId: accountId,
					FriendId:  req.ToID,
				}).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// 好友申请列表
func (s *Service) FriendApplyList(
	ctx context.Context,
	accountId uint,
) ([]*models.FriendApply, error) {
	friendApplies := make([]*models.FriendApply, 0)
	err := s.mysqlClient.Db().
		Model(&models.FriendApply{}).
		Where("to_account_id = ?", accountId).
		Where("status = ?", models.FriendApplyStatusWait).
		Order("id desc").
		Preload("FromAccount").
		Find(&friendApplies).Error
	if err != nil {
		return nil, err
	}

	for _, v := range friendApplies {
		v.ID = 0
		v.UpdatedAt = nil
		if v.FromAccount != nil {
			v.FromAccount.CreatedAt = nil
			v.FromAccount.UpdatedAt = nil
		}
	}

	return friendApplies, nil
}

// 好友列表
func (s *Service) FriendList(
	ctx context.Context,
	accountId uint,
	req *models.FriendListReq,
) ([]*models.Friend, error) {
	friends := make([]*models.Friend, 0)
	err := s.mysqlClient.Db().
		Model(&models.Friend{}).
		Where("account_id = ?", accountId).
		Where("blacklist = ?", req.BlackList).
		Preload("Account").
		Preload("FriendsGroups").
		Find(&friends).Error
	if err != nil {
		return nil, err
	}

	for _, v := range friends {
		v.ID = 0
		v.CreatedAt = nil
		v.UpdatedAt = nil
		if v.Account != nil {
			v.Account.CreatedAt = nil
			v.Account.UpdatedAt = nil
		}
		for _, v2 := range v.FriendsGroups {
			v2.CreatedAt = nil
			v2.UpdatedAt = nil
		}
	}

	return friends, nil
}

// 单个好友信息
func (s *Service) FriendInfo(
	ctx context.Context,
	accountId uint,
	req *models.ToIDReq,
) (*models.Friend, error) {
	friend := &models.Friend{}
	err := s.mysqlClient.Db().
		Model(&models.Friend{}).
		Where("account_id = ?", accountId).
		Where("friend_id = ?", req.ToID).
		Preload("Account").
		Preload("FriendsGroups").
		First(friend).Error
	if err != nil {
		return nil, err
	}
	if friend.ID == 0 {
		return nil, errors.New(resp.FRIEND_NOT_EXISTS)
	}

	friend.ID = 0
	friend.CreatedAt = nil
	friend.UpdatedAt = nil
	if friend.Account != nil {
		friend.Account.CreatedAt = nil
		friend.Account.UpdatedAt = nil
	}
	for _, v2 := range friend.FriendsGroups {
		v2.CreatedAt = nil
		v2.UpdatedAt = nil
	}

	return friend, nil
}

// 删除好友
func (s *Service) DeleteFriend(
	ctx context.Context,
	accountId uint,
	req *models.ToIDReq,
) error {
	db := s.mysqlClient.Db()

	friend := &models.Friend{}
	err := db.Model(&models.Friend{}).
		Where("account_id = ?", accountId).
		Where("friend_id = ?", req.ToID).
		First(friend).Error
	if err != nil {
		s.log.Errorf("DeleteFriend select %v", err)
		return err
	}
	if friend.ID == 0 {
		return errors.New(resp.FRIEND_NOT_EXISTS)
	}

	err = db.Delete(friend).Error
	if err != nil {
		s.log.Errorf("DeleteFriend delete %v", err)
		return err
	}

	return nil
}

// 添加黑名单
func (s *Service) AddBlacklist(
	ctx context.Context,
	accountId uint,
	req *models.ToIDReq,
) error {
	db := s.mysqlClient.Db()

	friend := &models.Friend{}
	err := db.Model(&models.Friend{}).
		Where("account_id = ?", accountId).
		Where("friend_id = ?", req.ToID).
		First(friend).Error
	if err != nil {
		s.log.Errorf("AddBlacklist select %v", err)
		return err
	}

	if friend.ID == 0 {
		err = db.Create(&models.Friend{
			AccountId: accountId,
			FriendId:  req.ToID,
			Blacklist: true,
		}).Error
	} else {
		err = db.Model(&models.Friend{}).
			Where("id = ?", friend.ID).
			Select("blacklist").
			Updates(&models.Friend{
				Blacklist: true,
			}).Error
	}
	if err != nil {
		s.log.Errorf("AddBlacklist create/update %v", err)
		return err
	}

	return nil
}

// 移除黑名单
func (s *Service) DeleteBlacklist(
	ctx context.Context,
	accountId uint,
	req *models.ToIDReq,
) error {
	db := s.mysqlClient.Db()

	friend := &models.Friend{}
	err := db.Model(&models.Friend{}).
		Where("account_id = ?", accountId).
		Where("friend_id = ?", req.ToID).
		Where("blacklist = ?", true).
		First(friend).Error
	if err != nil {
		s.log.Errorf("DeleteBlacklist select %v", err)
		return err
	}
	if friend.ID == 0 {
		return errors.New(resp.BLACKLIST_NOT_EXISTS)
	}

	err = db.Delete(friend).Error
	if err != nil {
		s.log.Errorf("DeleteBlacklist delete %v", err)
		return err
	}

	return nil
}

// 设置好友备注
func (s *Service) SetFriendRemark(
	ctx context.Context,
	accountId uint,
	req *models.SetFriendRemarkReq,
) error {
	db := s.mysqlClient.Db()

	friend := &models.Friend{}
	err := db.Model(&models.Friend{}).
		Where("account_id = ?", accountId).
		Where("friend_id = ?", req.ToID).
		First(friend).Error
	if err != nil {
		s.log.Errorf("FriendRemark select %v", err)
		return err
	}
	if friend.ID == 0 {
		return errors.New(resp.FRIEND_NOT_EXISTS)
	}

	err = db.Model(&models.Friend{}).
		Where("id = ?", friend.ID).
		UpdateColumn("remark", req.Remark).
		Error
	if err != nil {
		s.log.Errorf("FriendRemark update %v", err)
		return err
	}

	return nil
}

// 设置好友标签（自定义字段）
func (s *Service) SetFriendLabel(
	ctx context.Context,
	accountId uint,
	req *models.SetFriendLabelReq,
) error {
	db := s.mysqlClient.Db()

	friend := &models.Friend{}
	err := db.Model(&models.Friend{}).
		Where("account_id = ?", accountId).
		Where("friend_id = ?", req.ToID).
		First(friend).Error
	if err != nil {
		s.log.Errorf("FriendLabel select %v", err)
		return err
	}
	if friend.ID == 0 {
		return errors.New(resp.FRIEND_NOT_EXISTS)
	}

	err = db.Model(&models.Friend{}).
		Where("id = ?", friend.ID).
		UpdateColumn("label", req.Label).
		Error
	if err != nil {
		s.log.Errorf("FriendLabel update %v", err)
		return err
	}

	return nil
}

// 创建好友分组
func (s *Service) CreateFriendGroup(
	ctx context.Context,
	accountId uint,
	req *models.CreateFriendGroupReq,
) error {
	db := s.mysqlClient.Db()

	tempFG := &models.FriendGroup{}
	err := db.Model(&models.FriendGroup{}).
		Where("account_id = ?", accountId).
		Where("name = ?", req.Name).
		First(tempFG).Error
	if err != nil {
		s.log.Errorf("CreateFriendGroup select %v", err)
		return err
	}
	if tempFG.ID != 0 {
		return errors.New(resp.FRIEND_GROUP_EXISTS)
	}

	friendGroup := &models.FriendGroup{
		AccountId: accountId,
		Name:      req.Name,
	}
	err = db.Create(friendGroup).Error
	if err != nil {
		s.log.Errorf("CreateFriendGroup create %v", err)
		return err
	}

	if len(req.Members) > 0 {
		err = s.OperateFriendGroupMembers(ctx, accountId, &models.FriendGroupMembersReq{
			GroupId: friendGroup.ID,
			Members: req.Members,
		}, "add")
		if err != nil {
			return err
		}
	}

	return nil
}

// 删除好友分组
func (s *Service) DeleteFriendGroup(
	ctx context.Context,
	accountId uint,
	req *models.FriendGroupReq,
) error {
	db := s.mysqlClient.Db()

	friendGroup := &models.FriendGroup{}
	err := db.Model(&models.FriendGroup{}).
		Where("id = ?", req.GroupId).
		Where("account_id = ?", accountId).
		First(friendGroup).Error
	if err != nil {
		s.log.Errorf("DeleteFriendGroup select %v", err)
		return err
	}
	if friendGroup.ID == 0 {
		return errors.New(resp.FRIEND_GROUP_NOT_EXISTS)
	}

	err = db.Select("Friends").
		Delete(friendGroup).Error
	if err != nil {
		s.log.Errorf("DeleteFriendGroup delete %v", err)
		return err
	}

	return nil
}

// 好友分组添加/删除成员
func (s *Service) OperateFriendGroupMembers(
	ctx context.Context,
	accountId uint,
	req *models.FriendGroupMembersReq,
	operate string,
) error {
	db := s.mysqlClient.Db()

	friendGroup := &models.FriendGroup{}
	err := db.Model(&models.FriendGroup{}).
		Where("id = ?", req.GroupId).
		Where("account_id = ?", accountId).
		First(friendGroup).Error
	if err != nil {
		s.log.Errorf("AddFriendGroupMembers select %v", err)
		return err
	}
	if friendGroup.ID == 0 {
		return errors.New(resp.FRIEND_GROUP_NOT_EXISTS)
	}

	friends := make([]*models.Friend, 0)
	for _, v := range req.Members {
		friends = append(friends, &models.Friend{FriendId: v})
	}
	switch operate {
	case "add":
		err = db.Model(friendGroup).
			Omit("Friends.*").
			Association("Friends").
			Append(friends)
	case "del":
		err = db.Model(friendGroup).
			Omit("Friends.*").
			Association("Friends").
			Delete(friends)
	}
	if err != nil {
		s.log.Errorf("AddFriendGroupMembers operate %v", err)
		return err
	}

	return nil
}

// 获取好友分组
func (s *Service) GetFriendGroups(
	ctx context.Context,
	accountId uint,
) ([]*models.FriendGroup, error) {
	db := s.mysqlClient.Db()

	friendGroups := make([]*models.FriendGroup, 0)
	err := db.Model(&models.FriendGroup{}).
		Where("account_id = ?", accountId).
		Preload("Friends").
		Preload("Friends.Account").
		Find(&friendGroups).Error
	if err != nil {
		s.log.Errorf("GetFriendGroups %v", err)
		return nil, err
	}

	for _, v := range friendGroups {
		v.CreatedAt = nil
		v.UpdatedAt = nil
		for _, v2 := range v.Friends {
			v2.CreatedAt = nil
			v2.UpdatedAt = nil
			if v2.Account != nil {
				v2.Account.CreatedAt = nil
				v2.Account.UpdatedAt = nil
			}
		}
	}

	return friendGroups, nil
}

// 获取指定好友分组
func (s *Service) GetFriendGroup(
	ctx context.Context,
	accountId uint,
	req *models.FriendGroupReq,
) (*models.FriendGroup, error) {
	db := s.mysqlClient.Db()

	friendGroup := &models.FriendGroup{}
	err := db.Model(&models.FriendGroup{}).
		Where("id = ?", req.GroupId).
		Where("account_id = ?", accountId).
		Preload("Friends").
		Preload("Friends.Account").
		Find(friendGroup).Error
	if err != nil {
		s.log.Errorf("GetFriendGroup %v", err)
		return nil, err
	}
	if friendGroup.ID == 0 {
		return nil, errors.New(resp.FRIEND_GROUP_NOT_EXISTS)
	}

	friendGroup.CreatedAt = nil
	friendGroup.UpdatedAt = nil
	for _, v := range friendGroup.Friends {
		v.CreatedAt = nil
		v.UpdatedAt = nil
		if v.Account != nil {
			v.Account.CreatedAt = nil
			v.Account.UpdatedAt = nil
		}
	}

	return friendGroup, nil
}

// 重命名好友分组
func (s *Service) RenameFriendGroup(
	ctx context.Context,
	accountId uint,
	req *models.RenameFriendGroupReq,
) error {
	db := s.mysqlClient.Db()

	friendGroup := &models.FriendGroup{}
	err := db.Model(&models.FriendGroup{}).
		Where("id <> ?", req.GroupId).
		Where("account_id = ?", accountId).
		Where("name = ?", req.Name).
		First(friendGroup).Error
	if err != nil {
		s.log.Errorf("RenameFriendGroup select %v", err)
		return err
	}
	if friendGroup.ID != 0 {
		return errors.New(resp.FRIEND_GROUP_EXISTS)
	}

	err = db.Model(&models.FriendGroup{}).
		Where("id = ?", req.GroupId).
		Where("account_id = ?", accountId).
		Update("name", req.Name).Error
	if err != nil {
		s.log.Errorf("RenameFriendGroup update %v", err)
		return err
	}

	return nil
}

// 校验好友(或黑名单)
func (s *Service) VerifyFriend(
	ctx context.Context,
	accountId uint,
	req *models.ToIDReq,
) (*models.VerifyFriendRes, error) {
	res := &models.VerifyFriendRes{}

	friend := &models.Friend{}
	err := s.mysqlClient.Db().
		Model(&models.Friend{}).
		Where("account_id = ?", req.ToID).
		Where("friend_id = ?", accountId).
		First(friend).Error
	if err != nil {
		s.log.Errorf("VerifyFriend select %v", err)
		return nil, err
	}

	if friend.ID != 0 {
		res.IsFriend = true
		return res, nil
	}
	if friend.Blacklist {
		res.IsBlacklist = true
	}

	return res, nil
}
