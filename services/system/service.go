package system

import (
	"context"

	"im/models"
	"im/pkg/database"
	"im/services/system/config"

	"go.uber.org/zap"
)

type Service struct {
	log         *zap.SugaredLogger
	config      *config.Config
	mysqlClient *database.Client
}

func NewService(
	config *config.Config,
	mysqlClient *database.Client,
) *Service {
	return &Service{
		log:         zap.S().With("module", "services.system.service"),
		config:      config,
		mysqlClient: mysqlClient,
	}
}

// CreateOperateLog 创建操作日志
func (s *Service) CreateOperateLog(
	ctx context.Context,
	log *models.OperateLogs,
) error {
	err := s.mysqlClient.Db().WithContext(ctx).
		Model(&models.OperateLogs{}).
		Create(log).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetConfig() *config.Config {
	return s.config
}
