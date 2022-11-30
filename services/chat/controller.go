package chat

import (
	"im/models"
	"im/pkg/resp"
	"im/services/system"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinController struct {
	ChatSvc   *Service
	SystemSvc *system.Service
}

func NewGinController(chatSvc *Service, systemSvc *system.Service) *GinController {
	return &GinController{
		ChatSvc:   chatSvc,
		SystemSvc: systemSvc,
	}
}

func (ctrl *GinController) ConnectWebsocket(ctx *gin.Context) {
	ctrl.ChatSvc.ConnectWebsocket(ctx)
}

// 发送消息
func (ctrl *GinController) SendMessage(c *gin.Context) {
	req := &models.SendMessageReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	res, err := ctrl.ChatSvc.SendMessage(
		c.Request.Context(),
		c.GetUint("id"),
		req,
	)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, &resp.Response{
		Result:  res,
		Message: resp.SEND_SUCCESS,
	})
}

// 撤回消息
func (ctrl *GinController) RevocationMessage(c *gin.Context) {
	req := &models.RevocationMessageReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.ChatSvc.RevocationMessage(
		c.Request.Context(),
		c.GetUint("id"),
		req,
	)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, &resp.Response{Message: resp.SUCCESS})
}
