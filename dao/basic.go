package dao

import "im/pkg/database"

type Dao struct {
	Account   *Account
	Friend    *Friend
	ChatGroup *ChatGroup
	Chat      *Chat
	Config    *Config
	Other     *Other
}

func NewDao(mysqlClient *database.Client) *Dao {
	return &Dao{
		Account:   NewAccount(mysqlClient),
		Friend:    NewFriend(mysqlClient),
		ChatGroup: NewChatGroup(mysqlClient),
		Chat:      NewChat(mysqlClient),
		Config:    NewConfig(mysqlClient),
		Other:     NewOther(mysqlClient),
	}
}
