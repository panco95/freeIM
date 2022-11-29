package models

type SendMessageReq struct {
	ToIDReq
	Ope       int    `form:"ope" binding:"required"`  //消息通道
	Type      int    `form:"type" binding:"required"` //消息类型
	Body      string `form:"body" binding:"required"` //消息内容
	IsPrivate bool   `form:"isPrivate"`               //是否加密(私密聊天)
}

type RevocationMessageReq struct {
	ToIDReq
	MessageId string `form:"id" binding:"required"`  //消息ID
	Ope       int    `form:"ope" binding:"required"` //消息通道
}
