package db_gorm

import (
	"errors"
	dbstruct "github.com/gucooing/BaPs/db/struct"
)

func (x *DbGorm) GetYostarAccountByYostarAccount(yostarAccount string) *dbstruct.YostarAccount {
	var account *dbstruct.YostarAccount
	x.sql.Model(&dbstruct.YostarAccount{}).Where("yostar_account = ?", yostarAccount).First(&account)
	if account.YostarUid == 0 {
		return nil
	}
	return account
}

func (x *DbGorm) GetYostarAccountByYostarUid(yostarUid int64) *dbstruct.YostarAccount {
	var account *dbstruct.YostarAccount
	err := x.sql.Model(&dbstruct.YostarAccount{}).Where("yostar_uid = ?", yostarUid).First(&account).Error
	if err != nil {
		return nil
	}
	return account
}

func (x *DbGorm) AddYostarAccountByYostarAccount(yostarAccount string) (*dbstruct.YostarAccount, error) {
	account := new(dbstruct.YostarAccount)
	account.YostarAccount = yostarAccount
	err := x.sql.Select("yostar_account", yostarAccount).Create(&account).Error
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (x *DbGorm) UpdateYostarAccount(data *dbstruct.YostarAccount) error {
	if data == nil || data.YostarUid == 0 {
		return errors.New("YostarAccountData Nil")
	}
	return x.sql.Model(&dbstruct.YostarAccount{}).Where("yostar_uid = ?", data.YostarUid).Updates(data).Error
}
