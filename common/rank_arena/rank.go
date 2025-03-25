package ranar

import (
	"sync"
)

const (
	NoRank           = -1
	DefaultArenaRank = 15000
)

type RankArena struct {
	dict    map[int64]int64 // 排名 uid
	dictUid map[int64]int64 // uid 排名
	lock    sync.RWMutex
	length  int64
}

func New() *RankArena {
	r := &RankArena{
		dict:    make(map[int64]int64),
		dictUid: make(map[int64]int64),
		lock:    sync.RWMutex{},
	}

	return r
}

func (r *RankArena) Length() int64 {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.length
}

func (r *RankArena) Set(rank, uid int64) {
	r.lock.Lock()
	defer r.lock.Unlock()
	v, ok := r.dict[rank]      // 查询这个位置上是否有其他用户
	old, ok2 := r.dictUid[uid] // 查询我是否有其他位置的排名

	r.dict[rank] = uid
	r.dictUid[uid] = rank
	if ok {
		// 如果该排名上已有玩家，并且不是自己，则交换uid
		if v != uid {
			if ok2 {
				r.dict[old] = v
				r.dictUid[v] = old
			} else {
				r.length--
				delete(r.dictUid, v)
			}
		}

	} else if ok2 {
		delete(r.dict, old)
	} else {
		r.length++
	}
}

// GetUid 通过排名拉取uid
func (r *RankArena) GetUid(rank int64) int64 {
	r.lock.RLock()
	defer r.lock.RUnlock()
	v, ok := r.dict[rank]
	if !ok {
		return NoRank
	}
	return v
}

// GetRank 通过uid拉取排名
func (r *RankArena) GetRank(uid int64) int64 {
	r.lock.RLock()
	defer r.lock.RUnlock()
	v, ok := r.dictUid[uid]
	if !ok {
		return DefaultArenaRank
	}
	return v
}

// GetAll 获取全部数据
func (r *RankArena) GetAll(f func(uid, rank int64)) {
	list := make(map[int64]int64)

	r.lock.RLock()
	for k, v := range r.dictUid {
		list[k] = v
	}
	r.lock.RUnlock()
	for k, v := range list {
		f(k, v)
	}
}
