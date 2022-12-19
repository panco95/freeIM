package models

type Config struct {
	Model
	Key string `gorm:"column:key;not null;default:'';type:varchar(50);index:key"`
	Val string `gorm:"column:val;not null;default:'';type:varchar(2000)"`
}

func (Config) TableName() string {
	return "im_configs"
}
