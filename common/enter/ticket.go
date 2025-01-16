package enter

import (
	"time"
)

type TicketInfo struct {
	AccountServerId int64
	YostarUID       int64
	Ticket          string
	EndTime         time.Time
}

// 定时检查一次是否有通行证过期
func (e *EnterSet) checkTicket() {
	for ticket, info := range GetAllEnterTicketInfo() {
		if time.Now().After(info.EndTime) {
			DelEnterTicket(ticket)
		}
	}
}

// AddEnterTicket 有锁 添加EnterTicket
func AddEnterTicket(accountServerId, yostarUID int64, ticker string) bool {
	e := getEnterSet()

	info := &TicketInfo{
		AccountServerId: accountServerId,
		YostarUID:       yostarUID,
		Ticket:          ticker,
		EndTime:         time.Now().Add(30 * time.Minute), // 30分钟有效期
	}
	e.ticketSync.Lock()
	defer e.ticketSync.Unlock()
	if e.EnterTicketMap == nil {
		e.EnterTicketMap = make(map[string]*TicketInfo)
	}
	e.EnterTicketMap[ticker] = info

	return true
}

// DelEnterTicket 有锁 删除Ticket,如果此Ticket存在返回true,不存在返回false
func DelEnterTicket(ticker string) bool {
	e := getEnterSet()
	e.ticketSync.Lock()
	defer e.ticketSync.Unlock()
	if e.EnterTicketMap == nil {
		e.EnterTicketMap = make(map[string]*TicketInfo)
	}
	if _, ok := e.EnterTicketMap[ticker]; ok {
		delete(e.EnterTicketMap, ticker)
		return true
	}
	return false
}

// GetEnterTicketInfo 有锁 通过ticker获取登录信息
func GetEnterTicketInfo(ticker string) *TicketInfo {
	e := getEnterSet()
	e.ticketSync.RLock()
	defer e.ticketSync.RUnlock()
	if info, ok := e.EnterTicketMap[ticker]; ok {
		if time.Now().After(info.EndTime) {
			return nil
		}
		return info
	}
	return nil
}

// GetAllEnterTicketInfo 有锁 获取全部登录信息
func GetAllEnterTicketInfo() map[string]*TicketInfo {
	allTicketInfo := make(map[string]*TicketInfo)
	e := getEnterSet()
	e.ticketSync.RLock()
	defer e.ticketSync.RUnlock()
	for k, v := range e.EnterTicketMap {
		allTicketInfo[k] = v
	}
	return allTicketInfo
}
