package db

import (
	"errors"

	"gorm.io/gorm/clause"
)

type YostarGame struct {
	AccountServerId int64 `gorm:"unique"`
	BinData         []byte
}

// GetYostarGameByAccountServerId 使用AccountServerId拉取数据
func GetYostarGameByAccountServerId(accountServerId int64) *YostarGame {
	var data *YostarGame
	SQL.Model(&YostarGame{}).Where("account_server_id = ?", accountServerId).First(&data)
	if data.AccountServerId == 0 {
		return nil
	}
	return data
}

// AddYostarGameByYostarUid 指定YostarUid创建数据
func AddYostarGameByYostarUid(accountServerId int64) (*YostarGame, error) {
	data := new(YostarGame)
	data.AccountServerId = accountServerId
	err := SQL.Select("account_server_id", accountServerId).Create(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

// UpdateYostarGame 更新账号数据
func UpdateYostarGame(data *YostarGame) error {
	if data == nil || data.AccountServerId == 0 {
		return errors.New("YostarGame Nil")
	}
	return SQL.Model(&YostarGame{}).Where("account_server_id = ?", data.AccountServerId).Updates(data).Error
}

// UpAllYostarGame 批量覆盖保存账号数据
func UpAllYostarGame(x []*YostarGame) error {
	err := SQL.Clauses(clause.OnConflict{
		UpdateAll: true, // 如果主键冲突,则更新所有字段
	}).CreateInBatches(&x, 200).Error
	return err
}
