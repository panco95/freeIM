package account

import (
	"bytes"
	"context"
	"errors"
	"image/png"
	"time"

	"im/models"
	"im/pkg/database"
	"im/pkg/email"
	"im/pkg/jwt"
	"im/pkg/resp"
	"im/pkg/sms"
	"im/pkg/utils"
	"im/services/system/config"

	"github.com/afocus/captcha"
	redisCache "github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	log         *zap.SugaredLogger
	mysqlClient *database.Client
	redisClient *redis.Client
	cacheClient *redisCache.Cache
	emailClient *email.Mail
	smsClient   sms.Sms
	jwt         *jwt.Jwt
	config      *config.Config
}

func NewService(
	mysqlClient *database.Client,
	redisClient *redis.Client,
	cacheClient *redisCache.Cache,
	emailClient *email.Mail,
	smsClient sms.Sms,
	jwt *jwt.Jwt,
	config *config.Config,
) *Service {
	return &Service{
		log:         zap.S().With("module", "services.account.service"),
		mysqlClient: mysqlClient,
		redisClient: redisClient,
		cacheClient: cacheClient,
		emailClient: emailClient,
		smsClient:   smsClient,
		jwt:         jwt,
		config:      config,
	}
}

// 获取图片验证码
func (s *Service) GetImageCaptcha(
	ctx context.Context,
	captchaType models.CaptchaType,
) (*models.BasicCaptchaRes, error) {
	key := uuid.New().String()
	cap := captcha.New()
	err := cap.AddFontFromBytes(utils.GetDefaultFont())
	if err != nil {
		s.log.Errorf("GetCaptcha cap.AddFontFromBytes err=%v", err)
		return nil, err
	}
	img, code := cap.Create(4, captcha.NUM)
	content := bytes.NewBuffer([]byte{})
	err = png.Encode(content, img)
	if err != nil {
		s.log.Errorf("GetCaptcha png.Encode err=%v", err)
		return nil, err
	}

	cacheKey := "captcha:" + string(captchaType) + ":" + key
	err = s.cacheClient.Set(&redisCache.Item{
		Ctx:   ctx,
		Key:   cacheKey,
		Value: code,
		TTL:   time.Minute * 5,
	})
	if err != nil {
		s.log.Errorf("GetCaptcha cacheClient.Set err=%v", err)
		return nil, err
	}

	return &models.BasicCaptchaRes{
		Key:     key,
		Captcha: content.Bytes(),
	}, nil
}

// 校验验证码
func (s *Service) CheckCaptcha(
	ctx context.Context,
	captchaKey string,
	captcha string,
	captchaType models.CaptchaType,
) error {
	captchaCache := ""
	cacheKey := "captcha:" + string(captchaType) + ":" + captchaKey
	err := s.cacheClient.Get(ctx, cacheKey, &captchaCache)
	if err != nil {
		if err == redisCache.ErrCacheMiss {
			return errors.New(resp.CAPTCHA_EXPIRED)
		}
		return err
	}
	if captchaCache == "" {
		return errors.New(resp.CAPTCHA_EXPIRED)
	}
	if captcha != captchaCache {
		return errors.New(resp.CAPTCHA_ERROR)
	}

	go func() {
		err := s.cacheClient.Delete(ctx, cacheKey)
		if err != nil {
			s.log.Errorf("CheckCaptcha cacheClient.Delete %v", err)
		}
	}()

	return nil
}

