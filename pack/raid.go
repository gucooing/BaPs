package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func RaidLogin(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.RaidLoginResponse)

	rsp.SeasonType = proto.RaidSeasonType_Open
}

func RaidLobby(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.RaidLobbyResponse)

	rsp.SeasonType = proto.RaidSeasonType_Open
	rsp.RaidLobbyInfoDB = &proto.SingleRaidLobbyInfoDB{
		ClearDifficulty: make([]proto.Difficulty, 0),
		RaidLobbyInfoDB: &proto.RaidLobbyInfoDB{
			NextSeasonId:          2,
			NextSeasonStartDate:   "2025-03-05T11:00:00",
			NextSeasonEndDate:     "2025-03-12T03:59:59",
			NextSettlementEndDate: "2025-03-12T23:59:59",
			PlayableHighestDifficulty: map[string]proto.Difficulty{
				"EN0010": proto.Difficulty_Hard,
			},
			ReceiveRewardIds: make([]int64, 0),
			RemainFailCompensation: map[int32]bool{
				0: true,
			},
			SeasonId:                 1,
			SeasonStartDate:          "2025-01-27T11:00:00",
			SeasonEndDate:            "2025-02-03T03:59:59",
			SettlementEndDate:        "2025-02-03T23:59:59",
			SweepPointByRaidUniqueId: make(map[int64]int64),
			Tier:                     1,

			Ranking:                       0,
			BestRankingPoint:              0,
			TotalRankingPoint:             0,
			ReceivedRankingRewardId:       0,
			CanReceiveRankingReward:       false,
			PlayingRaidDB:                 nil,
			ReceiveLimitedRewardIds:       nil,
			ParticipateCharacterServerIds: nil,
			ClanAssistUseInfo:             nil,
		},
	}
}
