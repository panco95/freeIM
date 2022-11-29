package chat

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

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

// 连接结构体
type Connection struct {
	ClientId  string
	AccountId uint
	Conn      *websocket.Conn
	Channel   chan []byte
	Platform  string
	Ctx       *gin.Context
}

// 消息结构体
type Message struct {
	Id     string `json:"id"`     //消息ID
	FromId uint   `json:"fromId"` //来源id
	ToId   uint   `json:"toId"`   //接收id
	Ope    int    `json:"ope"`    //消息通道
	Type   int    `json:"type"`   //消息类型
	Body   string `json:"body"`   //消息内容
}

// 消息通道
const (
	OpeFriend = 1 //好友消息
	OpeGroup  = 2 //群消息
	OpeSystem = 3 //系统消息
)

// 消息类型
const (
	TypeText       = 1  //文字
	TypePicture    = 2  //图片
	TypeVoice      = 3  //语音
	TypeVideo      = 4  //视频
	TypeGeo        = 5  //地理位置
	TypeFile       = 6  //文件
	TypeRevocation = -1 //撤回
)
