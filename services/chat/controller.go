package chat

import (
	"im/models"
	"im/pkg/gin/middlewares"
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

// 已读消息
func (ctrl *GinController) ReadMessage(c *gin.Context) {
	req := &models.ToIDReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.ChatSvc.ReadMessage(
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

// 获取消息记录
func (ctrl *GinController) GetMessagLogs(c *gin.Context) {
	req := &models.GetMessagesReq{}
	p, _ := c.Get("pagination")
	page := p.(*middlewares.Pagination)

	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	items, total, err := ctrl.ChatSvc.GetMessagLogs(
		c.Request.Context(),
		c.GetUint("id"),
		req,
		page,
	)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	result := resp.PageResult{
		Items: items,
		Total: total,
	}

	c.JSON(http.StatusOK, &resp.Response{Result: result})
}