// 账号登录
func (s *Service) BasicLogin(
	ctx context.Context,
	req *models.LoginOrRegisterReq,
	ip string,
) (string, *models.Account, error) {
	if s.config.Get("login_captcha") == "true" {
		err := s.CheckCaptcha(ctx, req.CaptchaKey, req.Captcha, models.CaptchaTypeLogin)
		if err != nil {
			return "", nil, err
		}
	}

	account := &models.Account{}
	err := account.Query(ctx, s.mysqlClient.Db(), &models.Account{
		Username: req.Account,
		Email:    req.Account,
		Mobile:   req.Account,
	})
	if err != nil {
		return "", nil, err
	}
	if account.ID == 0 {
		return "", nil, errors.New(resp.ACCOUNT_NOT_FOUND)
	}
	if account.Status == models.AccountStatusLock {
		return "", nil, errors.New(resp.ACCOUNT_LOCKED)
	}
	if utils.Md5(utils.Md5(req.Password)+account.PasswordSalt) != account.Password {
		return "", nil, errors.New(resp.ACCOUNT_PWD_ERROR)
	}

	token, err := s.jwt.BuildToken(
		account.ID,
		models.LoginExpired,
	)
	if err != nil {
		s.log.Errorf("Login jwt.BuildToken %v", err)
		return "", nil, errors.New(resp.SERVER_ERROR)
	}

	go func() {
		now := time.Now()
		err = s.mysqlClient.Db().
			Model(&account).
			Updates(models.Account{
				LastLoginTime: &now,
				LastLoginIp:   ip,
				LoginTimes:    account.LoginTimes + 1,
			}).
			Error
		if err != nil {
			s.log.Errorf("Login update account %v", err)
		}
	}()

	return token, account, nil
}

// 账号注册
func (s *Service) BasicRegister(
	ctx context.Context,
	req *models.LoginOrRegisterReq,
	ip string,
) (string, *models.Account, error) {
	if s.config.Get("login_captcha") == "true" {
		err := s.CheckCaptcha(ctx, req.CaptchaKey, req.Captcha, models.CaptchaTypeRegister)
		if err != nil {
			return "", nil, err
		}
	}
	if utils.IsChinese(req.Account) {
		return "", nil, errors.New(resp.ACCOUNT_HAS_CHINESE)
	}

	exists := &models.Account{}
	err := s.mysqlClient.Db().
		Model(&models.Account{}).
		Where("username = ?", req.Account).
		First(exists).Error
	if err != nil {
		return "", nil, err
	}
	if exists.ID != 0 {
		return "", nil, errors.New(resp.ACCOUNT_EXISTS)
	}

	now := time.Now()
	account := &models.Account{}
	account.Username = req.Account
	account.Nikcname = req.Account
	account.LastLoginTime = &now
	account.LoginTimes = 1
	account.PasswordSalt = utils.RandStr(6)
	account.Password = utils.Md5(utils.Md5(req.Password) + account.PasswordSalt)
	account.LastLoginIp = ip

	err = s.mysqlClient.Db().
		Model(&models.Account{}).
		Create(account).Error
	if err != nil {
		return "", nil, err
	}

	token, err := s.jwt.BuildToken(
		account.ID,
		models.LoginExpired,
	)
	if err != nil {
		s.log.Errorf("Register jwt.BuildToken %v", err)
		return "", nil, errors.New(resp.SERVER_ERROR)
	}

	return token, account, nil
}

// 发送验证码
func (s *Service) SendCaptcha(
	ctx context.Context,
	captchaType models.CaptchaType,
	captchaKey string,
) error {
	vcode := utils.RandNumber(5)
	cacheKey := "captcha:" + string(captchaType) + ":" + captchaKey
	err := s.cacheClient.Set(&redisCache.Item{
		Ctx:   ctx,
		Key:   cacheKey,
		Value: vcode,
		TTL:   time.Minute * 5,
	})
	if err != nil {
		s.log.Errorf("SendCaptcha cacheClient.Set %v", err)
		return err
	}

	switch captchaType {
	case models.CaptchaTypeEmail:
		err = s.emailClient.SendEmail("IM邮箱登录验证码", "您的验证码为："+vcode, []string{captchaKey})
	case models.CaptchaTypeMobile:
		err = s.smsClient.Send(ctx, captchaKey, "【IM】您的验证码为："+vcode)
	}
	if err != nil {
		s.log.Errorf("SendCaptcha send %v", err)
	}

	return err
}

