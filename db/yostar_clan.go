package db

import (
	"errors"

	"gorm.io/gorm/clause"
)

type YostarClan struct {
	ServerId int64  `gorm:"primarykey;AUTO_INCREMENT"`
	ClanName string `gorm:"unique"`
	ClanInfo string
}

// GetYostarClanByServerId 使用ServerId拉取数据
func GetYostarClanByServerId(serverId int64) *YostarClan {
	var data *YostarClan
	SQL.Model(&YostarClan{}).Where("server_id = ?", serverId).First(&data)
	if data.ServerId == 0 {
		return nil
	}
	return data
}

// GetYostarClanByClanName 使用ClanName拉取数据
func GetYostarClanByClanName(clanName string) *YostarClan {
	var data *YostarClan
	SQL.Model(&YostarClan{}).Where("clan_name = ?", clanName).First(&data)
	if data.ServerId == 0 {
		return nil
	}
	return data
}

// AddYostarClanByClanName 指定ClanName创建数据
func AddYostarClanByClanName(clanName string) (*YostarClan, error) {
	data := new(YostarClan)
	data.ClanName = clanName
	err := SQL.Select("clan_name", clanName).Create(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

// UpdateYostarClan 更新社团数据
func UpdateYostarClan(data *YostarClan) error {
	if data == nil || data.ServerId == 0 {
		return errors.New("YostarClan Nil")
	}
	return SQL.Model(&YostarClan{}).Where("server_id = ?", data.ServerId).Updates(data).Error
}

// UpAllYostarClan 批量覆盖保存社团数据
func UpAllYostarClan(x []*YostarClan) error {
	err := SQL.Clauses(clause.OnConflict{
		UpdateAll: true, // 如果主键冲突,则更新所有字段
	}).CreateInBatches(&x, 200).Error
	return err
}
