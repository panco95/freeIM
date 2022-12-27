package config

import (
	"context"
	"im/models"
	"im/pkg/database"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Config struct {
	mapping       map[string]interface{}
	mappingLocker sync.RWMutex
	mysqlClient   *database.Client
	log           *zap.SugaredLogger
}

func NewConfig(
	mysqlClient *database.Client,
) *Config {
	return &Config{
		mapping:       make(map[string]interface{}),
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

	m := make(map[string]interface{}, 0)
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

func (c *Config) GetAll() map[string]interface{} {
	c.mappingLocker.RLock()
	t := c.mapping
	defer c.mappingLocker.RUnlock()

	return t
}

func (c *Config) SetAll(mapping map[string]interface{}) {
	c.mappingLocker.Lock()
	defer c.mappingLocker.Unlock()

	c.mapping = mapping
}

func (c *Config) Get(key string) interface{} {
	c.mappingLocker.RLock()
	defer c.mappingLocker.RUnlock()

	if v, ok := c.mapping[key]; ok {
		return v
	}

	return nil
}

func (c *Config) GetString(key string) string {
	c.mappingLocker.RLock()
	defer c.mappingLocker.RUnlock()

	if v, ok := c.mapping[key]; ok {
		if t, ok := v.(string); ok {
			return t
		}
	}

	return ""
}

func (c *Config) GetInt(key string) int {
	c.mappingLocker.RLock()
	defer c.mappingLocker.RUnlock()

	if v, ok := c.mapping[key]; ok {
		if t, ok := v.(int); ok {
			return t
		}
	}

	return 0
}

func (c *Config) GetFloat64(key string) float64 {
	c.mappingLocker.RLock()
	defer c.mappingLocker.RUnlock()

	if v, ok := c.mapping[key]; ok {
		if t, ok := v.(float64); ok {
			return t
		}
	}

	return 0
}

func (c *Config) GetBool(key string) bool {
	c.mappingLocker.RLock()
	defer c.mappingLocker.RUnlock()

	if v, ok := c.mapping[key]; ok {
		if t, ok := v.(bool); ok {
			return t
		}
	}

	return false
}

func (c *Config) InitDBConfItem(name, key, val string) error {
	return c.mysqlClient.Db().
		Model(&models.Config{}).
		FirstOrCreate(&models.Config{
			Name: name,
			Key:  key,
			Val:  val,
		}, &models.Config{
			Key: key,
		}).Error
}

func (c *Config) InitDBConf() {
	g := func(name, key, val string) {
		err := c.InitDBConfItem(name, key, val)
		if err != nil {
			c.log.Errorf("InitDBConf %v", err)
		}
	}

	g("登陆验证码", "login_captcha", "true")
	g("附近的人最大距离", "near_friend_distance", "1000")
	g("附近的群最大距离", "near_chatgroup_distance", "1000")
	g("群默认最大成员数", "chatgroup_members_limit", "500")
	g("IP注册频率", "ip_register_rate", "0")
	g("IP限制时长", "ip_register_limit", "0")
	g("IP注册白名单", "ip_register_whitelist", "")
	g("IP注册黑名单", "ip_register_blacklist", "")
	g("手机端注册开关", "register_mobile", "true")
	g("PC注册开关", "register_pc", "true")
	g("注册邀请码开关", "register_invite", "false")
	g("提示新用户上传头像开关", "new_upload_avatar", "true")
	g("用户默认头像开关", "default_avatar", "true")
	g("启用短信注册开关", "register_sms", "true")
	g("添加好友模式", "add_friend_mode", "all")
	g("限制添加好友请求次数", "add_friend_apply_limit", "0")
	g("用户搜索模式", "search_account_mode", "all")
	g("用户添加管理号免认证开关", "add_manager_free", "true")
	g("管理号ip限制开关", "manager_ip_limit", "false")
	g("管理号ip登陆白名单", "manager_ip_whitelist", "")
	g("用户创建群开关", "create_chatgroup", "true")
	g("用户添加好友开关", "add_freind", "true")
	g("用户登录错误次数上限", "login_fail_limit", "0")
	g("显示在线状态开关", "display_online_status", "true")
	g("显示消息阅读状态开关", "chat_display_read", "true")
	g("显示是否输入状态开关", "chat_display_input", "true")
	g("是否显示聊天记录按钮开关", "chat_display_log_button", "true")
	g("重复消息最低间隔时间", "chat_repeat_rate", "0")
	g("普通消息最低间隔时间", "chat_rate", "0")
	g("文字消息字数限制", "chat_text_max_length", "0")
	g("消息列表时间标注间隔", "chat_list_time_limit", "300")
	g("敏感词列表", "sensitive_words", "")
	g("敏感词是否可以发送开关", "sensitive_send", "true")
	g("敏感词替代字符编辑", "sensitive_replace", "??")
	g("是否可以搜索群号加群开关", "chatgroup_search_id", "true")
	g("是否允许普通群成员退群开关", "chatgroup_exit", "true")
	g("群信息界面是否显示群成员开关", "chatgroup_display_members", "true")
	g("是否允许群主清屏开关", "chatgroup_owner_clean_message", "true")
	g("群标题是否显示群人数开关", "chatgroup_display_title_members", "true")
	g("显示邀请入群消息开关", "chatgroup_display_invite", "true")
	g("是否显示假的群人数开关", "chatgroup_display_false_members", "false")
	g("允许普通用户查看群信息开关", "chatgroup_display_info", "true")
	g("站外链接底部按钮名称", "discover_name", "发现")
}
