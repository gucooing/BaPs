package enter

import (
	"math/rand"
	"sync"
	"time"

	"github.com/gucooing/BaPs/common/rank"
	"github.com/gucooing/BaPs/gdconf"
)

type ArenaInfo struct {
	ArenaUserList   []*ArenaUser
	BattleArenaUser *ArenaUser // 被打的玩家
}

type ArenaUser struct {
	IsNpc       bool  // 是否是机器人
	Uid         int64 // 玩家情况下的uid
	Index       int64 // 机器人情况下的index/conf
	CharacterId int64 // 机器人情况下
	Rank        int64 // 排名
}

func (x *Session) GetArenaInfo() *ArenaInfo {
	if x.arenaInfo == nil {
		x.arenaInfo = &ArenaInfo{}
	}
	return x.arenaInfo
}

// GenArenaUserList 生成预战斗三人表
func (x *Session) GenArenaUserList(seasonId int64) {
	info := x.GetArenaInfo()
	info.ArenaUserList = make([]*ArenaUser, 0)

	r := rank.GetArenaRank(seasonId, x.AccountServerId)

	besaR := r

	ranks := make(map[int64]bool, 0)
	if r <= 3 {
		for i := int64(1); i < 5; i++ {
			if i == r {
				continue
			}
			ranks[i] = true
		}
	} else {
		for i := 0; i < 3; i++ {
			rr := rand.Int63n(besaR/5-1) + 1 + (besaR * 4 / 5)
			if ranks[rr] {
				i--
				continue
			}
			besaR = rr
			ranks[rr] = true
		}
	}

	for aernaRank := range ranks {
		uid := rank.GetArenaUidByRank(seasonId, aernaRank)
		au := &ArenaUser{
			IsNpc: false,
			Uid:   uid,
			Index: 0,
			Rank:  aernaRank,
		}

		if ps := GetSessionByUid(uid); ps != nil {

		} else {
			au.IsNpc = true
			au.Index = gdconf.RandGetArenaNPC().Index
			au.CharacterId = gdconf.RandCharacter()
		}

		info.ArenaUserList = append(info.ArenaUserList, au)
	}
}

func (x *Session) GetArenaUserList() []*ArenaUser {
	return x.GetArenaInfo().ArenaUserList
}

func (x *Session) GetArenaUserByIndex(index int32) *ArenaUser {
	list := x.GetArenaUserList()
	if int32(len(list)) < index+1 {
		return nil
	}
	return list[index]
}

func (x *Session) SetBattleArenaUser(au *ArenaUser) {
	x.arenaInfo.BattleArenaUser = au
}

func (x *Session) GetBattleArenaUser() *ArenaUser {
	return x.arenaInfo.BattleArenaUser
}

var arenaBattleList map[int64]bool
var arenaBattleSync sync.RWMutex

// AddArenaBattleRank 将排名添加到战斗进行中的列表
func AddArenaBattleRank(rank int64) bool {
	arenaBattleSync.Lock()
	defer arenaBattleSync.Unlock()
	if arenaBattleList == nil {
		arenaBattleList = make(map[int64]bool)
	}
	if arenaBattleList[rank] {
		return false // 重复添加
	}
	arenaBattleList[rank] = true

	return true
}

// DelArenaBattleRank 战斗结束删除队列
func DelArenaBattleRank(rank int64) bool {
	arenaBattleSync.Lock()
	defer arenaBattleSync.Unlock()
	if arenaBattleList[rank] {
		arenaBattleList[rank] = false
		return true
	}
	return false
}

// CheckArenaBattleRank 判断角色是否正在战斗中
func CheckArenaBattleRank(rank int64) bool {
	arenaBattleSync.RLock()
	defer arenaBattleSync.RUnlock()
	return arenaBattleList[rank]
}

// CheckArenaBattle 检查战斗是否结束/释放锁！！！！！！！！
func CheckArenaBattle(ticker *time.Ticker, attackRank, defenceRank int64) {
	select {
	case <-ticker.C:
		ticker.Stop()
		DelArenaBattleRank(attackRank)
		DelArenaBattleRank(defenceRank)
	}
}
