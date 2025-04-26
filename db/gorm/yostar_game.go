package db_gorm

import (
	"errors"
	dbstruct "github.com/gucooing/BaPs/db/struct"

	"gorm.io/gorm/clause"
)

func (x *DbGorm) GetYostarGameByAccountServerId(accountServerId int64) *dbstruct.YostarGame {
	var data *dbstruct.YostarGame
	x.sql.Model(&dbstruct.YostarGame{}).Where("account_server_id = ?", accountServerId).First(&data)
	if data.AccountServerId == 0 {
		return nil
	}
	return data
}

func (x *DbGorm) AddYostarGameByYostarUid(accountServerId int64) (*dbstruct.YostarGame, error) {
	data := new(dbstruct.YostarGame)
	data.AccountServerId = accountServerId
	err := x.sql.Select("account_server_id", accountServerId).Create(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (x *DbGorm) UpdateYostarGame(data *dbstruct.YostarGame) error {
	if data == nil || data.AccountServerId == 0 {
		return errors.New("YostarGame Nil")
	}
	return x.sql.Model(&dbstruct.YostarGame{}).Where("account_server_id = ?", data.AccountServerId).Updates(data).Error
}

func (x *DbGorm) UpAllYostarGame(list []*dbstruct.YostarGame) error {
	err := x.sql.Clauses(clause.OnConflict{
		UpdateAll: true, // 如果主键冲突,则更新所有字段
	}).CreateInBatches(&list, 200).Error
	return err
}
