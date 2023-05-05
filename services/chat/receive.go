package chat

import (
	"im/models"
)

func (s *Service) SetRouter(name string, fn MsgFunc) {
	s.msgRouter[name] = fn
}

func (s *Service) SetRouters() {
	s.SetRouter(string(models.MessageTypePing), s.Ping)
}

func (s *Service) Ping(c *Connection, message *models.Message) error {
	message.FromId = c.AccountId
	err := c.Conn.WriteJSON(message)
	return err
}
