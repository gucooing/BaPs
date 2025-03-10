package db

import (
	"database/sql"
	"errors"
)

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

func GetAllYostarMail() []*YostarMail {
	list := make([]*YostarMail, 0)
	SQL.Model(&YostarMail{}).Find(&list)
	return list
}

// AddYostarMailBySender 指定Sender创建数据
func AddYostarMailBySender(sender string) (*YostarMail, error) {
	data := new(YostarMail)
	data.Sender = sender
	err := SQL.Select("sender", sender).Create(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

// UpdateYostarMail 更新邮件数据
func UpdateYostarMail(data *YostarMail) error {
	if data == nil || data.MailIndex == 0 {
		return errors.New("YostarMail Nil")
	}
	return SQL.Model(&YostarMail{}).Where("mail_index = ?", data.MailIndex).Updates(data).Error
}
