package dao

import "im/pkg/database"

type ChatGroup struct {
	mysqlClient *database.Client
}

func NewChatGroup(mysqlClient *database.Client) *ChatGroup {
	return &ChatGroup{
		mysqlClient: mysqlClient,
	}
}
