package game

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/mx/proto"
	"github.com/gucooing/BaPs/pkg/logger"
)

var MinElapsedRealtime float64 = 10 // 最低战斗时间

func BattleCheck(s *enter.Session, info *proto.BattleSummary) {
	if s == nil || info == nil {
		return
	}
	level := 0
	// 战斗时间验证
	if info.ElapsedRealtime != 0 &&
		info.ElapsedRealtime < MinElapsedRealtime {
		level++
	}

	if level >= 1 {
		logger.Info("[UID:%v]玩家判定为战斗违规", s.AccountServerId)
	}
}
