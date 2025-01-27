package db

import (
	"errors"

	"gorm.io/gorm"
)

type YostarFriend struct {
	AccountServerId int64 `gorm:"unique"`
	FriendInfo      string
}

// GetYostarFriendByAccountServerId 使用AccountServerId拉取数据
func GetYostarFriendByAccountServerId(accountServerId int64) *YostarFriend {
	var data *YostarFriend
	err := SQL.Model(&YostarFriend{}).Where("account_server_id = ?", accountServerId).First(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		data, _ = AddYostarFriendByYostarUid(accountServerId)
	}
	if data.AccountServerId == 0 {
		return nil
	}
	return data
}

// AddYostarFriendByYostarUid 指定YostarUid创建数据
func AddYostarFriendByYostarUid(accountServerId int64) (*YostarFriend, error) {
	data := new(YostarFriend)
	data.AccountServerId = accountServerId
	err := SQL.Select("account_server_id", accountServerId).Create(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

// UpdateYostarFriend 更新账号数据
func UpdateYostarFriend(data *YostarFriend) error {
	if data == nil || data.AccountServerId == 0 {
		return errors.New("YostarFriend Nil")
	}
	return SQL.Model(&YostarFriend{}).Where("account_server_id = ?", data.AccountServerId).Updates(data).Error
}
