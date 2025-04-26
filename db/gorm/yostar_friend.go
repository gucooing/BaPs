package db_gorm

import (
	"errors"
	dbstruct "github.com/gucooing/BaPs/db/struct"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (x *DbGorm) GetYostarFriendByAccountServerId(accountServerId int64) *dbstruct.YostarFriend {
	var data *dbstruct.YostarFriend
	err := x.sql.Model(&dbstruct.YostarFriend{}).Where("account_server_id = ?", accountServerId).First(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		data, _ = x.AddYostarFriendByYostarUid(accountServerId)
	}
	if data.AccountServerId == 0 {
		return nil
	}
	return data
}

func (x *DbGorm) AddYostarFriendByYostarUid(accountServerId int64) (*dbstruct.YostarFriend, error) {
	data := new(dbstruct.YostarFriend)
	data.AccountServerId = accountServerId
	err := x.sql.Select("account_server_id", accountServerId).Create(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (x *DbGorm) UpdateYostarFriend(data *dbstruct.YostarFriend) error {
	if data == nil || data.AccountServerId == 0 {
		return errors.New("YostarFriend Nil")
	}
	return x.sql.Model(&dbstruct.YostarFriend{}).Where("account_server_id = ?", data.AccountServerId).Updates(data).Error
}

func (x *DbGorm) UpAllYostarFriend(list []*dbstruct.YostarFriend) error {
	err := x.sql.Clauses(clause.OnConflict{
		UpdateAll: true, // 如果主键冲突,则更新所有字段
	}).CreateInBatches(&list, 200).Error
	return err
}
