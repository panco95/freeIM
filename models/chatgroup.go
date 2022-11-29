package models

import "time"

var (
	DefaultChatgroupMaxMembers uint = 500 //默认群成员最大人数
)

// 群聊表
type ChatGroup struct {
	Model
	Name              string             `gorm:"column:name;not null;default:'';type:varchar(50)" json:"name"`                //群名称
	Avatar            string             `gorm:"column:avatar;not null;default:'';type:varchar(500)" json:"avatar"`           //群头像
	Intro             string             `gorm:"column:intro;not null;default:'';type:varchar(1000)" json:"intro"`            //群介绍
	Members           uint               `gorm:"column:members;not null;default:0;" json:"members"`                           //当前成员数
	MembersLimit      uint               `gorm:"column:members_limit;not null;default:0;" json:"members_limit"`               //群内成员数上限(0：无限)
	MembersList       []*ChatGroupMember `gorm:"-" json:"membersList,omitempty"`                                              //成员列表
	SelfInfo          *ChatGroupSelfInfo `gorm:"-" json:"selfInfo,omitempty"`                                                 //自身信息
	DisableAddMember  bool               `gorm:"column:disable_add_member;not null;default:false;" json:"disableAddMember"`   //禁止加成员好友
	DisableViewMember bool               `gorm:"column:disable_view_member;not null;default:false;" json:"disableViewMember"` //禁用查看成员资料
	DisbaleAddGroup   bool               `gorm:"column:disable_add_group;not null;default:false;" json:"disbaleAddGroup"`     //禁用主动申请入群
	EnbaleBeforeMsg   bool               `gorm:"column:enable_before_msg;not null;default:false;" json:"enbaleBeforeMsg"`     //是否开启加群之前的漫游消息
}

// 群聊自身信息
type ChatGroupSelfInfo struct {
	Role     ChatGroupMemberRole `json:"chatGroupMemberRole"`
	IsBanned bool                `json:"isBanned"`
}

type ChatGroupMemberRole string // 群聊成员角色类型
var (
	ChatGroupMemberRoleOwner   ChatGroupMemberRole = "owner"   //群主
	ChatGroupMemberRoleManager ChatGroupMemberRole = "manager" //管理员
	ChatGroupMemberRoleGeneral ChatGroupMemberRole = "general" //普通成员
)

// 群聊成员表
type ChatGroupMember struct {
	Model
	ChatGroupId    uint                `gorm:"column:chat_group_id;not null;default:0;index:group_id" json:"-"` //群组ID
	ChatGroup      *ChatGroup          `gorm:"foreignKey:ChatGroupId" json:"chatGroup,omitempty"`
	AccountId      uint                `gorm:"column:account_id;not null;default:0;index:account_id" json:"-"` //账号ID
	Account        *Account            `gorm:"foreignKey:AccountId" json:"account,omitempty"`
	Role           ChatGroupMemberRole `gorm:"column:role;not null;default:'general';" json:"role"` //成员角色
	Remark         string              `gorm:"column:remark;not null;default:'';" json:"remark"`    //群内称呼
	BannedOverTime *time.Time          `gorm:"column:banned_over_time;" json:"-"`                   //禁言到期时间
	IsBanned       bool                `gorm:"-" json:"isBanned"`                                   //是否禁言中
}

// 加群申请表
type ChatGroupJoin struct {
	Model
	AccountId          uint        `gorm:"column:account_id;not null;default:0;index:account_id;" json:"-"` //发起请求账号id
	Account            *Account    `gorm:"foreignKey:AccountId" json:"account,omitempty"`
	ChatGroupId        uint        `gorm:"column:chat_group_id;not null;default:0;index:chat_group_id;" json:"-"` //接受请求群聊id
	ChatGroup          *ChatGroup  `gorm:"foreignKey:ChatGroupId" json:"chatGroup,omitempty"`
	ManagerAccountList string      `gorm:"column:manager_account_list;not null;default:'';type:varchar(5000)" json:"-"`   //管理员ID列表
	ApplyReason        string      `gorm:"column:apply_reason;not null;default:'';type:varchar(100);" json:"applyReason"` //申请理由
	DenyReason         string      `gorm:"column:deny_reason;not null;default:'';type:varchar(100);" json:"denyReason"`   //拒绝原因
	Status             ApplyStatus `gorm:"column:status;not null;default:'wait';type:varchar(10);" json:"status"`         //申请状态
	Replytime          *time.Time  `gorm:"column:reply_time;" json:"replyTime"`                                           //回复时间
}

type SearchChatGroupReq struct {
	Name string `form:"name" binding:"required"`
}

type CreateChatGroupReq struct {
	Name              string `form:"name" binding:"required"`
	Avatar            string `form:"avatar"`
	Intro             string `form:"intro" binding:"required"`
	DisableAddMember  bool   `form:"disableAddMember"`
	DisableViewMember bool   `form:"disableViewMember"`
	DisbaleAddGroup   bool   `form:"disbaleAddGroup"`
	EnbaleBeforeMsg   bool   `form:"enbaleBeforeMsg"`
}

type EditChatGroupReq struct {
	GroupIdReq
	CreateChatGroupReq
}

type GroupIdReq struct {
	GroupId uint `form:"groupId" binding:"required"`
}

type JoinChatGroupReq struct {
	GroupIdReq
	Reason string `form:"reason"`
}

type JoinChatGroupReplyReq struct {
	GroupIdReq
	AccountId uint        `form:"accountId" binding:"required"`
	Status    ApplyStatus `form:"status" binding:"required,oneof=pass deny"`
	Reason    string      `form:"reason"`
}

type ChatGroupToIDReq struct {
	GroupIdReq
	ToIDReq
}

type ChatGroupSetManagerReq struct {
	GroupIdReq
	ToIDReq
	IsManager bool `form:"isManager"`
}

type ChatGroupBannedMemberReq struct {
	GroupIdReq
	ToIDReq
	Minute uint `form:"minute" binding:"required"`
}
