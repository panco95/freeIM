package chat

import (
	"context"
	"im/models"
	"im/pkg/database"
	"im/services/chatgroup"
	"im/services/system/config"
	"sync"

	"github.com/google/uuid"
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
) error {
	msg := &Message{
		Id:     uuid.New().String(),
		FromId: accountId,
		ToId:   req.ToID,
		Ope:    req.Ope,
		Type:   req.Type,
		Body:   req.Body,
	}
	s.Send(ctx, msg)
	return nil
}

// 撤回消息
func (s *Service) RevocationMessage(
	ctx context.Context,
	accountId uint,
	req *models.RevocationMessageReq,
) error {
	msg := &Message{
		Id:     req.MessageId,
		FromId: accountId,
		ToId:   req.ToID,
		Ope:    req.Ope,
		Type:   TypeRevocation,
	}
	s.Send(ctx, msg)
	return nil
}
