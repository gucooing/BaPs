package db_struct

import "database/sql"

type BlackDevice struct {
	DeviceId string `gorm:"unique"`
}

type YostarAccount struct {
	YostarUid     int64 `gorm:"primarykey;AUTO_INCREMENT"`
	YostarAccount string
	YostarToken   string
}

type YostarClan struct {
	ServerId int64  `gorm:"primarykey;AUTO_INCREMENT"`
	ClanName string `gorm:"unique"`
	ClanInfo string
}

type YostarFriend struct {
	AccountServerId int64 `gorm:"unique"`
	FriendInfo      string
}

type YostarGame struct {
	AccountServerId int64 `gorm:"unique"`
	BinData         []byte
}

type YostarMail struct {
	MailIndex         int64         `gorm:"primarykey;AUTO_INCREMENT"` // 唯一id
	Sender            string        // 发件人
	Comment           string        // 内容
	SendDate          sql.NullTime  // 开始时间
	ExpireDate        sql.NullTime  // 过期时间
	ParcelInfoListSql string        // 附件列表
	ParcelInfoList    []*ParcelInfo `gorm:"-"` // 附件
}

type ParcelInfo struct {
	Type int32 `json:"type"`
	Id   int64 `json:"id"`
	Num  int64 `json:"num"`
}

type YostarUser struct {
	Uid       int64 `gorm:"primarykey;AUTO_INCREMENT"`
	Token     string
	YostarUid int64
	DeviceId  string // 设备码
	ChannelId string
}

type YostarUserLogin struct {
	AccountServerId  int64 `gorm:"primarykey;AUTO_INCREMENT"`
	YostarUid        int64 `gorm:"unique"`
	YostarLoginToken string
	Ban              bool
	BanMsg           string
}
