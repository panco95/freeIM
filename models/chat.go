package models

type Message struct {
	Model
	FromId    uint          `gorm:"column:from_id;not null;default:0;index:from_id" json:"fromId"` //发送方id
	ToId      uint          `gorm:"column:to_id;not null;default:0;index:to_id" json:"toId"`       //接收方id
	Ope       MessageOpe    `gorm:"column:ope;not null;default:'friend'" json:"ope"`               //消息通道
	Type      MessageType   `gorm:"column:type;not null;default:'text'" json:"type"`               //消息类型
	Body      string        `gorm:"column:body;type:text" json:"body"`                             //消息内容
	IsPrivate bool          `gorm:"column:is_private;not null;default:0" json:"isPrivate"`         //是否私密消息
	Status    MessageStatus `gorm:"column:status;not null;default:'normal'" json:"status"`         //消息状态
}

type MessageStatus string
type MessageType string
type MessageOpe string

var (
	MessageStatusNormal     MessageStatus = "normal"     //正常
	MessageStatusRevocation MessageStatus = "revocation" //已撤回

	MessageTypeText  MessageType = "text"  //文字
	MessageTypePic   MessageType = "pic"   //图片
	MessageTypeVoice MessageType = "voice" //语音
	MessageTypeVideo MessageType = "video" //视频
	MessageTypeGeo   MessageType = "geo"   //地理位置
	MessageTypeFile  MessageType = "file"  //文件

	MessageOpeFriend MessageOpe = "friend" //好友
	MessageOpeGroup  MessageOpe = "group"  //群聊
	MessageOpeSystem MessageOpe = "system" //系统消息
)

type SendMessageReq struct {
	ToIDReq
	Ope       MessageOpe  `form:"ope" binding:"required"`  //消息通道
	Type      MessageType `form:"type" binding:"required"` //消息类型
	Body      string      `form:"body" binding:"required"` //消息内容
	IsPrivate bool        `form:"isPrivate"`               //是否加密(私密聊天)
}

type RevocationMessageReq struct {
	ToIDReq
	Id  uint       `form:"id" binding:"required"`  //消息ID
	Ope MessageOpe `form:"ope" binding:"required"` //消息通道
}
