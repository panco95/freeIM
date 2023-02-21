package system

import (
	"im/pkg/resp"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinController struct {
	SystemSvc *Service
}

// NewGinController ...
func NewGinController(svc *Service) *GinController {
	return &GinController{
		SystemSvc: svc,
	}
}

func (ctrl *GinController) GetConfigs(c *gin.Context) {
	configs := ctrl.SystemSvc.config.GetAll()
	uploadToken, accessDomain := ctrl.SystemSvc.GetQiniuUploadToken()
	c.JSON(http.StatusOK, &resp.Response{Result: gin.H{
		"qiniu": gin.H{
			"uploadToken":  uploadToken,
			"accessDomain": accessDomain,
		},
		"configs": configs,
	}})
}
