package dao

import "im/pkg/database"

type Friend struct {
	mysqlClient *database.Client
}

func NewFriend(mysqlClient *database.Client) *Friend {
	return &Friend{
		mysqlClient: mysqlClient,
	}
}
