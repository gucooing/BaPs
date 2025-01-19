package db

import (
	"database/sql"
)

type YostarMail struct {
	Index             int64         `gorm:"primarykey;AUTO_INCREMENT"` // 唯一id
	Sender            string        // 发件人
	Comment           string        // 内容
	SendDate          sql.NullTime  // 开始时间
	ExpireDate        sql.NullTime  // 过期时间
	ParcelInfoListSql string        // 附件列表
	ParcelInfoList    []*ParcelInfo `gorm:"-"` // 附件
}

type ParcelInfo struct {
	Type int32
	Id   int64
	Num  int64
}

func GetAllYostarMail() []*YostarMail {
	list := make([]*YostarMail, 0)
	SQL.Model(&YostarMail{}).Find(&list)
	return list
}
