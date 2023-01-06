package dao

import "im/pkg/database"

type Chat struct {
	mysqlClient *database.Client
}

func NewChat(mysqlClient *database.Client) *Chat {
	return &Chat{
		mysqlClient: mysqlClient,
	}
}