// 邮箱或手机号登录
func (s *Service) EmailOrMobileLogin(
	ctx context.Context,
	captchaType models.CaptchaType,
	captcha string,
	account *models.Account,
	ip string,
) (string, bool, error) {
	captchatKey := ""
	switch captchaType {
	case models.CaptchaTypeEmail:
		captchatKey = account.Email
		account.Nikcname = account.Email
	case models.CaptchaTypeMobile:
		captchatKey = account.Mobile
		account.Nikcname = account.Mobile
	}
	err := s.CheckCaptcha(ctx, captchatKey, captcha, captchaType)
	if err != nil {
		return "", false, err
	}

	err = account.Query(ctx, s.mysqlClient.Db(), account)
	if err != nil {
		return "", false, err
	}
	if account.ID == 0 {
		account.LastLoginIp = ip
		return s.AutoRegister(ctx, account)
	}
	if account.Status == models.AccountStatusLock {
		return "", false, errors.New(resp.ACCOUNT_LOCKED)
	}

	token, err := s.jwt.BuildToken(
		account.ID,
		models.LoginExpired,
	)
	if err != nil {
		s.log.Errorf("Login jwt.BuildToken %v", err)
		return "", false, errors.New(resp.SERVER_ERROR)
	}

	go func() {
		now := time.Now()
		err = s.mysqlClient.Db().
			Model(&models.Account{}).
			Where("id = ?", account.ID).
			Updates(models.Account{
				LastLoginTime: &now,
				LastLoginIp:   ip,
				LoginTimes:    account.LoginTimes + 1,
			}).
			Error
		if err != nil {
			s.log.Errorf("EmailLogin update account %v", err)
		}
	}()

	needSetPassword := false
	if account.Password == "" {
		needSetPassword = true
	}

	return token, needSetPassword, nil
}

// 自动注册
func (s *Service) AutoRegister(
	ctx context.Context,
	account *models.Account,
) (string, bool, error) {
	now := time.Now()
	account.LastLoginTime = &now
	account.LoginTimes = 1
	account.Username = account.Email + account.Mobile

	err := s.mysqlClient.Db().
		Model(&models.Account{}).
		Create(account).Error
	if err != nil {
		return "", false, err
	}

	token, err := s.jwt.BuildToken(
		account.ID,
		models.LoginExpired,
	)
	if err != nil {
		s.log.Errorf("Register jwt.BuildToken %v", err)
		return "", false, errors.New(resp.SERVER_ERROR)
	}

	return token, true, nil
}

// 邮箱或手机号重置密码
func (s *Service) EmailOrMobileResetPassword(
	ctx context.Context,
	captchaType models.CaptchaType,
	captcha string,
	account *models.Account,
	password string,
) error {
	captchatKey := ""
	switch captchaType {
	case models.CaptchaTypeEmail:
		captchatKey = account.Email
	case models.CaptchaTypeMobile:
		captchatKey = account.Mobile
	}
	err := s.CheckCaptcha(ctx, captchatKey, captcha, captchaType)
	if err != nil {
		return err
	}

	err = account.Query(ctx, s.mysqlClient.Db(), account)
	if err != nil {
		return err
	}
	if account.ID == 0 {
		return errors.New(resp.ACCOUNT_NOT_FOUND)
	}

	account.Password = password
	account.PasswordSalt = utils.RandStr(6)
	account.Password = utils.Md5(utils.Md5(account.Password) + account.PasswordSalt)
	err = s.mysqlClient.Db().
		Model(&models.Account{}).
		Where("id = ?", account.ID).
		Updates(models.Account{
			Password:     account.Password,
			PasswordSalt: account.PasswordSalt,
		}).
		Error
	if err != nil {
		s.log.Errorf("EmailOrMobileResetPassword update %v", err)
	}

	return nil
}

// 账号信息
func (s *Service) Info(
	ctx context.Context,
	accountId uint,
) (*models.InfoRes, error) {
	account := &models.Account{}
	err := account.Query(ctx, s.mysqlClient.Db(), &models.Account{
		Model: models.Model{ID: accountId},
	})
	if err != nil {
		return nil, err
	}

	account.CreatedAt = nil
	account.UpdatedAt = nil
	result := &models.InfoRes{
		Account: account,
	}

	return result, nil
}

// 设置密码
func (s *Service) UpdatePassword(
	ctx context.Context,
	accountId uint,
	password string,
) error {
	newPasswordSalt := utils.RandStr(6)
	newPassword := utils.Md5(utils.Md5(password) + newPasswordSalt)
	err := s.mysqlClient.Db().
		Model(&models.Account{}).
		Where("id = ?", accountId).
		Updates(models.Account{
			Password:     newPassword,
			PasswordSalt: newPasswordSalt,
		}).Error
	if err != nil {
		s.log.Errorf("UpdatePassword %v", err)
		return err
	}

	return nil
}
