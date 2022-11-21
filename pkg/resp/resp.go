package resp

const (
	ERROR = 1
)

const (
	SERVER_ERROR        = "服务器错误"
	PARAM_INVALID       = "参数不合法"
	ACCOUNT_NOT_FOUND   = "账号不存在"
	ACCOUNT_PWD_ERROR   = "账号或密码错误"
	ACCOUNT_LOCKED      = "账号被封禁"
	CAPTCHA_ERROR       = "验证码错误"
	CAPTCHA_EXPIRED     = "验证码过期，请刷新后重试"
	TIMEOUT             = "登录超时"
	ACCOUNT_EXISTS      = "账号已存在"
	ACCOUNT_NOT_EXISTS  = "账号不存在"
	ACCOUNT_HAS_CHINESE = "用户名不能包含中文"

	FRIEND_NOT_EXISTS       = "对方还不是您的好友"
	FRIEND_EXISTS           = "对方已经是您的好友"
	FRIEND_APPLY_EXISTS     = "对方还未验证您的好友请求"
	FRIEND_APPLY_NOT_EXISTS = "请求不存在"
	BLACKLIST_NOT_EXISTS    = "黑名单不存在"
	FRIEND_GROUP_EXISTS     = "分组已存在"
	FRIEND_GROUP_NOT_EXISTS = "分组不存在"

	SUCCESS          = "成功"
	REGISTER_SUCCESS = "注册成功"
	LOGIN_SUCCESS    = "登陆成功"
	SETTING_SUCCESS  = "设置成功"
	SEND_SUCCESS     = "发送成功"
	DELETE_SUCCESS   = "删除成功"
	REMOVE_SUCCESS   = "移除成功"
	ADD_SUCCESS      = "添加成功"
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
