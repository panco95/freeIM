package models

type OperateLogs struct {
	Model
	AccountID    uint             `gorm:"column:account_id;not null;default:0" json:"accountId"`
	Module       string           `gorm:"column:module;not null;default:'';type:varchar(50)" json:"module"`
	IP           string           `gorm:"column:ip;not null;default:'';type:varchar(50)" json:"ip"`
	Content      string           `gorm:"column:content;type:varchar(10000)" json:"content"`
	Detail       string           `gorm:"column:detail;type:longtext;comment:操作详情" json:"detail"`
	Fields       ArrayFieldString `gorm:"column:fields;type:longtext;comment:字段" json:"fields"`
	FieldsBefore ArrayFieldString `gorm:"column:fields_before;type:json;comment:修改前字段" json:"beforeFields"`
	FieldsAfter  ArrayFieldString `gorm:"column:fields_after;type:json;comment:修改后字段" json:"afterFields"`
}

func (OperateLogs) TableName() string {
	return "im_operate_logs"
}
