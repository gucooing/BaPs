package enter

import (
	"time"
)

type Session struct {
	AccountServerId int64
	MxToken         string
	Time            time.Time
}

// 定时检查一次是否有用户长时间离线
func (e *EnterSet) checkSession() {
	for accountServerId, info := range GetAllSession() {
		if time.Now().After(info.Time.Add(30 * time.Minute)) {
			DelSession(accountServerId)
		}
	}
}

// 有锁 检查并处理重复登录
func (e *EnterSet) checkSessionRepeat(accountServerId int64) {
	e.sessionSync.Lock()
	defer e.sessionSync.Unlock()
	if e.SessionMap == nil {
		e.SessionMap = make(map[int64]*Session)
	}
	if _, ok := e.SessionMap[accountServerId]; ok {
		delete(e.SessionMap, accountServerId)
	}
}

// GetAllSession 获取全部在线玩家
func GetAllSession() map[int64]*Session {
	allSession := make(map[int64]*Session)
	e := getEnterSet()
	e.sessionSync.RLock()
	defer e.sessionSync.RUnlock()
	for k, v := range e.SessionMap {
		allSession[k] = v
	}
	return allSession
}

// DelSession 删除指定在线玩家
func DelSession(accountServerId int64) bool {
	e := getEnterSet()
	e.sessionSync.Lock()
	defer e.sessionSync.Unlock()
	if e.SessionMap == nil {
		e.SessionMap = make(map[int64]*Session)
	}
	if _, ok := e.SessionMap[accountServerId]; ok {
		delete(e.SessionMap, accountServerId)
		return true
	}
	return false
}
