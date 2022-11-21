package models

import (
	"context"

	"gorm.io/gorm"
)

type Config struct {
	Model
	Key string `gorm:"column:key;not null;default:'';type:varchar(50);index:key"`
	Val string `gorm:"column:val;not null;default:'';type:varchar(2000)"`
}

func (c Config) GetAll(ctx context.Context, db *gorm.DB) ([]*Config, error) {
	cfgs := make([]*Config, 0)
	err := db.Model(&Config{}).
		Order("id asc").
		Find(&cfgs).
		Error
	return cfgs, err
}
