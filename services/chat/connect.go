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
	Connected bool
}

// 连接websocket
func (s *Service) ConnectWebsocket(ctx *gin.Context) {
	w := ctx.Writer
	r := ctx.Request
	var conn *websocket.Conn
	var err error

	conn, err = websocketUpgrade.Upgrade(w, r, nil)
	if err != nil {
		s.log.Errorf("ConnectWebsocket websocketUpgrade.Upgrade %v", err)
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
		Connected: true,
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

	close(c.Channel)
	c.Ctx = nil
	c.Connected = false
	go s.UpdateOnlineStatus(c.AccountId, models.OnlineStatusOffline)

}

// 消息处理协程
func (s *Service) ConnChannel(c *Connection) {
	go s.UpdateOnlineStatus(c.AccountId, models.OnlineStatusOnline)

	go func() {
		for {
			if !c.Connected {
				return
			}
			err := c.Conn.WriteJSON(&models.Message{
				Type: models.MessageTypePing,
			})
			if err != nil {
				s.log.Errorf("ConnChannel PING %v", err)
			}
			time.Sleep(time.Second * 10)
		}
	}()

	for msg := range c.Channel {
		message := &models.Message{}
		if err := json.Unmarshal(msg, message); err != nil {
			s.log.Error(err)
			continue
		}

		fn := s.msgRouter[string(message.Type)]
		err := fn(c, message)
		if err != nil {
			s.log.Errorf("ConnChannel msgRouter call %v", err)
		}
	}
}

// 生成客户端id
func (s *Service) GenClientId() string {
	return utils.Md5(fmt.Sprintf("%d", s.counter.Load()))
}

// 更新账号在线状态
func (s *Service) UpdateOnlineStatus(accountId uint, onlineStatus models.OnlineStatus) {
	db := s.mysqlClient.Db()
	err := db.Model(&models.Account{}).
		Where("id = ?", accountId).
		UpdateColumn("online_status", models.OnlineStatusOnline).
		Error
	if err != nil {
		s.log.Errorf("ConnChannel UpdateOnlineStatus %v", err)
	}
}

var websocketUpgrade = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	// 取消ws跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
