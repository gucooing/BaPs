package db_gorm

import (
	"errors"
	dbstruct "github.com/gucooing/BaPs/db/struct"
)

func (x *DbGorm) GetYostarUserByUid(uid int64) *dbstruct.YostarUser {
	var user *dbstruct.YostarUser
	err := x.sql.Model(&dbstruct.YostarUser{}).Where("uid = ?", uid).First(&user).Error
	if err != nil {
		return nil
	}
	return user
}

func (x *DbGorm) GetYostarUserByYostarUid(yostarUid int64) *dbstruct.YostarUser {
	var user *dbstruct.YostarUser
	err := x.sql.Model(&dbstruct.YostarUser{}).Where("yostar_uid = ?", yostarUid).First(&user).Error
	if err != nil {
		return nil
	}
	return user
}

func (x *DbGorm) AddYostarUserByYostarUid(yostarUid int64) (*dbstruct.YostarUser, error) {
	user := new(dbstruct.YostarUser)
	user.YostarUid = yostarUid
	err := x.sql.Select("yostar_uid", yostarUid).Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (x *DbGorm) UpdateYostarUser(data *dbstruct.YostarUser) error {
	if data == nil || data.YostarUid == 0 {
		return errors.New("YostarUserData Nil")
	}
	return x.sql.Model(&dbstruct.YostarUser{}).Where("uid = ?", data.Uid).Updates(data).Error
}
