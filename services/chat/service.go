package chat

import (
	"im/pkg/database"
	"im/pkg/jwt"
	"im/services/system/config"
	"sync"

	redislib "github.com/go-redis/redis/v8"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

type Service struct {
	log         *zap.SugaredLogger
	mysqlClient *database.Client
	redisClient *redislib.Client
	jwt         *jwt.Jwt
	config      *config.Config

	connections       map[uint]*Connection // 连接表
	connectionsLocker sync.RWMutex         // 连接并发锁
	counter           atomic.Uint64
}

func NewService(
	mysqlClient *database.Client,
	redisClient *redislib.Client,
	config *config.Config,
	jwt *jwt.Jwt,
) *Service {
	return &Service{
		log:         zap.S().With("module", "services.chat.service"),
		mysqlClient: mysqlClient,
		redisClient: redisClient,
		config:      config,
		jwt:         jwt,

		connections:       make(map[uint]*Connection),
		connectionsLocker: sync.RWMutex{},
		counter:           atomic.Uint64{},
	}
}
