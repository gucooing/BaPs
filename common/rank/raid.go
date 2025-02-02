package rank

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/db"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/zset"
)

// SettlementRaid 结算线程
func (x *RankInfo) SettlementRaid(seasonId int64) {
	logger.Info("总力战赛季结算开始,此时请勿关闭服务器,否则排行可能混乱")

	x.raidSync.Lock()
	s := x.raidRankZset[seasonId]
	delete(x.raidRankZset, seasonId)
	x.raidSync.Unlock()

	all := make([]*db.YostarRank, 0)
	ranking := int64(0)
	s.RevRange(0, -1, func(score float64, k int64) {
		ranking++
		all = append(all, &db.YostarRank{
			SeasonId: seasonId,
			Uid:      k,
			Score:    score,
		})
		ac := enter.GetSessionByAccountServerId(k)
		upAc := false
		if ac == nil {
			upAc = true
			ac = enter.GetSessionByUid(k)
		}
		if ac == nil {
			return
		}
		info := ac.PlayerBin.GetRankBin().GetCurRaidInfo()
		if info == nil || info.SeasonId != seasonId {
			return
		}
		info.Ranking = ranking
		if upAc {
			ac.UpDate()
			enter.DelSession(k)
		}
	})
	err := db.UpAllYostarRank(x.SQL, all, seasonId)
	if err != nil {
		logger.Error("总力战排名保存失败SeasonId:%v,err:%s", seasonId, err.Error())
	}
	logger.Info("总力战赛季结算结束")
}

func GetRaidRankZset(seasonId int64) *zset.SortedSet[int64] {
	if RANKINFO == nil {
		return nil
	}
	RANKINFO.raidSync.RLock()
	defer RANKINFO.raidSync.RUnlock()
	if _, ok := RANKINFO.raidRankZset[seasonId]; !ok {
		RANKINFO.NewRaidRank(seasonId)
	}
	return RANKINFO.raidRankZset[seasonId]
}

func GetRaidRank(seasonId, uid int64) int64 {
	s := GetRaidRankZset(seasonId)
	if s == nil {
		return 0
	}
	rank, _ := s.GetRank(uid, true)
	return rank + 1
}

func GetRaidScore(seasonId, uid int64) float64 {
	s := GetRaidRankZset(seasonId)
	if s == nil {
		return 0
	}
	_, score := s.GetRank(uid, true)
	return score
}

func SetRaidScore(seasonId, uid int64, score float64) {
	s := GetRaidRankZset(seasonId)
	if s == nil {
		return
	}
	s.Set(score, uid)
}

func GetRaidRankAndScore(seasonId, uid int64) (int64, float64) {
	s := GetRaidRankZset(seasonId)
	if s == nil {
		return 0, 0
	}
	rank, score := s.GetRank(uid, true)
	return rank + 1, score
}

// GetUidByRank 获取指定排名uid和分数
func GetUidByRank(seasonId, rank int64) (int64, float64) {
	s := GetRaidRankZset(seasonId)
	if s == nil {
		return 0, 0
	}
	return s.GetDataByRank(rank-1, true)
}
