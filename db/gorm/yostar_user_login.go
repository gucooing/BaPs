package db_gorm

import (
	"errors"
	dbstruct "github.com/gucooing/BaPs/db/struct"
)

func (x *DbGorm) GetYoStarUserLoginByYostarUid(yostarUid int64) *dbstruct.YostarUserLogin {
	var data *dbstruct.YostarUserLogin
	x.sql.Model(&dbstruct.YostarUserLogin{}).Where("yostar_uid = ?", yostarUid).First(&data)
	if data.AccountServerId == 0 {
		return nil
	}
	return data
}

func (x *DbGorm) AddYoStarUserLoginByYostarUid(yostarUid int64) (*dbstruct.YostarUserLogin, error) {
	data := new(dbstruct.YostarUserLogin)
	data.YostarUid = yostarUid
	err := x.sql.Select("yostar_uid", yostarUid).Create(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (x *DbGorm) UpdateYoStarUserLogin(data *dbstruct.YostarUserLogin) error {
	if data == nil || data.YostarUid == 0 {
		return errors.New("YoStarUserLoginData Nil")
	}
	return x.sql.Model(&dbstruct.YostarUserLogin{}).Where("account_server_id = ?", data.AccountServerId).Updates(data).Error
}
