package chat

import (
	"context"
	"errors"
	"im/models"
	"im/pkg/database"
	"im/pkg/resp"
	"im/services/chatgroup"
	"im/services/system/config"
	"sync"
	"time"

	"go.uber.org/atomic"
	"go.uber.org/zap"
)

type Service struct {
	log          *zap.SugaredLogger
	mysqlClient  *database.Client
	config       *config.Config
	chatGroupSvc *chatgroup.Service

	connections       map[uint][]*Connection // 连接表
	connectionsLocker sync.RWMutex           // 连接并发锁
	counter           atomic.Uint64
}

func NewService(
	mysqlClient *database.Client,
	config *config.Config,
	chatGroupSvc *chatgroup.Service,
) *Service {
	return &Service{
		log:          zap.S().With("module", "services.chat.service"),
		mysqlClient:  mysqlClient,
		config:       config,
		chatGroupSvc: chatGroupSvc,

		connections:       make(map[uint][]*Connection),
		connectionsLocker: sync.RWMutex{},
		counter:           atomic.Uint64{},
	}
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
		return nil, err
	}

	s.Send(ctx, message)
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
	s.Send(ctx, message)

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
