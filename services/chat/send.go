package chat

import (
	"context"
	"im/models"
)

// 发送信息
func (s *Service) Send(ctx context.Context, msg *models.Message) {
	switch msg.Ope {
	case models.MessageOpeFriend, models.MessageOpeSystem:
		go s.SendToAccountId(ctx, msg.ToId, msg)
	case models.MessageOpeGroup:
		idList, err := s.FindChatGroupMemberIDList(ctx, msg.ToId)
		if err != nil {
			s.log.Errorf("Send OpeGroup findMembers %v", err)
			return
		}
		for _, accountId := range idList {
			go s.SendToAccountId(ctx, accountId, msg)
		}
	}
}

// 发送信息给用户
func (s *Service) SendToAccountId(ctx context.Context, accountId uint, msg *models.Message) {
	if _, ok := s.connections[accountId]; ok {
		for _, conn := range s.connections[accountId] {
			err := s.SendProtocol(conn, msg)
			if err != nil {
				s.log.Errorf("Send OpeFriend %v", err)
			}
		}
	}
}

// 找出群聊群员ID列表
func (s *Service) FindChatGroupMemberIDList(
	ctx context.Context,
	chatGroupId uint,
) ([]uint, error) {
	db := s.mysqlClient.Db()
	chatGroupMembers := make([]*models.ChatGroupMember, 0)
	err := db.Model(&models.ChatGroupMember{}).
		Where("chat_group_id = ?", chatGroupId).
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
