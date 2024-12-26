package db

import (
	"errors"
)

type YostarAccount struct {
	YostarUid     int64 `gorm:"primarykey;AUTO_INCREMENT"`
	YostarAccount string
	YostarToken   string
}

// GetYostarAccountByYostarAccount 使用YostarAccount拉取数据
func GetYostarAccountByYostarAccount(yostarAccount string) *YostarAccount {
	var account *YostarAccount
	SQL.Model(&YostarAccount{}).Where("yostar_account = ?", yostarAccount).First(&account)
	if account.YostarUid == 0 {
		return nil
	}
	return account
}

// GetYostarAccountByYostarUid 使用YostarUid拉取数据
func GetYostarAccountByYostarUid(yostarUid int64) *YostarAccount {
	var account *YostarAccount
	err := SQL.Model(&YostarAccount{}).Where("yostar_uid = ?", yostarUid).First(&account).Error
	if err != nil {
		return nil
	}
	return account
}

// AddYostarAccountByYostarAccount 指定YostarAccount创建数据
func AddYostarAccountByYostarAccount(yostarAccount string) (*YostarAccount, error) {
	account := new(YostarAccount)
	account.YostarAccount = yostarAccount
	err := SQL.Select("yostar_account", yostarAccount).Create(&account).Error
	if err != nil {
		return nil, err
	}
	return account, nil
}

// UpdateYostarAccount 更新账号数据
func UpdateYostarAccount(data *YostarAccount) error {
	if data == nil || data.YostarUid == 0 {
		return errors.New("YostarAccountData Nil")
	}
	return SQL.Model(&YostarAccount{}).Where("yostar_uid = ?", data.YostarUid).Updates(data).Error
}
