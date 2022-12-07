package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"im/models"
	"im/pkg/gin/middlewares"
	"im/pkg/validator"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	rpcx_logger "github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	_ "google.golang.org/grpc/encoding/gzip" // Install the gzip compressor
)

type serverFlags struct {
	port string
}

var sFlags serverFlags

func init() {
	serverCmd.Flags().StringVar(&sFlags.port, "port", "8008", "listen port")
	_ = viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))

	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"server"},
	Short:   "Start server application",
	Run: func(cmd *cobra.Command, args []string) {
		log := zap.S().With("cmd", "server")
		pkgs := NewPackages("server")
		svcs := NewServices(pkgs)
		ginCtrls := NewGinControllers(pkgs, svcs)

		// 创建mysql数据库表
		if err := pkgs.mysqlClient.AutoMigrate(
			&models.Account{},
			&models.Config{},
			&models.Friend{},
			&models.FriendApply{},
			&models.FriendGroup{},
			&models.ChatGroup{},
			&models.ChatGroupMember{},
			&models.ChatGroupJoin{},
			&models.Message{},
			&models.OperateLogs{},
		); err != nil {
			log.Fatalf("Mysql AutoMigrate %v", err)
		}
		if err := pkgs.mysqlClient.SetAutoIncrementID("accounts", 10000000); err != nil {
			log.Errorf("Mysql SetAutoIncrementID `accounts` %v", err)
		}

		// 自动刷新配置
		viper.SetDefault("config.autoRefreshInterval", "10s")
		go svcs.systemSvc.GetConfig().AutoRefresh(viper.GetDuration("config.autoRefreshInterval"))

		var httpPublicServer *http.Server
		var rpcPrivateServer *server.Server

		var eg errgroup.Group
		eg.Go(func() error {
			engine, err := GetGinPublicEngine(ginCtrls, pkgs)
			if err != nil {
				return err
			}

			port := ":" + viper.GetString("http.public.port")
			log.Infof("HTTP public server listen on %s", port)
			httpPublicServer = &http.Server{
				Addr:    port,
				Handler: engine,
			}

			return httpPublicServer.ListenAndServe()
		})

		eg.Go(func() error {
			rpcPrivateServer = server.NewServer()
			rpcx_logger.SetLogger(log)
			if err := rpcPrivateServer.RegisterName(svcs.rpc.ProjectName, svcs.rpc, ""); err != nil {
				return err
			}
			port := ":" + viper.GetString("rpc.private.port")
			log.Infof("RPC private server listen on: %s", port)
			if err := rpcPrivateServer.Serve("tcp", port); err != nil {
				return err
			}

			return nil
		})

		eg.Go(func() error {
			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
			sig := <-signalChan

			log.Warnf("System about to exit because of signal=%s", sig.String())

			var wg sync.WaitGroup
			servers := []*http.Server{
				httpPublicServer,
			}
			rpcServers := []*server.Server{
				rpcPrivateServer,
			}

			for _, svr := range servers {
				if svr == nil {
					continue
				}

				svr := svr
				wg.Add(1)

				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					if err := svr.Shutdown(ctx); err != nil {
						log.Errorf("HTTP server shutdown %v", err)
					}
					wg.Done()
				}()
			}

			for _, svr := range rpcServers {
				if svr == nil {
					continue
				}

				svr := svr
				wg.Add(1)

				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					if err := svr.Shutdown(ctx); err != nil {
						log.Errorf("RPC server shutdown %v", err)
					}
					wg.Done()
				}()
			}

			wg.Wait()

			return nil
		})

		if err := eg.Wait(); err != nil {
			log.Fatalf("Server %v", err)
		}
	},
}

