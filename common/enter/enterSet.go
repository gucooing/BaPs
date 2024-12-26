package enter

import (
	"sync"
	"time"
)

var es *EnterSet

type EnterSet struct {
	SessionMap     map[int64]*Session // 玩家信息-缓存
	sessionSync    sync.RWMutex
	EnterTicketMap map[string]*TicketInfo // 登录通行证字典-缓存
	ticketSync     sync.RWMutex
}

func getEnterSet() *EnterSet {
	if es == nil {
		es = &EnterSet{
			sessionSync: sync.RWMutex{},
			ticketSync:  sync.RWMutex{},
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
	}
}
