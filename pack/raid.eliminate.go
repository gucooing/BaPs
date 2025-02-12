package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func EliminateRaidLogin(s *enter.Session, request, response proto.Message) {
	rsp := response.(*proto.EliminateRaidLoginResponse)

	rsp.SeasonType = game.GetEliminateRaidSeasonType()
	rsp.SweepPointByRaidUniqueId = make(map[int64]int64) // 扫荡信息
}

func EliminateRaidLobby(s *enter.Session, request, response proto.Message) {
	rsp := response.(*proto.EliminateRaidLobbyResponse)

	rsp.SeasonType = game.GetEliminateRaidSeasonType()
	// rsp.RaidGiveUpDB = &proto.RaidGiveUpDB{
	// 	Ranking:          0,
	// 	RankingPoint:     0,
	// 	BestRankingPoint: 0,
	// }
	rsp.RaidLobbyInfoDB = &proto.EliminateRaidLobbyInfoDB{
		RaidLobbyInfoDB: &proto.RaidLobbyInfoDB{
			SeasonId:                      1,
			Tier:                          0,
			Ranking:                       0,
			BestRankingPoint:              0,
			TotalRankingPoint:             0,
			ReceivedRankingRewardId:       0,
			CanReceiveRankingReward:       false,
			PlayingRaidDB:                 nil,
			ReceiveRewardIds:              nil,
			ReceiveLimitedRewardIds:       nil,
			ParticipateCharacterServerIds: nil,
			PlayableHighestDifficulty:     nil,
			SweepPointByRaidUniqueId:      nil,
			SeasonStartDate:               mx.MxTime{},
			SeasonEndDate:                 mx.MxTime{},
			SettlementEndDate:             mx.MxTime{},
			NextSeasonId:                  0,
			NextSeasonStartDate:           mx.MxTime{},
			NextSeasonEndDate:             mx.MxTime{},
			NextSettlementEndDate:         mx.MxTime{},
			ClanAssistUseInfo:             nil,
			RemainFailCompensation:        nil,
		},
		OpenedBossGroups:             make([]string, 0),
		BestRankingPointPerBossGroup: make(map[string]int64),
	}
}
