package config

import (
	"context"
	"im/models"
	"im/pkg/database"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Config struct {
	mapping       map[string]string
	mappingLocker sync.RWMutex
	mysqlClient   *database.Client
	log           *zap.SugaredLogger
}

func NewConfig(
	mysqlClient *database.Client,
) *Config {
	return &Config{
		mapping:       make(map[string]string),
		mappingLocker: sync.RWMutex{},
		mysqlClient:   mysqlClient,
		log:           zap.S().With("module", "services.system.config"),
	}
}

func (c *Config) AutoRefresh(interval time.Duration) {
	for {
		err := c.RefreshConfigs(context.Background())
		if err != nil {
			time.Sleep(interval)
			continue
		}

		time.Sleep(interval)
	}
}

func (c *Config) RefreshConfigs(ctx context.Context) error {
	cfgs, err := c.SelectAll(ctx)
	if err != nil {
		c.log.Errorf("RefreshConfigs %v", err)
		return err
	}

	m := make(map[string]string, 0)
	for _, v := range cfgs {
		m[v.Key] = v.Val
	}

	c.SetAll(m)

	return nil
}

func (c *Config) SelectAll(ctx context.Context) ([]*models.Config, error) {
	db := c.mysqlClient.Db()
	cfgs := make([]*models.Config, 0)
	err := db.Model(&models.Config{}).
		Order("id asc").
		Find(&cfgs).
		Error
	return cfgs, err
}

func (c *Config) GetAll() map[string]string {
	c.mappingLocker.RLock()
	t := c.mapping
	defer c.mappingLocker.RUnlock()

	return t
}

func (c *Config) SetAll(mapping map[string]string) {
	c.mappingLocker.Lock()
	defer c.mappingLocker.Unlock()

	c.mapping = mapping
}

func (c *Config) Get(key string) string {
	c.mappingLocker.RLock()
	defer c.mappingLocker.RUnlock()

	if v, ok := c.mapping[key]; ok {
		return v
	}

	return ""
}

func (c *Config) GetString(key string) string {
	c.mappingLocker.RLock()
	defer c.mappingLocker.RUnlock()

	if v, ok := c.mapping[key]; ok {
		return v
	}

	return ""
}

func (c *Config) GetInt(key string) int {
	c.mappingLocker.RLock()
	defer c.mappingLocker.RUnlock()

	if v, ok := c.mapping[key]; ok {
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0
		}
		return i
	}

	return 0
}

func (c *Config) GetFloat64(key string) float64 {
	c.mappingLocker.RLock()
	defer c.mappingLocker.RUnlock()

	if v, ok := c.mapping[key]; ok {
		f, err := strconv.ParseFloat(v, 2)
		if err != nil {
			return 0
		}
		return f
	}

	return 0
}

func (c *Config) InitDBConfItem(name, key, val, intro string) error {
	return c.mysqlClient.Db().
		Model(&models.Config{}).
		FirstOrCreate(&models.Config{
			Name:  name,
			Key:   key,
			Val:   val,
			Intro: intro,
		}, &models.Config{
			Key: key,
		}).Error
}