func GetGinPublicEngine(ctrls *GinControllers, pkgs *Packages) (*gin.Engine, error) {
	if os.Getenv("GO_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	_ = router.SetTrustedProxies(viper.GetStringSlice("http.trustedProxies"))

	router.Use(gin.Recovery())
	router.Use(cors.AllowAll())
	router.Use(middlewares.HTTPGzipEncoding)

	api := router.Group("/api/v1")
	pkgs.prom.Use(router)
	api.Use(pkgs.prom.Instrument("public"))
	pprof.RouteRegister(api)

	gzipExcludePaths := []string{}
	api.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPaths(gzipExcludePaths)))

	uni, err := validator.GetUniversalTranslator()
	if err != nil {
		return nil, err
	}
	eh := middlewares.NewWithStatusHandler(uni)
	api.Use(eh.HandleErrors)

	api.Use(middlewares.Logger(zap.S()))
	api.Use(middlewares.NewPaginationMiddleware())
	api.Use(middlewares.NewI18nMiddleware())
	api.Use(middlewares.Tracing(middlewares.TracingComponentName("gin")))
	api.Use(middlewares.NewOperateLogger(zap.S(), pkgs.mysqlClient))

	captcha := api.Group("captcha")
	captcha.GET("image", ctrls.accountCtrl.ImageCaptcha)
	captcha.GET("email", ctrls.accountCtrl.EmailCaptcha)
	captcha.GET("mobile", ctrls.accountCtrl.MobileCaptcha)
	register := api.Group("register")
	register.POST("basic", ctrls.accountCtrl.BasicRegister)
	login := api.Group("login")
	login.POST("basic", ctrls.accountCtrl.BasicLogin)
	login.POST("email", ctrls.accountCtrl.EmailLogin)
	login.POST("mobile", ctrls.accountCtrl.MobileLogin)
	reset := api.Group("reset")
	reset.POST("password/email", ctrls.accountCtrl.EmailResetPassword)
	reset.POST("password/mobile", ctrls.accountCtrl.MobileResetPassword)

	api.Use(middlewares.NewJwtCheckMiddleware(pkgs.jwt, pkgs.mysqlClient, pkgs.cacheClient))
	api.GET("ws", ctrls.chatCtrl.ConnectWebsocket)                      //连接webscket
	api.GET("upload/qiniu/token", ctrls.systemCtrl.GetQiniuUploadToken) //获取七牛云上传Token

	api.GET("me/info", ctrls.accountCtrl.Info)                      //个人信息
	api.PUT("me/update/password", ctrls.accountCtrl.UpdatePassword) //修改密码
	api.PUT("me/update/info", ctrls.accountCtrl.UpdateAccountInfo)  //更新个人信息

	api.GET("friends/search", ctrls.friendCtrl.SearchFriend)                     //查找用户
	api.POST("friends/add", ctrls.friendCtrl.AddFriend)                          //添加好友
	api.POST("friends/add/reply", ctrls.friendCtrl.AddFriendReply)               //同意/拒绝好友请求
	api.GET("friends/add/applies", ctrls.friendCtrl.FriendApplyList)             //好友申请列表
	api.GET("friends", ctrls.friendCtrl.FriendList)                              //好友(黑名单)列表
	api.GET("friends/info", ctrls.friendCtrl.FriendInfo)                         //单个好友信息
	api.DELETE("friends", ctrls.friendCtrl.DeleteFriend)                         //删除好友
	api.POST("friends/blaklist", ctrls.friendCtrl.AddBlacklist)                  //添加黑名单
	api.DELETE("friends/blaklist", ctrls.friendCtrl.DeleteBlacklist)             //移除黑名单
	api.GET("friends/verify", ctrls.friendCtrl.VerifyFriend)                     //校验好友（是否你是他的好友或者黑名单）
	api.PUT("friends/remark", ctrls.friendCtrl.SetFriendRemark)                  //设置好友备注
	api.PUT("friends/label", ctrls.friendCtrl.SetFriendLabel)                    //设置好友标签(自定义字段)
	api.POST("friends/groups", ctrls.friendCtrl.CreateFriendGroup)               //创建好友分组
	api.DELETE("friends/groups", ctrls.friendCtrl.DeleteFriendGroup)             //删除好友分组
	api.POST("friends/groups/members", ctrls.friendCtrl.AddFriendGroupMembers)   //好友分组添加成员
	api.DELETE("friends/groups/members", ctrls.friendCtrl.DelFriendGroupMembers) //好友分组删除成员
	api.GET("friends/groups", ctrls.friendCtrl.GetFriendGroups)                  //获取好友分组列表
	api.GET("friends/group", ctrls.friendCtrl.GetFriendGroup)                    //获取指定好友分组
	api.PUT("friends/groups/name", ctrls.friendCtrl.RenameFriendGroup)           //重命名好友分组

	api.POST("chatGroups/search", ctrls.chatGroupCtrl.SearchChatGroup)        //搜索群聊
	api.POST("chatGroups", ctrls.chatGroupCtrl.CreateChatGroup)               //创建群聊
	api.PUT("chatGroups", ctrls.chatGroupCtrl.EditChatGroup)                  //修改群聊资料
	api.POST("chatGroups/join", ctrls.chatGroupCtrl.JoinChatGroup)            //加群申请
	api.POST("chatGroups/join/reply", ctrls.chatGroupCtrl.JoinChatGroupReply) //加群审批
	api.GET("chatGroups/join", ctrls.chatGroupCtrl.JoinChatGroupList)         //加群审批列表
	api.GET("chatGroups", ctrls.chatGroupCtrl.ChatGroupList)                  //我的群聊列表
	api.GET("chatGroups/info", ctrls.chatGroupCtrl.ChatGroupInfo)             //群聊信息(包括成员)
	api.POST("chatGroups/dissolve", ctrls.chatGroupCtrl.DissolveChatGroup)    //解散群聊
	api.POST("chatGroups/exit", ctrls.chatGroupCtrl.ExitChatGroup)            //退出群聊
	api.POST("chatGroups/transfer", ctrls.chatGroupCtrl.TransferChatGroup)    //转让群聊
	api.POST("chatGroups/kick", ctrls.chatGroupCtrl.ChatGroupKickMember)      //踢出群员
	api.POST("chatGroups/manager", ctrls.chatGroupCtrl.ChatGroupSetManager)   //设置管理员
	api.POST("chatGroups/banned", ctrls.chatGroupCtrl.ChatGroupBannedMember)  //成员禁言

	api.POST("chat/send", ctrls.chatCtrl.SendMessage)             //发送消息
	api.POST("chat/revocation", ctrls.chatCtrl.RevocationMessage) //撤回消息
	api.POST("chat/read", ctrls.chatCtrl.ReadMessage)             //已读消息
	api.GET("chat/messages", ctrls.chatCtrl.GetMessagLogs)        //获取聊天记录

	return router, nil
}
