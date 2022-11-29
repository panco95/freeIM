package chatgroup

import (
	"im/models"
	"im/pkg/resp"
	"im/services/system"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinController struct {
	ChatGroupSvc *Service
	SystemSvc    *system.Service
}

func NewGinController(chatGroupSvc *Service, systemSvc *system.Service) *GinController {
	return &GinController{
		ChatGroupSvc: chatGroupSvc,
		SystemSvc:    systemSvc,
	}
}

// 搜索群聊
func (ctrl *GinController) SearchChatGroup(c *gin.Context) {
	req := &models.SearchChatGroupReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	items, err := ctrl.ChatGroupSvc.SearchChatGroup(
		c.Request.Context(),
		c.GetUint("id"),
		req,
	)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	result := resp.PageResult{
		Items: items,
		Total: int64(len(items)),
	}

	c.JSON(http.StatusOK, &resp.Response{Result: result})
}

// 创建群聊
func (ctrl *GinController) CreateChatGroup(c *gin.Context) {
	req := &models.CreateChatGroupReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.ChatGroupSvc.CreateChatGroup(
		c.Request.Context(),
		c.GetUint("id"),
		req,
	)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, &resp.Response{Message: resp.CREATE_SUCCESS})
}

// 修改群资料
func (ctrl *GinController) EditChatGroup(c *gin.Context) {
	req := &models.EditChatGroupReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.ChatGroupSvc.EditChatGroup(
		c.Request.Context(),
		c.GetUint("id"),
		req,
	)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, &resp.Response{Message: resp.EDIT_SUCCESS})
}

// 加群申请
func (ctrl *GinController) JoinChatGroup(c *gin.Context) {
	req := &models.JoinChatGroupReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.ChatGroupSvc.JoinChatGroup(
		c.Request.Context(),
		c.GetUint("id"),
		req,
	)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, &resp.Response{Message: resp.APPLY_SUCCESS})
}

// 加群审批
func (ctrl *GinController) JoinChatGroupReply(c *gin.Context) {
	req := &models.JoinChatGroupReplyReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.ChatGroupSvc.JoinChatGroupReply(
		c.Request.Context(),
		c.GetUint("id"),
		req,
	)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, &resp.Response{Message: resp.PROCESS_SUCCESS})
}

// 加群审批列表
func (ctrl *GinController) JoinChatGroupList(c *gin.Context) {
	items, err := ctrl.ChatGroupSvc.JoinChatGroupList(
		c.Request.Context(),
		c.GetUint("id"),
	)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	result := resp.PageResult{
		Items: items,
		Total: int64(len(items)),
	}

	c.JSON(http.StatusOK, &resp.Response{Result: result})
}

// 我的群聊列表
func (ctrl *GinController) ChatGroupList(c *gin.Context) {
	items, err := ctrl.ChatGroupSvc.ChatGroupList(
		c.Request.Context(),
		c.GetUint("id"),
	)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	result := resp.PageResult{
		Items: items,
		Total: int64(len(items)),
	}

	c.JSON(http.StatusOK, &resp.Response{Result: result})
}

// 群聊信息(包括成员列表)
func (ctrl *GinController) ChatGroupInfo(c *gin.Context) {
	req := &models.GroupIdReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	info, err := ctrl.ChatGroupSvc.ChatGroupInfo(
		c.Request.Context(),
		c.GetUint("id"),
		req,
	)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, &resp.Response{Result: info})
}

// 退出群聊
func (ctrl *GinController) ExitChatGroup(c *gin.Context) {
	req := &models.GroupIdReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.ChatGroupSvc.ExitChatGroup(
		c.Request.Context(),
		c.GetUint("id"),
		req,
	)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, &resp.Response{Message: resp.EXIT_SUCCESS})
}

// 转让群聊
func (ctrl *GinController) TransferChatGroup(c *gin.Context) {
	req := &models.ChatGroupToIDReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.ChatGroupSvc.TransferChatGroup(
		c.Request.Context(),
		c.GetUint("id"),
		req,
	)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, &resp.Response{Message: resp.TRANSFER_SUCCESS})
}

// 解散群聊
func (ctrl *GinController) DissolveChatGroup(c *gin.Context) {
	req := &models.GroupIdReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.ChatGroupSvc.DissolveChatGroup(
		c.Request.Context(),
		c.GetUint("id"),
		req,
	)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, &resp.Response{Message: resp.PROCESS_SUCCESS})
}

// 群聊踢出成员
func (ctrl *GinController) ChatGroupKickMember(c *gin.Context) {
	req := &models.ChatGroupToIDReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.ChatGroupSvc.ChatGroupKickMember(
		c.Request.Context(),
		c.GetUint("id"),
		req,
	)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, &resp.Response{Message: resp.PROCESS_SUCCESS})
}

// 设置群聊管理员
func (ctrl *GinController) ChatGroupSetManager(c *gin.Context) {
	req := &models.ChatGroupSetManagerReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.ChatGroupSvc.ChatGroupSetManager(
		c.Request.Context(),
		c.GetUint("id"),
		req,
	)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, &resp.Response{Message: resp.SETTING_SUCCESS})
}

// 群聊成员禁言
func (ctrl *GinController) ChatGroupBannedMember(c *gin.Context) {
	req := &models.ChatGroupBannedMemberReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.ChatGroupSvc.ChatGroupBannedMember(
		c.Request.Context(),
		c.GetUint("id"),
		req,
	)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, &resp.Response{Message: resp.PROCESS_SUCCESS})
}
