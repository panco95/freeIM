package chat

import (
	"context"
	"errors"
	"im/models"
	"im/pkg/database"
	"im/pkg/etcd"
	"im/pkg/gin/middlewares"
	"im/pkg/resp"
	"im/services/system/config"
	"sync"
	"time"

	"go.uber.org/atomic"
	"go.uber.org/zap"
)

type MsgFunc func(conn *Connection, message *models.Message) error

type Service struct {
	log         *zap.SugaredLogger
	mysqlClient *database.Client
	manager     *etcd.Manager
	config      *config.Config

	connections       map[uint][]*Connection // 连接表
	connectionsLocker sync.RWMutex           // 连接并发锁
	counter           atomic.Uint64          // 链接计数器
	msgRouter         map[string]MsgFunc     // 消息路由表

	RPC *RPC
}

func NewService(
	mysqlClient *database.Client,
	manager *etcd.Manager,
	config *config.Config,
) *Service {
	s := &Service{
		log:         zap.S().With("module", "services.chat.service"),
		mysqlClient: mysqlClient,
		manager:     manager,
		config:      config,

		connections:       make(map[uint][]*Connection),
		connectionsLocker: sync.RWMutex{},
		counter:           atomic.Uint64{},
		msgRouter:         map[string]MsgFunc{},
	}

	s.SetRouters()
	return s
}

// 发送消息
func (s *Service) SendMessage(
	ctx context.Context,
	accountId uint,
	req *models.SendMessageReq,
) (*models.Message, error) {
	db := s.mysqlClient.Db()
	message := &models.Message{
		FromId:    accountId,
		ToId:      req.ToID,
		Ope:       req.Ope,
		Type:      req.Type,
		Body:      req.Body,
		IsPrivate: req.IsPrivate,
		Status:    models.MessageStatusNormal,
	}
	err := db.Create(message).Error
	if err != nil {
		s.log.Errorf("SendMessage Create %v", err)
		return nil, err
	}

	err = s.RPC.SendMessageCall(context.Background(), message)
	if err != nil {
		s.log.Errorf("SendMessage RPC.SendMessageCall %v", err)
		return nil, err
	}

	return message, nil
}

// 撤回消息
func (s *Service) RevocationMessage(
	ctx context.Context,
	accountId uint,
	req *models.RevocationMessageReq,
) error {
	db := s.mysqlClient.Db()
	message := &models.Message{}
	err := db.Model(&models.Message{}).
		Where("id = ?", req.Id).
		First(message).Error
	if err != nil {
		s.log.Errorf("RevocationMessage select %v", err)
		return err
	}
	if message.FromId != accountId {
		return errors.New(resp.MESSAGE_NOT_YOUR)
	}
	if message.CreatedAt.Unix()+120 < time.Now().Unix() {
		return errors.New(resp.MESSAGE_CANT_REVOCATION)
	}

	message.Status = models.MessageStatusRevocation
	err = s.RPC.SendMessageCall(context.Background(), message)
	if err != nil {
		s.log.Errorf("SendMessage RPC.SendMessageCall %v", err)
		return err
	}

	go func() {
		err := db.Model(&models.Message{}).
			Where("id = ?", req.Id).
			UpdateColumn("status", models.MessageStatusRevocation).
			Error
		if err != nil {
			s.log.Errorf("RevocationMessage updateStatus %v", err)
		}
	}()

	return nil
}

// 已读消息
func (s *Service) ReadMessage(
	ctx context.Context,
	accountId uint,
	req *models.ToIDReq,
) error {
	db := s.mysqlClient.Db()
	err := db.Model(&models.Message{}).
		Where("from_id = ?", req.ToID).
		Where("to_id = ?", accountId).
		UpdateColumn("is_read", true).Error
	if err != nil {
		s.log.Errorf("ReadMessage update %v", err)
		return err
	}

	err = s.RPC.SendMessageCall(context.Background(), &models.Message{
		FromId: accountId,
		ToId:   req.ToID,
		Ope:    models.MessageOpeFriend,
		Type:   models.MessageTypeRead,
		IsRead: true,
	})
	if err != nil {
		s.log.Errorf("ReadMessage RPC.SendMessageCall %v", err)
		return err
	}

	return nil
}

// 获取聊天记录
func (s *Service) GetMessagLogs(
	ctx context.Context,
	accountId uint,
	req *models.GetMessagesReq,
	page *middlewares.Pagination,
) ([]*models.MessageLog, int64, error) {
	db := s.mysqlClient.Db()
	query := db.Model(&models.Message{}).
		Where("(from_id = ? and to_id = ?) or (from_id = ? and to_id = ?)", accountId, req.ToID, req.ToID, accountId).
		Where("status <> ?", models.MessageStatusRevocation)
	if req.MessageType != "" {
		query = query.Where("type = ?", req.MessageType)
	}

	var total int64
	msgs := make([]*models.Message, 0)

	err := page.GetDocsAndTotal(query, &msgs, &total)
	if err != nil {
		s.log.Errorf("GetMessagLogs select %v", err)
		return nil, 0, err
	}

	logs := make([]*models.MessageLog, 0)
	for _, v := range msgs {
		v.UpdatedAt = nil
		logs = append(logs, &models.MessageLog{
			Message: v,
			Self:    accountId == v.FromId,
		})
	}

	return logs, total, nil
}
