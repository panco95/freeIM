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

func (c *Config) InitDBConf() {
	var err error

	err = c.mysqlClient.Db().
		Model(&models.Config{}).
		FirstOrCreate(&models.Config{
			Key: "login_captcha",
			Val: "true",
		}, &models.Config{
			Key: "login_captcha",
		}).Error
	if err != nil {
		c.log.Errorf("InitDBConf login_captcha %v", err)
	}

	err = c.mysqlClient.Db().
		Model(&models.Config{}).
		FirstOrCreate(&models.Config{
			Key: "near_friend_distance",
			Val: "1000000",
		}, &models.Config{
			Key: "near_friend_distance",
		}).Error
	if err != nil {
		c.log.Errorf("InitDBConf near_friend_distance %v", err)
	}

	err = c.mysqlClient.Db().
		Model(&models.Config{}).
		FirstOrCreate(&models.Config{
			Key: "near_chatgroup_distance",
			Val: "1000000",
		}, &models.Config{
			Key: "near_chatgroup_distance",
		}).Error
	if err != nil {
		c.log.Errorf("InitDBConf near_chatgroup_distance %v", err)
	}
}
