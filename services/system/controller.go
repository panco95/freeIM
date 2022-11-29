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

func (ctrl *GinController) GetQiniuUploadToken(c *gin.Context) {
	uploadToken := ctrl.SystemSvc.GetQiniuUploadToken()
	c.JSON(http.StatusOK, &resp.Response{Result: gin.H{
		"uploadToken": uploadToken,
	}})

}
