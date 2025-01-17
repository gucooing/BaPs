package pack

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func ArenaLogin(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.ArenaLoginResponse)

	rsp.ArenaPlayerInfoDB = &proto.ArenaPlayerInfoDB{
		// CurrentSeasonId:          7,                              // 当前赛季编号
		// PlayerGroupId:            1,                              // 玩家组
		// CurrentRank:              1,                              // 当前排名
		// SeasonRecord:             1,                              // 本赛季最高记录
		// AllTimeRecord:            1,                              // 全部赛季最高记录
		// CumulativeTimeReward:     0,                              // 积累的时间奖励
		// TimeRewardLastUpdateTime: time.Now(),                     // 奖励最后更新时间
		// BattleEnterActiveTime:    time.Now().Add(24 * time.Hour), // 战斗冷却结束时间
		DailyRewardActiveTime: time.Now().Add(99999 * time.Hour), // 下一个每日排名奖励可领取时间
	}
}
