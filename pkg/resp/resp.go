package resp

const (
	ERROR = 1
)

const (
	SERVER_ERROR           = "服务器错误"
	PARAM_INVALID          = "参数不合法"
	ACCOUNT_NOT_FOUND      = "账号不存在"
	ACCOUNT_PWD_ERROR      = "账号或密码错误"
	ACCOUNT_LOCKED         = "账号被封禁"
	CAPTCHA_ERROR          = "验证码错误"
	CAPTCHA_EXPIRED        = "验证码过期，请刷新后重试"
	TIMEOUT                = "登录超时"
	ACCOUNT_EXISTS         = "账号已存在"
	ACCOUNT_NOT_EXISTS     = "账号不存在"
	ACCOUNT_HAS_CHINESE    = "用户名不能包含中文"
	INVITE_CODE_NOT_EXISTS = "非法邀请码"
	IP_REGISTER_BLAKCLIST  = "您的IP无法注册"
	IP_REGISTER_LIMIT      = "您的IP受到注册限制"
	EMAIL_REGISTER_OFF     = "邮箱注册通道已关闭"
	MOBILE_REGISTER_OFF    = "手机号注册通道已关闭"
	ACCOUNTS_REGISTER_OFF  = "账号注册通道已关闭"

	FRIEND_NOT_EXISTS       = "对方还不是您的好友"
	FRIEND_EXISTS           = "对方已经是您的好友"
	FRIEND_APPLY_EXISTS     = "对方还未验证您的好友请求"
	FRIEND_APPLY_NOT_EXISTS = "请求不存在"
	FRIEND_APPLY_NOT_WAIT   = "此好友请求已处理过"
	BLACKLIST_NOT_EXISTS    = "黑名单不存在"
	FRIEND_GROUP_EXISTS     = "分组已存在"
	FRIEND_GROUP_NOT_EXISTS = "分组不存在"
	FRIEND_ADD_OFF          = "添加好友通道已关闭"

	CHAT_GROUP_NOT_EXISTS        = "群聊不存在"
	CHAT_GROUP_DISABLE_ADD_GROUP = "当前群聊禁止新成员加入"
	CHAT_GROUP_MEMBERS_LIMIT     = "当前群聊人数已满"
	CHAT_GROUP_APPLY_NOT_EXISTS  = "加群申请不存在"
	CHAT_GROUP_APPLY_EXISTS      = "已发送过加群申请"
	CHAT_GROUP_APPLY_NOT_WAIT    = "此加群请求已处理过"
	CHAT_GROUP_IS_MEMBER         = "你已经是此群聊成员"
	CHAT_GROUP_NOT_MANAGER       = "您不是群聊管理员"
	CHAT_GROUP_NOT_OWNER         = "您不是群主"
	CHAT_GROUP_NOT_MEMBER        = "你不是是此群聊成员"
	CHAT_GROUP_OWNER_EXIT        = "群主无法退群，请转让群"
	CHAT_GROUP_NOT_ALLOW         = "不允许踢出群主或管理员"
	CHAT_GROUP_MEMBER_NOT_EXISTS = "成员不在本群聊"
	CHAT_GROUP_CREATE_LIMIT      = "您的创建群聊数量超过限制"
	CHAT_GROUP_CREATE_OFF        = "创建群聊通道已关闭"

	MESSAGE_NOT_YOUR        = "只能撤回自己的消息"
	MESSAGE_CANT_REVOCATION = "超过2分钟消息无法撤回"

	SUCCESS          = "成功"
	REGISTER_SUCCESS = "注册成功"
	LOGIN_SUCCESS    = "登陆成功"
	SETTING_SUCCESS  = "设置成功"
	SEND_SUCCESS     = "发送成功"
	DELETE_SUCCESS   = "删除成功"
	REMOVE_SUCCESS   = "移除成功"
	ADD_SUCCESS      = "添加成功"
	APPLY_SUCCESS    = "申请成功"
	PROCESS_SUCCESS  = "处理成功"
	CREATE_SUCCESS   = "创建成功"
	EXIT_SUCCESS     = "退出成功"
	TRANSFER_SUCCESS = "转让成功"
	DISSOLVE_SUCCESS = "解散成功"
	EDIT_SUCCESS     = "修改成功"
	SAVE_SUCCESS     = "保存成功"
)

type Response struct {
	Code    int         `json:"code"`
	Result  interface{} `json:"result"`
	Message string      `json:"message"`
}

type PageResult struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
}
