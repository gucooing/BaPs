package rank

import (
	"fmt"
	"sync"

	"github.com/gucooing/BaPs/config"
	"github.com/gucooing/BaPs/db"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/zset"
	"gorm.io/gorm"
)

var RANKINFO *RankInfo

type RankInfo struct {
	SQL *gorm.DB
	// 总力战
	raidRankZset map[int64]*zset.SortedSet[int64] // 赛季
	raidSync     sync.RWMutex
	// 大决战
	raidEliminateRankZset map[int64]*zset.SortedSet[int64] // 赛季
	raidEliminateSync     sync.RWMutex
}

func NewRank() *RankInfo {
	RANKINFO = &RankInfo{
		raidRankZset:          make(map[int64]*zset.SortedSet[int64]),
		raidSync:              sync.RWMutex{},
		raidEliminateRankZset: make(map[int64]*zset.SortedSet[int64]),
		raidEliminateSync:     sync.RWMutex{},
	}
	// 初始化数据库
	RANKINFO.SQL = db.NewYostarRank(config.GetRankDB())
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
	// 获取大决战信息
	for _, conf := range gdconf.GetRaidEliminateScheduleMap() {
		err := RANKINFO.SQL.Table(fmt.Sprintf("raid_eliminate_rank_%v", conf.SeasonId)).AutoMigrate(&db.YostarRank{})
		if err != nil {
			logger.Error(err.Error())
		}
	}
	// 拉取当期大决战排名数据
	logger.Info("开始拉取当期大决战排名")
	if cur := gdconf.GetCurRaidEliminateSchedule(); cur != nil {
		RANKINFO.NewRaidEliminateRank(cur.SeasonId)
	} else {
		logger.Warn("缺少大决战当期排期")
	}
	logger.Info("拉取当期大决战排名完成")
	return RANKINFO
}

func (x *RankInfo) Close() {
	if x == nil {
		return
	}
	// 保存总力战数据
	x.raidSync.Lock()
	for seasonId, s := range x.raidRankZset {
		all := make([]*db.YostarRank, 0)
		s.RevRange(0, -1, func(score float64, uid int64) {
			all = append(all, &db.YostarRank{
				SeasonId: seasonId,
				Uid:      uid,
				Score:    score,
			})
		})
		err := db.UpAllYostarRank(x.SQL, all, db.RaidUserTable(seasonId))
		if err != nil {
			logger.Error("总力战排名保存失败SeasonId:%v,err:%s", seasonId, err.Error())
		}
	}
	// 保存大决战数据
	x.raidEliminateSync.Lock()
	for seasonId, s := range x.raidEliminateRankZset {
		all := make([]*db.YostarRank, 0)
		s.RevRange(0, -1, func(score float64, uid int64) {
			all = append(all, &db.YostarRank{
				SeasonId: seasonId,
				Uid:      uid,
				Score:    score,
			})
		})
		err := db.UpAllYostarRank(x.SQL, all, db.RaidEliminateUserTable(seasonId))
		if err != nil {
			logger.Error("大决战排名保存失败SeasonId:%v,err:%s", seasonId, err.Error())
		}
	}
	logger.Info("排名数据保存完毕")
}
