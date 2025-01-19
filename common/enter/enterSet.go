package enter

import (
	"sync"
	"time"

	"github.com/gucooing/BaPs/db"
)

var es *EnterSet

type EnterSet struct {
	SessionMap     map[int64]*Session // 玩家信息-缓存
	sessionSync    sync.RWMutex
	EnterTicketMap map[string]*TicketInfo // 登录通行证字典-缓存
	ticketSync     sync.RWMutex
	MailMap        map[int64]*db.YostarMail // 全服邮件
	mailSync       sync.RWMutex
}

func InitEnterSet() {
	e := getEnterSet()
	e.checkMail()
}

func getEnterSet() *EnterSet {
	if es == nil {
		es = &EnterSet{
			sessionSync: sync.RWMutex{},
			ticketSync:  sync.RWMutex{},
			mailSync:    sync.RWMutex{},
		}
		go es.Check()
	}
	return es
}

func (e *EnterSet) Check() {
	ticker := time.NewTicker(time.Second * 300) // 五分钟验证一次
	for {
		<-ticker.C
		e.checkTicket()
		e.checkSession()
		e.checkMail()
	}
}
