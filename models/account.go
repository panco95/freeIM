package models

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Account struct {
	Model
	Username      string        `gorm:"column:username;not null;default:'';type:varchar(50);index:username" json:"username"` //用户名
	Email         string        `gorm:"column:email;not null;default:'';type:varchar(100)" json:"email"`                     //邮箱地址
	Mobile        string        `gorm:"column:mobile;not null;default:'';type:varchar(50)" json:"mobile"`                    //手机号
	Nikcname      string        `gorm:"column:nickname;not null;default:'';type:varchar(50)" json:"nickname"`                //昵称
	Avatar        string        `gorm:"column:avatar;not null;default:'';type:varchar(500)" json:"avatar"`                   //头像
	Gender        string        `gorm:"column:gender;not null;default:'';type:varchar(10)" json:"gender"`                    //性别
	Birth         string        `gorm:"column:birth;default:null;type:varchar(100)" json:"birth"`                            //生日
	Age           int           `gorm:"-" json:"age"`                                                                        //年龄
	Intro         string        `gorm:"column:intro;not null;default:'';type:varchar(1000)" json:"intro"`                    //个人介绍
	Password      string        `gorm:"column:password;not null;default:'';type:varchar(200)" json:"-"`                      //密码
	PasswordSalt  string        `gorm:"column:password_salt;not null;default:'';type:varchar(200)" json:"-"`                 //密码盐值
	Status        AccountStatus `gorm:"column:status;not null;default:'normal';type:varchar(20)" json:"-"`                   //状态
	LastLoginTime *time.Time    `gorm:"column:last_login_time;" json:"-"`                                                    //最后登陆时间
	LastLoginIp   string        `gorm:"column:last_login_ip;not null;default:'';type:varchar(20)" json:"-"`                  //最后登录IP
	LoginTimes    uint          `gorm:"column:login_times;not null;default:0;type:int(10)" json:"-"`                         //登录次数
}

const (
	LoginExpired = 7 * 24 * time.Hour
)

type AccountStatus string

var (
	AccountStatusNormal AccountStatus = "normal"
	AccountStatusLock   AccountStatus = "lock"
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
	Email    string   `form:"email" binding:"required"`
	Captcha  string   `form:"captcha" binding:"required"`
	Platform Platform `form:"platform" binding:"required"`
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
	Mobile   string   `form:"mobile" binding:"required"`
	Captcha  string   `form:"captcha" binding:"required"`
	Platform Platform `form:"platform" binding:"required"`
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

func (a *Account) Query(
	ctx context.Context,
	db *gorm.DB,
	account *Account,
) error {
	result := db.
		Where("(id <> 0 AND id = ?) OR (username <> '' AND username = ?) OR (email <> '' AND email = ?) OR (mobile <> '' AND mobile = ?)",
			account.ID,
			account.Username,
			account.Email,
			account.Mobile,
		).
		First(a)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
