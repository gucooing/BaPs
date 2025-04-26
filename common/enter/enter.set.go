package enter

import (
	dbstruct "github.com/gucooing/BaPs/db/struct"
	"sync"
	"time"

	"github.com/gucooing/BaPs/pkg/alg"
)

var es *EnterSet

type EnterSet struct {
	SessionMap     map[int64]*Session // 玩家信息-缓存
	sessionSync    sync.RWMutex
	EnterTicketMap map[string]*TicketInfo // 登录通行证字典-缓存
	ticketSync     sync.RWMutex
	MailMap        map[int64]*dbstruct.YostarMail // 全服邮件
	mailSync       sync.RWMutex
	FriendMap      map[int64]*AccountFriend // 全部玩家的好友关系
	friendSync     sync.RWMutex
	YostarClan     map[int64]*YostarClan // 全部缓存社团
	YostarClanHash map[string]int64
	ycSync         sync.RWMutex
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
			friendSync:  sync.RWMutex{},
			ycSync:      sync.RWMutex{},
		}
		go es.Check()
	}
	return es
}

func (e *EnterSet) Check() {
	ticker := time.NewTicker(time.Second * 300)        // 五分钟验证一次
	friendTicker := time.NewTimer(alg.GetEveryDayH(4)) // 每天四点
	for {
		select {
		case <-ticker.C:
			e.checkTicket()
			e.checkSession()
			e.checkMail()
		case <-friendTicker.C:
			e.checkAccountFriend()
			e.checkYostarClan()
			friendTicker.Reset(alg.GetEveryDayH(4))
		}
	}
}
