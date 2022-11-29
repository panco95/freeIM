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
	err := db.Model(&Config{}).
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
	t := ""
	c.mappingLocker.RLock()
	defer c.mappingLocker.RUnlock()

	if v, ok := c.mapping[key]; ok {
		t = v
	}

	return t
}

func (c *Config) GetInt(key string) int {
	t := 0
	c.mappingLocker.RLock()
	defer c.mappingLocker.RUnlock()

	if v, ok := c.mapping[key]; ok {
		t1 := v
		t2, err := strconv.Atoi(t1)
		if err == nil {
			t = t2
		}
	}

	return t
}
