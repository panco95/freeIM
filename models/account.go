package models

import (
	"time"
)

// 账号表
type Account struct {
	Model
	Username      string        `gorm:"column:username;not null;default:'';type:varchar(50);index:username" json:"username"`  //用户名
	Email         string        `gorm:"column:email;not null;default:'';type:varchar(100)" json:"email"`                      //邮箱地址
	Mobile        string        `gorm:"column:mobile;not null;default:'';type:varchar(50)" json:"mobile"`                     //手机号
	Nickname      string        `gorm:"column:nickname;not null;default:'';type:varchar(50)" json:"nickname"`                 //昵称
	Avatar        string        `gorm:"column:avatar;not null;default:'';type:varchar(500)" json:"avatar"`                    //头像
	Gender        string        `gorm:"column:gender;not null;default:'';type:varchar(10)" json:"gender"`                     //性别
	Birth         string        `gorm:"column:birth;default:null;type:varchar(100)" json:"birth"`                             //生日
	Age           int           `gorm:"-" json:"age"`                                                                         //年龄
	Intro         string        `gorm:"column:intro;not null;default:'';type:varchar(1000)" json:"intro"`                     //个人介绍
	Longitude     float64       `gorm:"column:longitude;not null;default:0;type:decimal(9,6)" json:"longitude"`               //经度
	Latitude      float64       `gorm:"column:latitude;not null;default:0;type:decimal(9,6)" json:"latitude"`                 //纬度
	Country       string        `gorm:"column:country;not null;default:'';type:varchar(50)" json:"country"`                   //国家
	Province      string        `gorm:"column:province;not null;default:'';type:varchar(50)" json:"province"`                 //省份
	City          string        `gorm:"column:city;not null;default:'';type:varchar(50)" json:"city"`                         //城市
	District      string        `gorm:"column:district;not null;default:'';type:varchar(50)" json:"district"`                 //区县
	Password      string        `gorm:"column:password;not null;default:'';type:varchar(200)" json:"-"`                       //密码
	PasswordSalt  string        `gorm:"column:password_salt;not null;default:'';type:varchar(200)" json:"-"`                  //密码盐值
	Status        AccountStatus `gorm:"column:status;not null;default:'normal';type:varchar(20)" json:"-"`                    //状态
	OnlineStatus  OnlineStatus  `gorm:"column:online_status;not null;default:'offline';type:varchar(20)" json:"onlineStatus"` //在线状态
	LastLoginTime *time.Time    `gorm:"column:last_login_time;" json:"-"`                                                     //最后登陆时间
	LastLoginIp   string        `gorm:"column:last_login_ip;not null;default:'';type:varchar(20)" json:"-"`                   //最后登录IP
	LoginTimes    uint          `gorm:"column:login_times;not null;default:0;type:int(10)" json:"-"`                          //登录次数
	InviteCode    string        `gorm:"column:invite_code;not null;default:'';type:varchar(100)" json:"-"`                    //邀请码
}

func (Account) TableName() string {
	return "im_accounts"
}

// 邀请码表
type InviteCode struct {
	Model
	Code   string           `gorm:"column:code;not null;default:'';type:varchar(50);index:code" json:"code"` //邀请码
	Times  uint             `gorm:"column:times;not null;default:0" json:"times"`                            //使用次数
	Status InviteCodeStatus `gorm:"column:status;not null;default:'';type:varchar(20)" json:"status"`        //状态
	Intro  string           `gorm:"column:intro;not null;default:'';type:varchar(1000)" json:"intro"`        //备注
}

func (InviteCode) TableName() string {
	return "im_invite_codes"
}

const (
	LoginExpired = 7 * 24 * time.Hour //登陆Token有效期
)

type InviteCodeStatus AccountStatus //邀请码状态
type AccountStatus string           //用户状态

var (
	AccountStatusNormal AccountStatus = "normal" //状态正常
	AccountStatusLock   AccountStatus = "lock"   //账户被封禁
)

type OnlineStatus string // 在线状态

var (
	OnlineStatusOnline  OnlineStatus = "online"  //在线
	OnlineStatusOffline OnlineStatus = "offline" //离线
)

type CaptchaType string

var (
	CaptchaTypeLogin    CaptchaType = "login"
	CaptchaTypeRegister CaptchaType = "register"
	CaptchaTypeEmail    CaptchaType = "email"
	CaptchaTypeMobile   CaptchaType = "mobile"
)

type Platform string

var (
	PlatformWeb     Platform = "web"
	PlatformPC      Platform = "pc"
	PlatformH5      Platform = "h5"
	PlatformIOS     Platform = "ios"
	PlatformAndroid Platform = "android"
	PlatformUnknown Platform = "unknown"
)

type BasicCaptchaReq struct {
	Type CaptchaType `form:"type" binding:"required"` //验证码类型
}

type BasicCaptchaRes struct {
	Key     string `json:"key"`
	Captcha []byte `json:"captcha"`
}

type EmailCaptchaReq struct {
	Email string `form:"email" binding:"required"`
}

type EmailLoginReq struct {
	Email      string   `form:"email" binding:"required"`
	Captcha    string   `form:"captcha" binding:"required"`
	Platform   Platform `form:"platform" binding:"required"`
	InviteCode string   `form:"inviteCode"`
}

type EmailResetPasswordReq struct {
	Email    string `form:"email" binding:"required"`
	Captcha  string `form:"captcha" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type MobileCaptchaReq struct {
	Mobile string `form:"mobile" binding:"required"`
}

type MobileLoginReq struct {
	Mobile     string   `form:"mobile" binding:"required"`
	Captcha    string   `form:"captcha" binding:"required"`
	Platform   Platform `form:"platform" binding:"required"`
	InviteCode string   `form:"inviteCode"`
}

type MobileResetPasswordReq struct {
	Mobile   string `form:"mobile" binding:"required"`
	Captcha  string `form:"captcha" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type LoginOrRegisterReq struct {
	Account    string   `form:"account" binding:"required"`
	Password   string   `form:"password" binding:"required,min=6,max=50"`
	CaptchaKey string   `form:"captchaKey"`
	Captcha    string   `form:"captcha"`
	Platform   Platform `form:"platform" binding:"required"`
	InviteCode string   `form:"inviteCode"`
}

type LoginOrRegisterRes struct {
	Token              string `json:"token"`
	NeedUpdatePassword bool   `json:"needUpdatePassword"`
}

type InfoRes struct {
	Account *Account `json:"account"`
}

type UpdatePasswordReq struct {
	Password string `form:"password" binding:"required"`
}

type UpdateAccountInfoReq struct {
	Avatar    string  `form:"avatar"`
	Nickname  string  `form:"nickname"`
	Intro     string  `form:"intro"`
	Gender    string  `form:"gender"`
	Longitude float64 `form:"longitude"`
	Latitude  float64 `form:"latitude"`
	Country   string  `form:"country"`
	Province  string  `form:"province"`
	City      string  `form:"city"`
	District  string  `form:"district"`
}
