package chat

import (
	"im/services/system"

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

func (s *GinController) ConnectWebsocket(ctx *gin.Context) {
	s.ChatSvc.ConnectWebsocket(ctx)
}
