package system

type GinController struct {
	SystemSvc *Service
}

// NewGinController ...
func NewGinController(svc *Service) *GinController {
	return &GinController{
		SystemSvc: svc,
	}
}
