package system

import (
	"context"

	"im/models"
	"im/pkg/database"
	"im/pkg/qiniu"
	"im/services/system/config"

	"go.uber.org/zap"
)

type Service struct {
	log         *zap.SugaredLogger
	config      *config.Config
	mysqlClient *database.Client
	qiniuClient *qiniu.Qiniu
}

func NewService(
	config *config.Config,
	mysqlClient *database.Client,
	qiniuClient *qiniu.Qiniu,
) *Service {
	return &Service{
		log:         zap.S().With("module", "services.system.service"),
		config:      config,
		mysqlClient: mysqlClient,
		qiniuClient: qiniuClient,
	}
}

// 创建操作日志
func (s *Service) CreateOperateLog(
	ctx context.Context,
	log *models.OperateLogs,
) error {
	db := s.mysqlClient.Db()
	err := db.Model(&models.OperateLogs{}).
		Create(log).
		Error
	if err != nil {
		return err
	}
	return nil
}

// 获取配置表
func (s *Service) GetConfig() *config.Config {
	return s.config
}

// 获取七牛云上传Token
func (s *Service) GetQiniuUploadToken() (string, string) {
	token := s.qiniuClient.GetUploadToken()
	domain := s.qiniuClient.GetDomain()
	return token, domain
}