func (c *Config) InitDBConf() {
	g := func(name, key, val, intro string) {
		err := c.InitDBConfItem(name, key, val, intro)
		if err != nil {
			c.log.Errorf("InitDBConf %v", err)
		}
	}

	g("短信宝API用户名", "smsbao_username", "a97392828", "短信宝用户名")
	g("短信宝API密码", "smsbao_password", "0c3b7189bb23477db90f6fa5352e7221", "短信宝密码")
	g("短信宝API通道", "smsbao_send_range", "local", "短信宝有两种接口通道，local：本地，world：全球")
	g("短信签名", "sms_sign", "IM", "会在短信开头加上：【签名】")
	g("短信模板", "sms_template", "您的验证码为{code}", "{code}为必填字符串，真实情况会替换验证码")
	g("登陆验证码", "login_captcha", "false", "账号注册是否需要图形验证码，需要：true，不需要： false")
	g("附近的人最大距离", "near_friend_distance", "1000", "附近的人最大距离，单位km")
	g("附近的群最大距离", "near_chatgroup_distance", "1000", "附近的群最大距离，单位km")
	g("创建群聊限制数量", "chatgroup_create_limit", "20", "单个用户可以创建的群聊数量")
	g("群聊最大成员数", "chatgroup_members_limit", "500", "群聊可以加入的人数限制")
	g("IP注册频率", "ip_register_rate", "0", "单个IP注册了多少个账号无法注册，0表示无限制")
	g("IP限制时长", "ip_register_limit", "0", "单个ip注册上限等待多少天后解除限制，单位：天")
	g("IP注册白名单", "ip_register_whitelist", "", "不限制注册限制的IP列表，用;符号分割多个")
	g("IP注册黑名单", "ip_register_blacklist", "", "不允许注册的IP列表，用;符号分割多个")
	g("注册邀请码开关", "register_invite", "false", "注册是否必填邀请码，开：true，关：false")
	g("提示新用户上传头像开关", "need_upload_avatar", "true", "开：true，关：false")
	g("用户默认头像", "default_account_avatar", "", "用户默认头像地址(七牛云)")
	g("群聊默认头像", "default_chatgroup_avatar", "", "群聊默认头像地址(七牛云)")
	g("手机端注册开关", "register_mobile", "true", "是否允许手机端注册，允许：true，不允许：false")
	g("PC注册开关", "register_pc", "true", "是否允许PC端注册，允许：true，不允许：false")
	g("短信注册开关", "register_sms", "true", "开：true，关：false")
	g("邮箱注册开关", "register_email", "true", "开：true，关：false")
	g("账号注册开关", "register_account", "true", "开：true，关：false")
	g("添加好友模式", "add_friend_mode", "all", "")
	g("限制添加好友请求次数", "add_friend_apply_limit", "0", "0：不限制")
	g("用户搜索模式", "search_account_mode", "all", "")
	g("用户添加管理号免认证开关", "add_manager_free", "true", "开：true，关：false")
	g("管理号ip限制开关", "manager_ip_limit", "false", "开：true，关：false")
	g("管理号ip登陆白名单", "manager_ip_whitelist", "", "IP列表，用;符号分割多个")
	g("用户创建群开关", "create_chatgroup", "true", "开：true，关：false")
	g("用户添加好友开关", "add_freind", "true", "开：true，关：false")
	g("用户登录错误次数上限", "login_fail_limit", "0", "0：不限制")
	g("显示在线状态开关", "display_online_status", "true", "开：true，关：false")
	g("显示消息阅读状态开关", "chat_display_read", "true", "开：true，关：false")
	g("显示是否输入状态开关", "chat_display_input", "true", "开：true，关：false")
	g("是否显示聊天记录按钮开关", "chat_display_log_button", "true", "开：true，关：false")
	g("重复消息最低间隔时间", "chat_repeat_rate", "0", "0：不限制，单位：秒")
	g("普通消息最低间隔时间", "chat_rate", "0", "0：不限制，单位：秒")
	g("文字消息字数限制", "chat_text_max_length", "0", "0：不限制")
	g("消息列表时间标注间隔", "chat_list_time_limit", "300", "单位：秒")
	g("敏感词列表", "sensitive_words", "", "多个用;分隔")
	g("敏感词是否可以发送开关", "sensitive_send", "true", "开：true，关：false")
	g("敏感词替代字符编辑", "sensitive_replace", "??", "")
	g("是否可以搜索群号加群开关", "chatgroup_search_id", "true", "开：true，关：false")
	g("是否允许普通群成员退群开关", "chatgroup_exit", "true", "开：true，关：false")
	g("群信息界面是否显示群成员开关", "chatgroup_display_members", "true", "开：true，关：false")
	g("是否允许群主清屏开关", "chatgroup_owner_clean_message", "true", "开：true，关：false")
	g("群标题是否显示群人数开关", "chatgroup_display_title_members", "true", "开：true，关：false")
	g("显示邀请入群消息开关", "chatgroup_display_invite", "true", "开：true，关：false")
	g("是否显示假的群人数开关", "chatgroup_display_false_members", "false", "开：true，关：false")
	g("允许普通用户查看群信息开关", "chatgroup_display_info", "true", "开：true，关：false")
	g("站外链接底部按钮名称", "discover_name", "发现", "")
}
