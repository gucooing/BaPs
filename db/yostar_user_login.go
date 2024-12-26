package db

import (
	"errors"
)

type YostarUserLogin struct {
	AccountServerId  int64 `gorm:"primarykey;AUTO_INCREMENT"`
	YostarUid        int64 `gorm:"unique"`
	YostarLoginToken string
	Ban              bool
	BanMsg           string
}

// GetYoStarUserLoginByYostarUid 使用YostarUid拉取数据
func GetYoStarUserLoginByYostarUid(yostarUid int64) *YostarUserLogin {
	var data *YostarUserLogin
	SQL.Model(&YostarUserLogin{}).Where("yostar_uid = ?", yostarUid).First(&data)
	if data.AccountServerId == 0 {
		return nil
	}
	return data
}

// AddYoStarUserLoginByYostarUid 指定YostarUid创建数据
func AddYoStarUserLoginByYostarUid(yostarUid int64) (*YostarUserLogin, error) {
	data := new(YostarUserLogin)
	data.YostarUid = yostarUid
	err := SQL.Select("yostar_uid", yostarUid).Create(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

// UpdateYoStarUserLogin 更新账号数据
func UpdateYoStarUserLogin(data *YostarUserLogin) error {
	if data == nil || data.YostarUid == 0 {
		return errors.New("YoStarUserLoginData Nil")
	}
	return SQL.Model(&YostarUserLogin{}).Where("account_server_id = ?", data.AccountServerId).Updates(data).Error
}
