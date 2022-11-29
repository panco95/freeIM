package chat

import (
	"context"
)

func (s *Service) Send(ctx context.Context, msg *Message) {
	switch msg.Ope {
	case OpeFriend, OpeSystem:
		go s.SendToAccountId(ctx, msg.ToId, msg)
	case OpeGroup:
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

func (s *Service) SendToAccountId(ctx context.Context, accountId uint, msg *Message) {
	if _, ok := s.connections[accountId]; ok {
		for _, conn := range s.connections[accountId] {
			err := conn.Conn.WriteJSON(msg)
			if err != nil {
				s.log.Errorf("Send OpeFriend %v", err)
			}
		}
	}
}
