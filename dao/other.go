package dao

import (
	"context"
	"im/models"
	"im/pkg/database"
)

type Other struct {
	mysqlClient *database.Client
}

func NewOther(mysqlClient *database.Client) *Other {
	return &Other{
		mysqlClient: mysqlClient,
	}
}

func (s *Other) GetDiscovers(ctx context.Context) ([]*models.Discover, error) {
	db := s.mysqlClient.Db()
	discovers := make([]*models.Discover, 0)
	err := db.Model(&models.Discover{}).
		Order("`order` asc").
		Find(&discovers).Error
	return discovers, err
}
