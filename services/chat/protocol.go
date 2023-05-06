package chat

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"im/models"
	"im/models/pb"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Protocol string

var (
	ProtocolJson  Protocol = "json"
	ProtocolProto Protocol = "proto"
)

var websocketUpgrade = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	// 取消ws跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 发送(协议层)
func (s *Service) SendProtocol(conn *Connection, msg *models.Message) error {
	switch conn.Protocol {
	case ProtocolJson:
		err := conn.Conn.WriteJSON(msg)
		if err != nil {
			return err
		}
	case ProtocolProto:
		data, err := s.ToMessageBytes(conn.Protocol, msg)
		if err != nil {
			return err
		}
		err = conn.Conn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			return err
		}
	}

	return nil
}

// 消息字节转为结构体
func (s *Service) ToMessage(protocol Protocol, msg []byte) (*models.Message, error) {
	message := &models.Message{}
	var err error

	switch protocol {
	case ProtocolJson:
		err = json.Unmarshal(msg, message)
	case ProtocolProto:
		messageProto := &pb.Message{}
		if err := proto.Unmarshal(msg, messageProto); err != nil {
			return nil, err
		}
		message.ID = uint(messageProto.Id)
		message.FromId = uint(messageProto.FromId)
		message.ToId = uint(messageProto.ToId)
		message.Ope = models.MessageOpe(messageProto.Ope)
		message.Type = models.MessageType(messageProto.Type)
		message.Body = messageProto.Body
		message.IsPrivate = messageProto.IsPrivate
		message.Status = models.MessageStatus(messageProto.Status)
		message.IsRead = messageProto.IsRead
		if messageProto.CreatedAt != nil {
			createdAt := messageProto.CreatedAt.AsTime()
			message.CreatedAt = &createdAt
		}
	}

	if message.Type == "" || message.Type != models.MessageTypePing {
		return nil, errors.New("unsupported message")
	}

	return message, err
}

// 消息结构体转为字节
func (s *Service) ToMessageBytes(protocol Protocol, message *models.Message) (data []byte, err error) {
	switch protocol {
	case ProtocolJson:
		data, err = json.Marshal(message)
	case ProtocolProto:
		messageProto := &pb.Message{}
		messageProto.Id = int64(message.ID)
		messageProto.FromId = int64(message.FromId)
		messageProto.ToId = int64(message.ToId)
		messageProto.Ope = string(message.Ope)
		messageProto.Type = string(message.Type)
		messageProto.Body = message.Body
		messageProto.IsPrivate = message.IsPrivate
		messageProto.Status = string(message.Status)
		messageProto.IsRead = message.IsRead
		if messageProto.CreatedAt != nil {
			messageProto.CreatedAt = timestamppb.New(*message.CreatedAt)
		}
		data, err = proto.Marshal(messageProto)
	}

	return data, err
}
