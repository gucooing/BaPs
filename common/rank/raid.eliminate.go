package rank

import (
	"time"

	"github.com/gucooing/BaPs/db"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/zset"
)

// NewRaidEliminateRank 此操作会将排名强制覆盖成冷数据中的排名，建议仅用于初始化这个赛季时拉取冷数据使用
func (x *RankInfo) NewRaidEliminateRank(seasonId int64) {
	conf := gdconf.GetRaidEliminateScheduleInfo(seasonId)
	if x == nil || conf == nil {
		return
	}
	x.raidEliminateSync.Lock()
	if x.raidEliminateRankZset == nil {
		x.raidEliminateRankZset = make(map[int64]*zset.SortedSet[int64])
	}
	s := zset.New[int64]()
	x.raidEliminateRankZset[seasonId] = s
	x.raidEliminateSync.Unlock()
	for _, dbInfo := range db.GetAllYostarRank(x.SQL, db.RaidEliminateUserTable(seasonId)) {
		s.Set(dbInfo.Score, dbInfo.Uid)
	}
	// 赛季结束
	nextConf := gdconf.GetRaidEliminateScheduleInfo(conf.NextSeasonId)
	if nextConf == nil {
		logger.Warn("大决战缺少下一个赛季配置")
		return
	}
	if conf.StartTime.After(time.Now()) {
		go func() {
			d := nextConf.StartTime.Add(1 * time.Hour).Sub(time.Now())
			ticker := time.NewTimer(d)
			logger.Debug("离下一个大决战赛季开始还有:%s", d.String())
			select {
			case <-ticker.C:
				x.SettlementRaidEliminate(conf.SeasonId)
			}
		}()
	}
}

// SettlementRaidEliminate 结算线程
func (x *RankInfo) SettlementRaidEliminate(seasonId int64) {
	logger.Info("新的大决战赛季开始,正在回收旧赛季数据,此时请勿关闭服务器,否则排行可能混乱")

	x.raidEliminateSync.Lock()
	s := x.raidEliminateRankZset[seasonId]
	delete(x.raidEliminateRankZset, seasonId)
	x.raidEliminateSync.Unlock()

	all := make([]*db.YostarRank, 0)
	s.RevRange(0, -1, func(score float64, k int64) {
		all = append(all, &db.YostarRank{
			SeasonId: seasonId,
			Uid:      k,
			Score:    score,
		})
	})
	err := db.UpAllYostarRank(x.SQL, all, db.RaidEliminateUserTable(seasonId))
	if err != nil {
		logger.Error("旧赛季大决战排名保存失败SeasonId:%v,err:%s", seasonId, err.Error())
	}
	logger.Info("大决战旧赛季回收结束")
}

// GetRaidEliminateRankZset 注意,如果赛季不存在法会强制前往数据库拉取数据
func GetRaidEliminateRankZset(seasonId int64) *zset.SortedSet[int64] {
	if RANKINFO == nil {
		return nil
	}
	RANKINFO.raidEliminateSync.RLock()
	if _, ok := RANKINFO.raidEliminateRankZset[seasonId]; !ok {
		RANKINFO.raidEliminateSync.RUnlock()
		RANKINFO.NewRaidEliminateRank(seasonId)
		RANKINFO.raidEliminateSync.RLock()
	}
	defer RANKINFO.raidEliminateSync.RUnlock()
	return RANKINFO.raidEliminateRankZset[seasonId]
}

func GetRaidEliminateRank(seasonId, uid int64) int64 {
	s := GetRaidEliminateRankZset(seasonId)
	if s == nil {
		return 0
	}
	rank, _ := s.GetRank(uid, true)
	return rank + 1
}

func SetRaidEliminateScore(seasonId, uid int64, score float64) {
	s := GetRaidEliminateRankZset(seasonId)
	if s == nil {
		return
	}
	s.Set(score, uid)
}

func GetRaidEliminateRankAndScore(seasonId, uid int64) (int64, float64) {
	s := GetRaidEliminateRankZset(seasonId)
	if s == nil {
		return 0, 0
	}
	rank, score := s.GetRank(uid, true)
	return rank + 1, score
}

// GetUidByEliminateRank 获取指定排名uid和分数
func GetUidByEliminateRank(seasonId, rank int64) (int64, float64) {
	s := GetRaidEliminateRankZset(seasonId)
	if s == nil {
		return 0, 0
	}
	return s.GetDataByRank(rank-1, true)
}
