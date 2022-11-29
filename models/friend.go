package models

import "time"

// 好友表
type Friend struct {
	Model
	AccountId     uint           `gorm:"column:account_id;not null;default:0;index:account_id;" json:"-"` //账号id
	FriendId      uint           `gorm:"column:friend_id;not null;default:0;index:friend_id" json:"-"`    //好友id
	Account       *Account       `gorm:"foreignKey:FriendId" json:"account"`
	FriendsGroups []*FriendGroup `gorm:"many2many:friend_group_friends;foreignKey:FriendId;joinForeignKey:FriendId;" json:"friendGroups"` //好友在哪些分组
	Remark        string         `gorm:"column:remark;not null;default:'';type:varchar(100);" json:"remark"`                              //好友备注
	Label         string         `gorm:"column:label;not null;default:'';type:varchar(1000);" json:"label"`                               //好友自定义字段(标签)
	Blacklist     bool           `gorm:"column:blacklist;not null" json:"-"`                                                              //是否黑名单
}

// 好友分组表
type FriendGroup struct {
	Model
	AccountId uint      `gorm:"column:account_id;not null;default:0;index:account_id;" json:"-"`                           //账号id
	Friends   []*Friend `gorm:"many2many:friend_group_friends;References:FriendId;joinReferences:FriendId" json:"friends"` //分组下的好友
	Name      string    `gorm:"column:name;not null;default:'';" json:"name"`                                              //分组名称
	Order     uint      `gorm:"column:order;not null;default:0;" json:"-"`                                                 //排序
}

type ApplyStatus string

var (
	ApplyStatusWait ApplyStatus = "wait" //等待验证
	ApplyStatusPass ApplyStatus = "pass" //已通过
	ApplyStatusDeny ApplyStatus = "deny" //已拒绝
)

// 好友申请表
type FriendApply struct {
	Model
	FromAccountId uint        `gorm:"column:from_account_id;not null;default:0;index:from_account_id;" json:"-"` //发起请求账号id
	FromAccount   *Account    `gorm:"foreignKey:FromAccountId" json:"fromAccount,omitempty"`
	ToAccountId   uint        `gorm:"column:to_account_id;not null;default:0;index:to_account_id;" json:"-"` //接受i请求账号id
	ToAccount     *Account    `gorm:"foreignKey:ToAccountId" json:"toAccount,omitempty"`
	ApplyReason   string      `gorm:"column:apply_reason;not null;default:'';type:varchar(100);" json:"applyReason"` //申请理由
	DenyReason    string      `gorm:"column:deny_reason;not null;default:'';type:varchar(100);" json:"denyReason"`   //拒绝原因
	Status        ApplyStatus `gorm:"column:status;not null;default:'wait';type:varchar(10);" json:"status"`         //申请状态
	Replytime     *time.Time  `gorm:"column:reply_time;" json:"replyTime"`                                           //回复时间
}

type SearchFriendReq struct {
	Account string `form:"account" binding:"required"`
}

type AddFriendReq struct {
	ToIDReq
	Reason string `form:"reason"`
}

type AddFriendReplyReq struct {
	ToIDReq
	Status ApplyStatus `form:"status" binding:"required,oneof=pass deny"`
	Reason string      `form:"reason"`
}

type FriendListReq struct {
	BlackList uint `form:"blacklist"`
}

type ToIDReq struct {
	ToID uint `form:"toId" binding:"required"`
}

type SetFriendRemarkReq struct {
	ToIDReq
	Remark string `form:"remark" binding:"required"`
}

type SetFriendLabelReq struct {
	ToIDReq
	Label string `form:"label" binding:"required"`
}

type CreateFriendGroupReq struct {
	Name    string `form:"name" binding:"required"`
	Members []uint `form:"members"`
}

type FriendGroupReq struct {
	GroupId uint `form:"groupId" binding:"required"`
}

type FriendGroupMembersReq struct {
	GroupId uint   `form:"groupId" binding:"required"`
	Members []uint `form:"members"`
}

type RenameFriendGroupReq struct {
	GroupId uint   `form:"groupId" binding:"required"`
	Name    string `form:"name" binding:"required,max=50"`
}

type VerifyFriendRes struct {
	IsFriend    bool `json:"isFriend"`
	IsBlacklist bool `json:"isBlacklist"`
}
