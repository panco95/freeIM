package account

import (
	"net/http"

	"im/models"
	"im/pkg/resp"
	"im/services/system"

	"github.com/gin-gonic/gin"
)

type GinController struct {
	AccountSvc *Service
	SystemSvc  *system.Service
}

func NewGinController(accountSvc *Service, systemSvc *system.Service) *GinController {
	return &GinController{
		AccountSvc: accountSvc,
		SystemSvc:  systemSvc,
	}
}

// 获取图形验证码
func (ctrl *GinController) ImageCaptcha(c *gin.Context) {
	req := &models.BasicCaptchaReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	result, err := ctrl.AccountSvc.GetImageCaptcha(c.Request.Context(), req.Type)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, &resp.Response{Result: result})
}

// 账号注册
func (ctrl *GinController) BasicRegister(c *gin.Context) {
	req := &models.LoginOrRegisterReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	token, _, err := ctrl.AccountSvc.BasicRegister(c.Request.Context(), req, c.ClientIP())
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	result := &models.LoginOrRegisterRes{
		Token: token,
	}
	c.JSON(http.StatusOK, &resp.Response{Result: result, Message: resp.REGISTER_SUCCESS})
}

// 账号登录
func (ctrl *GinController) BasicLogin(c *gin.Context) {
	req := &models.LoginOrRegisterReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	token, _, err := ctrl.AccountSvc.BasicLogin(c.Request.Context(), req, c.ClientIP())
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	result := &models.LoginOrRegisterRes{
		Token: token,
	}
	c.JSON(http.StatusOK, &resp.Response{Result: result, Message: resp.LOGIN_SUCCESS})
}

// 获取邮箱验证码
func (ctrl *GinController) EmailCaptcha(c *gin.Context) {
	req := &models.EmailCaptchaReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.AccountSvc.SendCaptcha(c.Request.Context(), models.CaptchaTypeEmail, req.Email)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, &resp.Response{Message: resp.SEND_SUCCESS})
}

// 邮箱登录
func (ctrl *GinController) EmailLogin(c *gin.Context) {
	req := &models.EmailLoginReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	token, needUpdatePassword, err := ctrl.AccountSvc.EmailOrMobileLogin(
		c.Request.Context(),
		models.CaptchaTypeEmail,
		req.Captcha,
		&models.Account{
			Email: req.Email,
		},
		c.GetString("ip"),
	)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	result := &models.LoginOrRegisterRes{
		Token:              token,
		NeedUpdatePassword: needUpdatePassword,
	}
	c.JSON(http.StatusOK, &resp.Response{Result: result, Message: resp.LOGIN_SUCCESS})
}

// 通过邮箱重置密码
func (ctrl *GinController) EmailResetPassword(c *gin.Context) {
	req := &models.EmailResetPasswordReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.AccountSvc.EmailOrMobileResetPassword(c.Request.Context(), models.CaptchaTypeEmail, req.Captcha, &models.Account{
		Email: req.Email,
	}, req.Password)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, &resp.Response{Message: resp.SETTING_SUCCESS})
}

// 获取手机验证码
func (ctrl *GinController) MobileCaptcha(c *gin.Context) {
	req := &models.MobileCaptchaReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.AccountSvc.SendCaptcha(c.Request.Context(), models.CaptchaTypeMobile, req.Mobile)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, &resp.Response{Message: resp.SEND_SUCCESS})
}

// 手机号登录
func (ctrl *GinController) MobileLogin(c *gin.Context) {
	req := &models.MobileLoginReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	token, needUpdatePassword, err := ctrl.AccountSvc.EmailOrMobileLogin(
		c.Request.Context(),
		models.CaptchaTypeMobile,
		req.Captcha,
		&models.Account{
			Mobile: req.Mobile,
		},
		c.GetString("ip"),
	)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	result := &models.LoginOrRegisterRes{
		Token:              token,
		NeedUpdatePassword: needUpdatePassword,
	}
	c.JSON(http.StatusOK, &resp.Response{Result: result, Message: resp.LOGIN_SUCCESS})
}

// 通过邮箱重置密码
func (ctrl *GinController) MobileResetPassword(c *gin.Context) {
	req := &models.MobileResetPasswordReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.AccountSvc.EmailOrMobileResetPassword(c.Request.Context(), models.CaptchaTypeMobile, req.Captcha, &models.Account{
		Mobile: req.Mobile,
	}, req.Password)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, &resp.Response{Message: resp.SETTING_SUCCESS})
}

// 查询当前登录账号信息
func (ctrl *GinController) Info(c *gin.Context) {
	result, err := ctrl.AccountSvc.Info(
		c.Request.Context(),
		c.GetUint("id"),
	)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, &resp.Response{Result: result})
}

// 设置密码
func (ctrl *GinController) UpdatePassword(c *gin.Context) {
	req := &models.UpdatePasswordReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.AccountSvc.UpdatePassword(
		c.Request.Context(),
		c.GetUint("id"),
		req.Password,
	)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, &resp.Response{Message: resp.SETTING_SUCCESS})
}
