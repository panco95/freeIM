package dao

import (
	"context"
	"im/models"
	"im/pkg/database"

	"gorm.io/gorm"
)

type Account struct {
	mysqlClient *database.Client
}

func NewAccount(mysqlClient *database.Client) *Account {
	return &Account{
		mysqlClient: mysqlClient,
	}
}

// 查询账号
func (s *Account) Query(
	ctx context.Context,
	account *models.Account,
) (*models.Account, error) {
	db := s.mysqlClient.Db()
	queryAccount := &models.Account{}
	err := db.
		Where("(id <> 0 AND id = ?) OR (username <> '' AND username = ?) OR (email <> '' AND email = ?) OR (mobile <> '' AND mobile = ?)",
			account.ID,
			account.Username,
			account.Email,
			account.Mobile,
		).
		First(queryAccount).Error
	return queryAccount, err
}

// 更新账号
func (s *Account) Update(
	ctx context.Context,
	id uint,
	account *models.Account,
) error {
	db := s.mysqlClient.Db()
	err := db.Model(&models.Account{}).
		Where("id = ?", id).
		Updates(account).
		Error
	return err
}

// 创建账号
func (s *Account) Create(
	ctx context.Context,
	account *models.Account,
) error {
	db := s.mysqlClient.Db()
	err := db.Model(&models.Account{}).
		Create(account).Error
	return err
}

// 获取邀请码
func (s *Account) GetInviteCode(
	ctx context.Context,
	code string,
) (*models.InviteCode, error) {
	db := s.mysqlClient.Db()
	inviteCode := &models.InviteCode{}
	err := db.Model(&models.InviteCode{}).
		Where("code = ?", code).
		First(inviteCode).Error
	return inviteCode, err
}

// 邀请码使用次数+1
func (s *Account) IncrInviteCodeTimes(
	ctx context.Context,
	id uint,
) error {
	db := s.mysqlClient.Db()
	err := db.Model(&models.InviteCode{}).
		Where("id = ?", id).
		UpdateColumn("times", gorm.Expr("times + 1")).
		Error
	return err
}
