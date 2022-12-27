package models

type Discover struct {
	Model
	Name  string `gorm:"column:name;not null;default:'';type:varchar(100)" json:"name"`
	Logo  string `gorm:"column:logo;not null;default:'';type:varchar(1000)" json:"logo"`
	Url   string `gorm:"column:url;not null;default:'';type:varchar(1000)" json:"url"`
	Order int    `gorm:"column:order;not null;default:0" json:"-"`
}

func (Discover) TableName() string {
	return "im_discovers"
}
