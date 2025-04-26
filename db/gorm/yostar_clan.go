package db_gorm

import (
	"errors"
	dbstruct "github.com/gucooing/BaPs/db/struct"

	"gorm.io/gorm/clause"
)

func (x *DbGorm) GetYostarClanByServerId(serverId int64) *dbstruct.YostarClan {
	var data *dbstruct.YostarClan
	x.sql.Model(&dbstruct.YostarClan{}).Where("server_id = ?", serverId).First(&data)
	if data.ServerId == 0 {
		return nil
	}
	return data
}

func (x *DbGorm) GetYostarClanByClanName(clanName string) *dbstruct.YostarClan {
	var data *dbstruct.YostarClan
	x.sql.Model(&dbstruct.YostarClan{}).Where("clan_name = ?", clanName).First(&data)
	if data.ServerId == 0 {
		return nil
	}
	return data
}

func (x *DbGorm) AddYostarClanByClanName(clanName string) (*dbstruct.YostarClan, error) {
	data := new(dbstruct.YostarClan)
	data.ClanName = clanName
	err := x.sql.Select("clan_name", clanName).Create(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (x *DbGorm) UpdateYostarClan(data *dbstruct.YostarClan) error {
	if data == nil || data.ServerId == 0 {
		return errors.New("YostarClan Nil")
	}
	return x.sql.Model(&dbstruct.YostarClan{}).Where("server_id = ?", data.ServerId).Updates(data).Error
}

func (x *DbGorm) UpAllYostarClan(list []*dbstruct.YostarClan) error {
	err := x.sql.Clauses(clause.OnConflict{
		UpdateAll: true, // 如果主键冲突,则更新所有字段
	}).CreateInBatches(&list, 200).Error
	return err
}
