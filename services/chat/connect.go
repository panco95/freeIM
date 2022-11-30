package chat

import (
	"encoding/json"
	"fmt"
	"im/models"
	"im/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 连接结构体
type Connection struct {
	ClientId  string
	AccountId uint
	Conn      *websocket.Conn
	Channel   chan []byte
	Platform  string
	Ctx       *gin.Context
}

// 连接websocket
func (s *Service) ConnectWebsocket(ctx *gin.Context) {
	w := ctx.Writer
	r := ctx.Request
	var conn *websocket.Conn
	var err error

	conn, err = websocketUpgrade.Upgrade(w, r, nil)
	if err != nil {
		s.log.Errorf("websocketUpgrade.Upgrade %v", err)
		return
	}

	accountId := ctx.GetUint("id")
	platform := ctx.GetString("platform")

	c := Connection{
		Conn:      conn,
		Channel:   make(chan []byte),
		Ctx:       ctx,
		AccountId: accountId,
		Platform:  platform,
	}
	go s.ConnChannel(&c)

	s.connectionsLocker.Lock()
	s.counter.Inc()
	c.ClientId = s.GenClientId()
	if _, ok := s.connections[accountId]; !ok {
		s.connections[accountId] = make([]*Connection, 0)
	}
	s.connections[accountId] = append(s.connections[accountId], &c)
	s.connectionsLocker.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		c.Channel <- msg
	}

	s.connectionsLocker.Lock()
	delete(s.connections, accountId)
	s.connectionsLocker.Unlock()
}

// 消息处理协程
func (s *Service) ConnChannel(c *Connection) {
	for msg := range c.Channel {
		s.log.Infof("%d: %s", c.AccountId, string(msg))

		var message models.Message
		if err := json.Unmarshal(msg, &message); err != nil {
			continue
		}
	}
}

// 生成客户端id
func (s *Service) GenClientId() string {
	return utils.Md5(fmt.Sprintf("%d", s.counter.Load()))
}

// 升级http为websocket服务
var websocketUpgrade = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	// 取消ws跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
