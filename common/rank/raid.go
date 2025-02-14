package rank

import (
	"time"

	"github.com/gucooing/BaPs/db"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/zset"
)

// NewRaidRank 此操作会将排名强制覆盖成冷数据中的排名，建议仅用于初始化这个赛季时拉取冷数据使用
func (x *RankInfo) NewRaidRank(seasonId int64) {
	conf := gdconf.GetRaidScheduleInfo(seasonId)
	if x == nil || conf == nil {
		return
	}
	x.raidSync.Lock()
	if x.raidRankZset == nil {
		x.raidRankZset = make(map[int64]*zset.SortedSet[int64])
	}
	s := zset.New[int64]()
	x.raidRankZset[seasonId] = s
	x.raidSync.Unlock()
	for _, dbInfo := range db.GetAllYostarRank(x.SQL, db.RaidUserTable(seasonId)) {
		s.Set(dbInfo.Score, dbInfo.Uid)
	}
	// 赛季结束
	nextConf := gdconf.GetRaidScheduleInfo(conf.NextSeasonId)
	if nextConf == nil {
		logger.Warn("总力战缺少下一个赛季配置")
		return
	}
	if conf.StartTime.After(time.Now()) {
		go func() {
			d := nextConf.StartTime.Add(1 * time.Hour).Sub(time.Now())
			ticker := time.NewTimer(d)
			logger.Debug("离下一个总力战赛季开始还有:%s", d.String())
			select {
			case <-ticker.C:
				x.SettlementRaid(conf.SeasonId)
			}
		}()
	}
}

// SettlementRaid 结算线程
func (x *RankInfo) SettlementRaid(seasonId int64) {
	logger.Info("新的总力战赛季开始,正在回收旧赛季数据,此时请勿关闭服务器,否则排行可能混乱")

	x.raidSync.Lock()
	s := x.raidRankZset[seasonId]
	delete(x.raidRankZset, seasonId)
	x.raidSync.Unlock()

	all := make([]*db.YostarRank, 0)
	s.RevRange(0, -1, func(score float64, k int64) {
		all = append(all, &db.YostarRank{
			SeasonId: seasonId,
			Uid:      k,
			Score:    score,
		})
	})
	err := db.UpAllYostarRank(x.SQL, all, db.RaidUserTable(seasonId))
	if err != nil {
		logger.Error("旧赛季总力战排名保存失败SeasonId:%v,err:%s", seasonId, err.Error())
	}
	logger.Info("总力战旧赛季回收结束")
}

func GetRaidRankZset(seasonId int64) *zset.SortedSet[int64] {
	if RANKINFO == nil {
		return nil
	}
	RANKINFO.raidSync.RLock()
	if _, ok := RANKINFO.raidRankZset[seasonId]; !ok {
		RANKINFO.raidSync.RUnlock()
		RANKINFO.NewRaidRank(seasonId)
		RANKINFO.raidSync.RLock()
	}
	defer RANKINFO.raidSync.RUnlock()
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
