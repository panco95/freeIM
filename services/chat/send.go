package chat

import (
	"context"
	"im/models"
)

func (s *Service) Send(ctx context.Context, msg *models.Message) {
	switch msg.Ope {
	case models.MessageOpeFriend, models.MessageOpeSystem:
		go s.SendToAccountId(ctx, msg.ToId, msg)
	case models.MessageOpeGroup:
		idList, err := s.chatGroupSvc.FindChatGroupManagerIDList(ctx, msg.ToId)
		if err != nil {
			s.log.Errorf("Send OpeGroup findMembers %v", err)
			return
		}
		for _, accountId := range idList {
			go s.SendToAccountId(ctx, accountId, msg)
		}
	}
}

func (s *Service) SendToAccountId(ctx context.Context, accountId uint, msg *models.Message) {
	if _, ok := s.connections[accountId]; ok {
		for _, conn := range s.connections[accountId] {
			err := conn.Conn.WriteJSON(msg)
			if err != nil {
				s.log.Errorf("Send OpeFriend %v", err)
			}
		}
	}
}
