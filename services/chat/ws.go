package chat

import (
	"encoding/json"
	"fmt"
	"im/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

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

	uid := ctx.GetUint("id")

	c := Connection{
		Conn:    conn,
		Channel: make(chan []byte),
		Ctx:     ctx,
		Uid:     uid,
	}
	go s.ConnChannel(&c)

	s.connectionsLocker.Lock()
	s.counter.Inc()
	c.ClientId = s.GenClientId()
	s.connections[uid] = &c
	s.connectionsLocker.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		c.Channel <- msg
	}

	s.connectionsLocker.Lock()
	delete(s.connections, uid)
	s.connectionsLocker.Unlock()
}

// 消息处理协程
func (s *Service) ConnChannel(c *Connection) {
	for msg := range c.Channel {
		s.log.Infof("%d: %s", c.Uid, string(msg))

		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			continue
		}
	}
}

// 生成客户端id
func (s *Service) GenClientId() string {
	return utils.Md5(fmt.Sprintf("%d", s.counter.Load()))
}
