package dao

import "im/pkg/database"

type Config struct {
	mysqlClient *database.Client
}

func NewConfig(mysqlClient *database.Client) *Config {
	return &Config{
		mysqlClient: mysqlClient,
	}
}
