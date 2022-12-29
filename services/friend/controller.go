package friend

import (
	"im/models"
	"im/pkg/resp"
	"im/services/system"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinController struct {
	FriendSvc *Service
	SystemSvc *system.Service
}

func NewGinController(friendSvc *Service, systemSvc *system.Service) *GinController {
	return &GinController{
		FriendSvc: friendSvc,
		SystemSvc: systemSvc,
	}
}

// 查找好友
func (ctrl *GinController) SearchFriend(c *gin.Context) {
	req := &models.SearchFriendReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	items, err := ctrl.FriendSvc.SearchFriend(c.Request.Context(), c.GetUint("id"), req)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	res := resp.PageResult{
		Items: items,
		Total: int64(len(items)),
	}

	c.JSON(http.StatusOK, &resp.Response{Result: res})
}

// 添加好友
func (ctrl *GinController) AddFriend(c *gin.Context) {
	req := &models.AddFriendReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.FriendSvc.AddFriend(c.Request.Context(), c.GetUint("id"), req)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, &resp.Response{Message: resp.SUCCESS})
}

// 同意/拒绝好友请求
func (ctrl *GinController) AddFriendReply(c *gin.Context) {
	req := &models.AddFriendReplyReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.FriendSvc.AddFriendReply(c.Request.Context(), c.GetUint("id"), req)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, &resp.Response{Message: resp.SUCCESS})
}

// 好友申请列表
func (ctrl *GinController) FriendApplyList(c *gin.Context) {
	items, err := ctrl.FriendSvc.FriendApplyList(c.Request.Context(), c.GetUint("id"))
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	res := resp.PageResult{
		Items: items,
		Total: int64(len(items)),
	}

	c.JSON(http.StatusOK, &resp.Response{Result: res})
}

// 好友列表
func (ctrl *GinController) FriendList(c *gin.Context) {
	req := &models.FriendListReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	items, err := ctrl.FriendSvc.FriendList(c.Request.Context(), c.GetUint("id"), req)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	res := resp.PageResult{
		Items: items,
		Total: int64(len(items)),
	}

	c.JSON(http.StatusOK, &resp.Response{Result: res})
}

// 单个好友信息
func (ctrl *GinController) FriendInfo(c *gin.Context) {
	req := &models.ToIDReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	result, err := ctrl.FriendSvc.FriendInfo(c.Request.Context(), c.GetUint("id"), req)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, &resp.Response{Result: result})
}

// 删除好友
func (ctrl *GinController) DeleteFriend(c *gin.Context) {
	req := &models.ToIDReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.FriendSvc.DeleteFriend(c.Request.Context(), c.GetUint("id"), req)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, &resp.Response{Message: resp.DELETE_SUCCESS})
}

// 添加黑名单
func (ctrl *GinController) AddBlacklist(c *gin.Context) {
	req := &models.ToIDReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.FriendSvc.AddBlacklist(c.Request.Context(), c.GetUint("id"), req)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, &resp.Response{Message: resp.ADD_SUCCESS})
}

// 移除黑名单
func (ctrl *GinController) DeleteBlacklist(c *gin.Context) {
	req := &models.ToIDReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.FriendSvc.DeleteBlacklist(c.Request.Context(), c.GetUint("id"), req)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, &resp.Response{Message: resp.REMOVE_SUCCESS})
}

// 设置好友备注
func (ctrl *GinController) SetFriendRemark(c *gin.Context) {
	req := &models.SetFriendRemarkReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.FriendSvc.SetFriendRemark(c.Request.Context(), c.GetUint("id"), req)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, &resp.Response{Message: resp.SETTING_SUCCESS})
}

// 设置好友标签（自定义字段）
func (ctrl *GinController) SetFriendLabel(c *gin.Context) {
	req := &models.SetFriendLabelReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.FriendSvc.SetFriendLabel(c.Request.Context(), c.GetUint("id"), req)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, &resp.Response{Message: resp.SETTING_SUCCESS})
}

// 创建好友分组
func (ctrl *GinController) CreateFriendGroup(c *gin.Context) {
	req := &models.CreateFriendGroupReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.FriendSvc.CreateFriendGroup(c.Request.Context(), c.GetUint("id"), req)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, &resp.Response{Message: resp.ADD_SUCCESS})
}

// 删除好友分组
func (ctrl *GinController) DeleteFriendGroup(c *gin.Context) {
	req := &models.FriendGroupReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.FriendSvc.DeleteFriendGroup(c.Request.Context(), c.GetUint("id"), req)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, &resp.Response{Message: resp.DELETE_SUCCESS})
}

// 好友分组添加成员
func (ctrl *GinController) AddFriendGroupMembers(c *gin.Context) {
	req := &models.FriendGroupMembersReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.FriendSvc.OperateFriendGroupMembers(c.Request.Context(), c.GetUint("id"), req, "add")
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, &resp.Response{Message: resp.ADD_SUCCESS})
}

// 好友分组删除成员
func (ctrl *GinController) DelFriendGroupMembers(c *gin.Context) {
	req := &models.FriendGroupMembersReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.FriendSvc.OperateFriendGroupMembers(c.Request.Context(), c.GetUint("id"), req, "del")
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, &resp.Response{Message: resp.DELETE_SUCCESS})
}

// 获取好友分组列表
func (ctrl *GinController) GetFriendGroups(c *gin.Context) {
	items, err := ctrl.FriendSvc.GetFriendGroups(c.Request.Context(), c.GetUint("id"))
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

// 获取指定好友分组
func (ctrl *GinController) GetFriendGroup(c *gin.Context) {
	req := &models.FriendGroupReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	res, err := ctrl.FriendSvc.GetFriendGroup(c.Request.Context(), c.GetUint("id"), req)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, &resp.Response{Result: res})
}

// 重命名好友分组
func (ctrl *GinController) RenameFriendGroup(c *gin.Context) {
	req := &models.RenameFriendGroupReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	err := ctrl.FriendSvc.RenameFriendGroup(c.Request.Context(), c.GetUint("id"), req)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, &resp.Response{Message: resp.SETTING_SUCCESS})
}

// 校验好友（是否你是他的好友或者黑名单）
func (ctrl *GinController) VerifyFriend(c *gin.Context) {
	req := &models.ToIDReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	result, err := ctrl.FriendSvc.VerifyFriend(c.Request.Context(), c.GetUint("id"), req)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, &resp.Response{Result: result})
}

// 附近的人
func (ctrl *GinController) NearFriends(c *gin.Context) {
	req := &models.NearFriendsReq{}
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	items, err := ctrl.FriendSvc.NearFriends(c.Request.Context(), c.GetUint("id"), req)
	if err != nil {
		_ = c.Error(err).
			SetType(gin.ErrorTypePublic)
		return
	}

	res := resp.PageResult{
		Items: items,
		Total: int64(len(items)),
	}

	c.JSON(http.StatusOK, &resp.Response{Result: res})
}
