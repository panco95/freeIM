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
	ClientId string
	Uid      uint
	Conn     *websocket.Conn
	Channel  chan []byte
	Ctx      *gin.Context
}

// 消息结构体
type Message struct {
	Cmd    int    //指令
	FromId int    //来源id
	ToId   int    //接收id
	Ope    int    //消息通道
	Type   int    //消息类型
	Body   string //消息内容
}

// 消息指令定义
const (
	CmdSuccess                = 1  //通用失败
	CmdFail                   = 2  //通用失败
	CmdAck                    = 3  //通用失败
	CmdSignSuccess            = 4  //登录成功
	CmdReceiveFriendMessage   = 5  //收到好友消息
	CmdWithdrawFriendMessage  = 6  //撤回好友消息
	CmdReceiveFriendAdd       = 7  //收到好友添加请求
	CmdReceiveFriendAddResult = 8  //收到好友请求结果
	CmdReceiveGroupMessage    = 9  //收到群消息
	CmdWithdrawGroupMessage   = 10 //撤回群消息
	CmdReceiveGroupJoin       = 11 //收到加入群组请求
	CmdReceiveGroupJoinResult = 12 //收到加入群组结果
	CmdReceiveGroupShot       = 13 //收到被踢出群组通知
)

// 消息通道
const (
	OpeFriend = 0 //好友消息
	OpeGroup  = 1 //群消息
	OpeSystem = 2 //系统消息
)

// 消息类型
const (
	TypeText    = 0
	TypePicture = 1
	TypeVoice   = 2
	TypeVideo   = 3
	TypeGeo     = 4
	TypeFile    = 6
	TypePrompt  = 10
)
