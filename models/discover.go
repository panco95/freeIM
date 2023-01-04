package models

// 发现页表（连接列表）
type Discover struct {
	Model
	Name     string `gorm:"column:name;not null;default:'';type:varchar(100)" json:"name"`  //名称
	Logo     string `gorm:"column:logo;not null;default:'';type:varchar(1000)" json:"logo"` //图标
	Url      string `gorm:"column:url;not null;default:'';type:varchar(1000)" json:"url"`   //链接地址
	Order    int    `gorm:"column:order;not null;default:0" json:"-"`                       //排序
	Password string `gorm:"column:password;not null;default:''" json:"password"`            //访问密码
}

func (Discover) TableName() string {
	return "im_discovers"
}
