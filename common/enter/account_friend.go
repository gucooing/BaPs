package enter

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/gucooing/BaPs/db"
	"github.com/gucooing/BaPs/pkg/logger"
)

var MaxCacheAccountFriendTime = 24 // 最大玩家缓存时间 单位: 小时

type AccountFriend struct {
	Uid                     int64          `json:"uid"`
	UpTime                  int64          `json:"-"`
	SyncAf                  sync.RWMutex   `json:"-"`
	AutoAcceptFriendRequest bool           `json:"auto_accept_friend_request"` // 是否自动同意好友申请
	FriendList              map[int64]bool `json:"friend_list"`                // 好友列表
	ReceivedList            map[int64]bool `json:"received_list"`              // 被申请好友列表
	SendReceivedList        map[int64]bool `json:"send_received_list"`         // 发送的好友申请列表
	BlockedList             map[int64]bool `json:"blocked_list"`               // 黑名单列表
}

// 每天4点检查一次是否有用户长时间离线然后离线掉好友数据
func (e *EnterSet) checkAccountFriend() {
	for accountServerId, info := range GetAllAccountFriend() {
		if time.Now().After(time.Unix(info.UpTime, 0).
			Add(time.Hour * time.Duration(MaxCacheAccountFriendTime))) {
			info.upDate()
			DelSession(accountServerId)
			logger.Debug("AccountId:%v,AccountFriend超时离线", accountServerId)
		}
	}
}

// GetAllAccountFriend 获取全部好友关系
func GetAllAccountFriend() map[int64]*AccountFriend {
	list := make(map[int64]*AccountFriend)
	e := getEnterSet()
	e.friendSync.RLock()
	defer e.friendSync.RUnlock()
	for v, k := range e.FriendMap {
		list[v] = k
	}
	return list
}

// GetAccountFriend 拉取好友关系
func GetAccountFriend(uid int64) *AccountFriend {
	s := getEnterSet()
	s.friendSync.RLock()
	if af, ok := s.FriendMap[uid]; ok {
		s.friendSync.RUnlock()
		return af
	}
	s.friendSync.RUnlock()
	af, err := DbGetAccountFriend(uid)
	if err != nil {
		return nil
	}
	s.friendSync.Lock()
	defer s.friendSync.Unlock()
	if s.FriendMap == nil {
		s.FriendMap = make(map[int64]*AccountFriend)
	}
	s.FriendMap[uid] = af
	return af
}

// DbGetAccountFriend 从db拉取数据
func DbGetAccountFriend(uid int64) (*AccountFriend, error) {
	af := new(AccountFriend)
	bin := db.GetYostarFriendByAccountServerId(uid)
	if bin == nil {
		return nil, errors.New("sql err")
	}
	json.Unmarshal([]byte(bin.FriendInfo), af)
	af.Uid = uid
	af.UpTime = time.Now().Unix()
	af.SyncAf = sync.RWMutex{}
	return af, nil
}

// UpDate 将玩家数据保存到数据库
func (x *AccountFriend) upDate() bool {
	if x == nil {
		return false
	}
	bin := &db.YostarFriend{
		AccountServerId: x.Uid,
	}
	friendInfo, err := json.Marshal(x)
	if err != nil {
		return false
	}
	bin.FriendInfo = string(friendInfo)
	err = db.UpdateYostarFriend(bin)
	if err != nil {
		return false
	}
	return true
}

func (x *AccountFriend) GetFriendList() map[int64]bool {
	if x == nil {
		return nil
	}
	x.SyncAf.RLock()
	defer x.SyncAf.RUnlock()
	friendList := make(map[int64]bool)
	for k, v := range x.FriendList {
		friendList[k] = v
	}
	return friendList
}
