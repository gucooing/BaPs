package enter

import (
	"encoding/json"
	dbstruct "github.com/gucooing/BaPs/db/struct"

	"github.com/gucooing/BaPs/db"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/proto"
)

func (e *EnterSet) checkMail() {
	e.mailSync.Lock()
	defer e.mailSync.Unlock()
	e.MailMap = make(map[int64]*dbstruct.YostarMail)
	for _, bin := range db.GetDBGame().GetAllYostarMail() {
		if bin.ParcelInfoListSql != "" {
			bin.ParcelInfoList = make([]*dbstruct.ParcelInfo, 0)
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

func GetYostarMail() map[int64]*dbstruct.YostarMail {
	e := getEnterSet()
	e.mailSync.RLock()
	defer e.mailSync.RUnlock()
	list := make(map[int64]*dbstruct.YostarMail, 0)
	for k, v := range e.MailMap {
		list[k] = v
	}
	return list
}

func AddYostarMail(mail *dbstruct.YostarMail) bool {
	e := getEnterSet()
	bin, err := db.GetDBGame().AddYostarMailBySender(mail.Sender)
	if err != nil {
		return false
	}
	mail.MailIndex = bin.MailIndex
	parcelInfoListSql, _ := json.Marshal(mail.ParcelInfoList)
	mail.ParcelInfoListSql = string(parcelInfoListSql)
	if err := db.GetDBGame().UpdateYostarMail(mail); err != nil {
		return false
	}
	e.mailSync.RLock()
	defer e.mailSync.RUnlock()
	if e.MailMap == nil {
		e.MailMap = make(map[int64]*dbstruct.YostarMail)
	}
	e.MailMap[mail.MailIndex] = mail
	return true
}
