package cmd

import (
	"im/dao"
	"im/pkg/database"
	"im/pkg/email"
	"im/pkg/etcd"
	"im/pkg/gin/middlewares"
	"im/pkg/jwt"
	"im/pkg/qiniu"
	"im/pkg/sms"
	"im/pkg/tracing"
	"im/pkg/utils"
	"im/services/account"
	"im/services/chat"
	"im/services/chatgroup"
	"im/services/friend"
	"im/services/system"
	"im/services/system/config"

	redisCache "github.com/go-redis/cache/v8"
	redislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type Packages struct {
	mysqlClient   *database.Client
	redisClient   *redislib.Client
	cacheClient   *redisCache.Cache
	redSyncClient *redsync.Redsync
	dao           *dao.Dao
	etcdClient    *etcd.Etcd
	manager       *etcd.Manager
	emailClient   *email.Mail
	smsClient     sms.Sms
	qiniuClient   *qiniu.Qiniu
	prom          *middlewares.Prometheus
	tracing       *tracing.TracingService
	jwt           *jwt.Jwt
}

func NewPackages(serverID string) (pkgs *Packages) {
	pkgs = &Packages{
		prom: middlewares.NewPromMiddleware(),
	}
	log := zap.S().With("module", "init")

	{
		pkgs.tracing = tracing.NewTracingService(
			viper.GetBool("tracing.ext.logging.enable"),
		)

		err := pkgs.tracing.InitGlobal(viper.Sub("tracing"))
		if err != nil {
			log.Fatalf("Init tracing err %+v", err)
		}
	}

	{
		viper.SetDefault("mysql.logLevel", logger.Info)
		mysqlClient, err := database.NewMysql(
			viper.GetString("mysql.addr"),
			viper.GetInt("mysql.maxIdleConns"),
			viper.GetInt("mysql.maxOpenConns"),
			viper.GetDuration("mysql.connMaxLifetime"),
			logger.LogLevel(viper.GetUint("mysql.logLevel")),
		)
		if err != nil {
			log.Fatalf("Init mysql %v", err)
		}
		pkgs.mysqlClient = mysqlClient
	}

	{
		viper.SetDefault("jwt.key", "suanzi")
		viper.SetDefault("jwt.issue", "SZKJ")
		pkgs.jwt = jwt.New(
			[]byte(viper.GetString("jwt.key")),
			viper.GetString("jwt.issue"),
		)
	}

	{
		viper.SetDefault("redis.uri", "192.168.16.131:6379")
		viper.SetDefault("redis.password", "")
		viper.SetDefault("redis.db", 0)
		rdb := redislib.NewClient(&redislib.Options{
			Addr:     viper.GetString("redis.uri"),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
		})
		pkgs.redisClient = rdb
		pkgs.cacheClient = redisCache.New(&redisCache.Options{
			Redis: rdb,
		})
		pkgs.redSyncClient = redsync.New(goredis.NewPool(pkgs.redisClient))
	}

	{
		viper.SetDefault("email.addr", "")
		viper.SetDefault("email.identity", "")
		viper.SetDefault("email.username", "")
		viper.SetDefault("email.password", "")
		viper.SetDefault("email.host", "")
		emailClient := email.New(
			viper.GetString("email.addr"),
			viper.GetString("email.identity"),
			viper.GetString("email.username"),
			viper.GetString("email.password"),
			viper.GetString("email.host"),
		)
		pkgs.emailClient = emailClient
	}

	// {
	// 	viper.SetDefault("aliyun.regionId", "")
	// 	viper.SetDefault("aliyun.accessKeyId", "")
	// 	viper.SetDefault("aliyun.accessKeySecret", "")
	// 	viper.SetDefault("aliyun.smsCaptcha.signName", "")
	// 	viper.SetDefault("aliyun.smsCaptcha.templateCode", "")
	// 	aliyunClient, err := aliyun.New(
	// 		viper.GetString("aliyun.regionId"),
	// 		viper.GetString("aliyun.accessKeyId"),
	// 		viper.GetString("aliyun.accessKeySecret"),
	// 		viper.GetString("aliyun.smsCaptcha.signName"),
	// 		viper.GetString("aliyun.smsCaptcha.templateCode"),
	// 	)
	// 	if err != nil {
	// 		log.Errorf("Init aliyun %v", err)
	// 	}
	// 	pkgs.aliyunClient = aliyunClient
	// }

	{
		viper.SetDefault("smsbao.username", "")
		viper.SetDefault("smsbao.password", "")
		viper.SetDefault("smsbao.sendRange", sms.SendRangeLocal)
		smsClient := sms.NewSmsBao(
			viper.GetString("smsbao.username"),
			viper.GetString("smsbao.password"),
			sms.SendRange(viper.GetString("smsbao.sendRange")),
		)
		pkgs.smsClient = smsClient
	}

	{
		viper.SetDefault("etcd.addrs", []string{})
		etcdClient, err := etcd.New(
			viper.GetStringSlice("etcd.addrs"),
			zap.S().With("client", "etcd").Desugar(),
		)
		if err != nil {
			log.Fatalf("Init etcd %v", err)
		}
		pkgs.etcdClient = etcdClient
	}

	{
		viper.SetDefault("etcd.prefix", "freeim")
		viper.SetDefault("local.addr", "")
		viper.SetDefault("rpc.private.port", "9000")
		viper.SetDefault("etcd.addrs", []string{})
		manager, err := etcd.NewManager(
			pkgs.etcdClient.Client,
			viper.GetString("local.addr"),
			viper.GetString("rpc.private.port"),
			viper.GetString("etcd.prefix"),
		)
		if err != nil {
			log.Fatalf("Init mananger %v", err)
		}
		pkgs.manager = manager
	}

	{
		viper.SetDefault("qiniu.ak", "")
		viper.SetDefault("qiniu.sk", "")
		viper.SetDefault("qiniu.bucket", "")
		viper.SetDefault("qiniu.domain", "")
		viper.SetDefault("qiniu.tokenExpireSec", 3600)
		qiniuClient := qiniu.New(
			viper.GetString("qiniu.ak"),
			viper.GetString("qiniu.sk"),
			viper.GetString("qiniu.bucket"),
			viper.GetString("qiniu.domain"),
			viper.GetUint64("qiniu.tokenExpireSec"),
		)
		pkgs.qiniuClient = qiniuClient
	}

	{
		dao := dao.NewDao(pkgs.mysqlClient)
		pkgs.dao = dao
	}

	return
}

type Services struct {
	accountSvc   *account.Service
	friendSvc    *friend.Service
	chatGroupSvc *chatgroup.Service
	chatSvc      *chat.Service
	systemSvc    *system.Service
	rpc          *chat.RPC
}

func NewServices(pkgs *Packages) *Services {
	config := config.NewConfig(pkgs.mysqlClient)

	systemSvc := system.NewService(
		config,
		pkgs.mysqlClient,
		pkgs.qiniuClient,
	)

	viper.SetDefault("login.rsaPrivateKey", utils.DefaultRSAPrivateKey)
	viper.SetDefault("login.rsaPublicKey", utils.DefaultRSAPublicKey)
	accountSvc := account.NewService(
		pkgs.redisClient,
		pkgs.cacheClient,
		pkgs.dao,
		pkgs.emailClient,
		pkgs.smsClient,
		pkgs.jwt,
		config,
	)

	chatSvc := chat.NewService(
		pkgs.mysqlClient,
		pkgs.manager,
		config,
	)

	friendSvc := friend.NewService(
		chatSvc,
		pkgs.mysqlClient,
		pkgs.redisClient,
		config,
	)

	chatGroupSvc := chatgroup.NewService(
		chatSvc,
		pkgs.mysqlClient,
		pkgs.redisClient,
		config,
	)

	viper.SetDefault("project.name", "freeim")
	rpc := chat.NewRPC(
		viper.GetString("project.name"),
		chatSvc,
	)

	chatSvc.RPC = rpc

	return &Services{
		accountSvc:   accountSvc,
		systemSvc:    systemSvc,
		chatSvc:      chatSvc,
		friendSvc:    friendSvc,
		chatGroupSvc: chatGroupSvc,
		rpc:          rpc,
	}
}

type GinControllers struct {
	accountCtrl   *account.GinController
	friendCtrl    *friend.GinController
	chatGroupCtrl *chatgroup.GinController
	chatCtrl      *chat.GinController
	systemCtrl    *system.GinController
}

func NewGinControllers(pkgs *Packages, svcs *Services) *GinControllers {
	return &GinControllers{
		accountCtrl: account.NewGinController(
			svcs.accountSvc,
			svcs.systemSvc,
		),
		friendCtrl: friend.NewGinController(
			svcs.friendSvc,
			svcs.systemSvc,
		),
		chatGroupCtrl: chatgroup.NewGinController(
			svcs.chatGroupSvc,
			svcs.systemSvc,
		),
		chatCtrl: chat.NewGinController(
			svcs.chatSvc,
			svcs.systemSvc,
		),
		systemCtrl: system.NewGinController(
			svcs.systemSvc,
		),
	}
}
