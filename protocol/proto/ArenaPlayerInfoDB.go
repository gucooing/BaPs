// Code auto-generated by gucooing & Hiro420 DO NOT EDIT.

package proto

import (
	"github.com/bytedance/sonic"
	"github.com/gucooing/BaPs/protocol/mx"
)

type ArenaPlayerInfoDB struct {
	CurrentSeasonId          int64
	PlayerGroupId            int64
	CurrentRank              int64
	SeasonRecord             int64
	AllTimeRecord            int64
	CumulativeTimeReward     int64
	TimeRewardLastUpdateTime mx.MxTime
	BattleEnterActiveTime    mx.MxTime
	DailyRewardActiveTime    mx.MxTime
}

func (x *ArenaPlayerInfoDB) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}
