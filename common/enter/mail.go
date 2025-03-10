package enter

import (
	"encoding/json"

	"github.com/gucooing/BaPs/db"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/proto"
)

func (e *EnterSet) checkMail() {
	e.mailSync.Lock()
	defer e.mailSync.Unlock()
	e.MailMap = make(map[int64]*db.YostarMail)
	for _, bin := range db.GetAllYostarMail() {
		if bin.ParcelInfoListSql != "" {
			bin.ParcelInfoList = make([]*db.ParcelInfo, 0)
			if err := json.Unmarshal([]byte(bin.ParcelInfoListSql),
				&bin.ParcelInfoList); err != nil {
				logger.Warn("解析邮件附件失败,请检查邮件配置index:%v", bin.MailIndex)
			}
			for _, info := range bin.ParcelInfoList {
				if _, ok := proto.ParcelType_name[info.Type]; !ok {
					logger.Error("未知的邮箱附件类型:%v", info.Type)
					return
				}
			}
		}
		e.MailMap[bin.MailIndex] = bin
	}
}

func GetYostarMail() map[int64]*db.YostarMail {
	e := getEnterSet()
	e.mailSync.RLock()
	defer e.mailSync.RUnlock()
	list := make(map[int64]*db.YostarMail, 0)
	for k, v := range e.MailMap {
		list[k] = v
	}
	return list
}

func AddYostarMail(mail *db.YostarMail) bool {
	e := getEnterSet()
	bin, err := db.AddYostarMailBySender(mail.Sender)
	if err != nil {
		return false
	}
	mail.MailIndex = bin.MailIndex
	parcelInfoListSql, _ := json.Marshal(mail.ParcelInfoList)
	mail.ParcelInfoListSql = string(parcelInfoListSql)
	if err := db.UpdateYostarMail(mail); err != nil {
		return false
	}
	e.mailSync.RLock()
	defer e.mailSync.RUnlock()
	if e.MailMap == nil {
		e.MailMap = make(map[int64]*db.YostarMail)
	}
	e.MailMap[mail.MailIndex] = mail
	return true
}
