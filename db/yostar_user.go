package db

import (
	"errors"
)

type YostarUser struct {
	Uid       int64 `gorm:"primarykey;AUTO_INCREMENT"`
	Token     string
	YostarUid int64
	DeviceId  string // 设备码
	ChannelId string
}

// GetYostarUserByUid 使用Uid拉取数据
func GetYostarUserByUid(uid int64) *YostarUser {
	var user *YostarUser
	err := SQL.Model(&YostarUser{}).Where("uid = ?", uid).First(&user).Error
	if err != nil {
		return nil
	}
	return user
}

// GetYostarUserByYostarUid 使用YostarUid拉取数据
func GetYostarUserByYostarUid(yostarUid int64) *YostarUser {
	var user *YostarUser
	err := SQL.Model(&YostarUser{}).Where("yostar_uid = ?", yostarUid).First(&user).Error
	if err != nil {
		return nil
	}
	return user
}

// AddYostarUserByYostarUid 指定YostarAccount创建数据
func AddYostarUserByYostarUid(yostarUid int64) (*YostarUser, error) {
	user := new(YostarUser)
	user.YostarUid = yostarUid
	err := SQL.Select("yostar_uid", yostarUid).Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateYostarUser 更新账号数据
func UpdateYostarUser(data *YostarUser) error {
	if data == nil || data.YostarUid == 0 {
		return errors.New("YostarUserData Nil")
	}
	return SQL.Model(&YostarUser{}).Where("uid = ?", data.Uid).Updates(data).Error
}
