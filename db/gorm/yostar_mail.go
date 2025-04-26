package db_gorm

import (
	"errors"
	dbstruct "github.com/gucooing/BaPs/db/struct"
)

func (x *DbGorm) GetAllYostarMail() []*dbstruct.YostarMail {
	list := make([]*dbstruct.YostarMail, 0)
	x.sql.Model(&dbstruct.YostarMail{}).Find(&list)
	return list
}

func (x *DbGorm) AddYostarMailBySender(sender string) (*dbstruct.YostarMail, error) {
	data := new(dbstruct.YostarMail)
	data.Sender = sender
	err := x.sql.Select("sender", sender).Create(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (x *DbGorm) UpdateYostarMail(data *dbstruct.YostarMail) error {
	if data == nil || data.MailIndex == 0 {
		return errors.New("YostarMail Nil")
	}
	return x.sql.Model(&dbstruct.YostarMail{}).Where("mail_index = ?", data.MailIndex).Updates(data).Error
}
