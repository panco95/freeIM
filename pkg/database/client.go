package database

import (
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Client struct {
	db *gorm.DB
}

func factory(
	db *gorm.DB,
	maxIdleConns int,
	maxOpenConns int,
	connMaxLifetime time.Duration,
) (*Client, error) {
	// gorm first方法忽略记录查不到err
	_ = db.Callback().Query().Before("gorm:query").Register("disable_raise_record_not_found", MaskNotDataError)

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	client := &Client{
		db: db,
	}
	return client, nil
}

func NewClickHouse(
	serverUrl string,
	maxIdleConns int,
	maxOpenConns int,
	connMaxLifetime time.Duration,
	logLevel logger.LogLevel,
) (*Client, error) {
	db, err := gorm.Open(clickhouse.New(clickhouse.Config{
		DSN: serverUrl,
	}), &gorm.Config{
		Logger: GetLogger(logLevel),
	})
	if err != nil {
		return nil, err
	}

	return factory(db, maxOpenConns, maxIdleConns, connMaxLifetime)
}

func NewMysql(
	serverUrl string,
	maxIdleConns int,
	maxOpenConns int,
	connMaxLifetime time.Duration,
	logLevel logger.LogLevel,
) (*Client, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: serverUrl,
	}), &gorm.Config{
		Logger: GetLogger(logLevel),
	})
	if err != nil {
		return nil, err
	}

	return factory(db, maxOpenConns, maxIdleConns, connMaxLifetime)
}

func (client *Client) Db() *gorm.DB {
	return client.db
}

func MaskNotDataError(gormDB *gorm.DB) {
	gormDB.Statement.RaiseErrorOnNotFound = false
}

func GetLogger(logLevel logger.LogLevel) logger.Interface {
	l := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logLevel,    // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,       // 禁用彩色打印
		},
	)
	return l
}

func (c *Client) AutoMigrate(models ...interface{}) error {
	return c.Db().
		AutoMigrate(models...)
}

func (c *Client) AutoMigrateWithSet(key, val string, models ...interface{}) error {
	return c.Db().
		Set(key, val).
		AutoMigrate(models...)
}

func (c *Client) SetAutoIncrementID(tableName string, autoID int) error {
	return c.Db().
		Exec("alter table " + tableName + " AUTO_INCREMENT " + strconv.Itoa(autoID)).
		Error
}
