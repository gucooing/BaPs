package proto

import (
	"github.com/gucooing/BaPs/pkg/mx"
)

type RaidLobbyInfoDB struct {
	SeasonId                      int64
	Tier                          int32
	Ranking                       int64
	BestRankingPoint              int64
	TotalRankingPoint             int64
	ReceivedRankingRewardId       int64
	CanReceiveRankingReward       bool
	PlayingRaidDB                 *RaidDB
	ReceiveRewardIds              []int64
	ReceiveLimitedRewardIds       []int64
	ParticipateCharacterServerIds []int64
	PlayableHighestDifficulty     map[string]Difficulty
	SweepPointByRaidUniqueId      map[int64]int64
	SeasonStartDate               mx.MxTime
	SeasonEndDate                 mx.MxTime
	SettlementEndDate             mx.MxTime
	NextSeasonId                  int64
	NextSeasonStartDate           mx.MxTime
	NextSeasonEndDate             mx.MxTime
	NextSettlementEndDate         mx.MxTime
	ClanAssistUseInfo             *ClanAssistUseInfo
	RemainFailCompensation        map[int32]bool
}
