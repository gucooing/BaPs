package rank

import (
	"time"

	"github.com/gucooing/BaPs/db"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/zset"
)

/*
竞技场抽象
使用 账号id-排名 储存
排名作为分数,正序储存
使用len方法获取最低的新排名
*/

const DefaultArenaRank = 15000

// NewArenaRank 此操作会将排名强制覆盖成冷数据中的排名，建议仅用于初始化这个赛季时拉取冷数据使用
func (x *RankInfo) NewArenaRank(seasonId int64) {
	conf := gdconf.GetArenaSeasonExcelTable(seasonId)
	if x == nil || conf == nil {
		return
	}
	x.arenaSync.Lock()
	if x.arenaRankZset == nil {
		x.arenaRankZset = make(map[int64]*zset.SortedSet[int64])
	}
	s := zset.New[int64]()
	x.arenaRankZset[seasonId] = s
	x.arenaSync.Unlock()
	for _, dbInfo := range db.GetAllYostarRank(x.SQL, db.ArenaUserTable(seasonId)) {
		s.Set(dbInfo.Score, dbInfo.Uid)
	}
	// 赛季结束
	nextConf := gdconf.GetArenaSeasonExcelTable(conf.PrevSeasonId)
	if nextConf == nil {
		logger.Warn("竞技场缺少下一个赛季配置")
		return
	}
	startTime, err := time.Parse("2006-01-02 15:04:05", nextConf.SeasonStartDate)
	if err != nil {
		return
	}
	if startTime.After(time.Now()) {
		go func() {
			d := startTime.Add(1 * time.Hour).Sub(time.Now())
			ticker := time.NewTimer(d)
			logger.Debug("离下一个竞技场赛季开始还有:%s", d.String())
			select {
			case <-ticker.C:
				x.SettlementAren(conf.UniqueId)
			}
		}()
	}
}

// SettlementAren 结算线程
func (x *RankInfo) SettlementAren(seasonId int64) {
	logger.Info("新的竞技场赛季开始,正在回收旧赛季数据,此时请勿关闭服务器,否则排行可能混乱")

	x.arenaSync.Lock()
	s := x.arenaRankZset[seasonId]
	delete(x.arenaRankZset, seasonId)
	x.arenaSync.Unlock()

	all := make([]*db.YostarRank, 0)
	s.RevRange(0, -1, func(score float64, k int64) {
		all = append(all, &db.YostarRank{
			SeasonId: seasonId,
			Uid:      k,
			Score:    score,
		})
	})
	err := db.UpAllYostarRank(x.SQL, all, db.ArenaUserTable(seasonId))
	if err != nil {
		logger.Error("旧赛季竞技场排名保存失败SeasonId:%v,err:%s", seasonId, err.Error())
		return
	}
	logger.Info("竞技场旧赛季回收结束")
}

func GetArenaRankZset(seasonId int64) *zset.SortedSet[int64] {
	if RANKINFO == nil {
		return nil
	}
	RANKINFO.arenaSync.RLock()
	if _, ok := RANKINFO.arenaRankZset[seasonId]; !ok {
		RANKINFO.arenaSync.RUnlock()
		RANKINFO.NewArenaRank(seasonId)
		RANKINFO.arenaSync.RLock()
	}
	defer RANKINFO.arenaSync.RUnlock()
	return RANKINFO.arenaRankZset[seasonId]
}

func GetArenaRank(seasonId, uid int64) int64 {
	s := GetArenaRankZset(seasonId)
	if s == nil {
		return 0
	}
	rank, _ := s.GetRank(uid, false)
	if rank == zset.NoRank {
		return DefaultArenaRank
	}
	return rank + 1
}

// GetArenaUidByRank 获取指定排名uid和分数
func GetArenaUidByRank(seasonId, rank int64) (int64, float64) {
	s := GetArenaRankZset(seasonId)
	if s == nil {
		return 0, 0
	}
	return s.GetDataByRank(rank-1, false)
}
