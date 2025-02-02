package rank

import (
	"fmt"
	"sync"
	"time"

	"github.com/gucooing/BaPs/config"
	"github.com/gucooing/BaPs/db"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/zset"
	"gorm.io/gorm"
)

var RANKINFO *RankInfo

type RankInfo struct {
	raidRankZset map[int64]*zset.SortedSet[int64] // 赛季
	raidSync     sync.RWMutex
	SQL          *gorm.DB
}

func NewRank() *RankInfo {
	RANKINFO = &RankInfo{
		raidRankZset: make(map[int64]*zset.SortedSet[int64]),
		raidSync:     sync.RWMutex{},
	}
	// 初始化数据库
	RANKINFO.SQL = db.NewYostarRank(config.GetRaidRankDB())
	if RANKINFO.SQL == nil {
		logger.Error("YostarRank 数据库初始化失败")
		return nil
	}
	// 获取总力战信息
	for _, conf := range gdconf.GetRaidScheduleMap() {
		err := RANKINFO.SQL.Table(fmt.Sprintf("raid_rank_%v", conf.SeasonId)).AutoMigrate(&db.YostarRank{})
		if err != nil {
			logger.Error(err.Error())
		}
	}
	// 拉取当期总力战排名数据
	logger.Info("开始拉取当期总力战排名")
	if cur := gdconf.GetCurRaidSchedule(); cur != nil {
		RANKINFO.NewRaidRank(cur.SeasonId)
	} else {
		logger.Warn("缺少总力战当期排期")
	}
	logger.Info("拉取当期总力战排名完成")
	return RANKINFO
}

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
	for _, dbInfo := range db.GetAllYostarRank(x.SQL, seasonId) {
		s.Set(dbInfo.Score, dbInfo.Uid)
	}
	if conf.EndTime.After(time.Now()) {
		go func() {
			d := conf.EndTime.Add(1 * time.Hour).Sub(time.Now())
			ticker := time.NewTimer(d)
			logger.Debug("离总力战赛季结算还有:%s", d.String())
			select {
			case <-ticker.C:
				x.SettlementRaid(conf.SeasonId)
			}
		}()

	}
}

func (x *RankInfo) Close() {
	if x == nil {
		return
	}
	// 保存总力战数据
	x.raidSync.RLock()
	for seasonId, s := range x.raidRankZset {
		all := make([]*db.YostarRank, 0)
		s.RevRange(0, -1, func(score float64, uid int64) {
			all = append(all, &db.YostarRank{
				SeasonId: seasonId,
				Uid:      uid,
				Score:    score,
			})
		})
		err := db.UpAllYostarRank(x.SQL, all, seasonId)
		if err != nil {
			logger.Error("总力战排名保存失败SeasonId:%v,err:%s", seasonId, err.Error())
		}
	}
}
