package pack

import (
	"github.com/gucooing/BaPs/protocol/mx"
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/config"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/protocol/proto"
)

func MissionList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.MissionListResponse)

	rsp.ClearedOrignalMissionIds = make([]int64, 0)                  // 已清除任务
	rsp.MissionHistoryUniqueIds = make([]int64, 0)                   // 完成且已领取奖励的任务
	rsp.ProgressDBs = make([]*proto.MissionProgressDB, 0)            // 进度数据
	rsp.DailySuddenMissionInfo = game.GetDailySuddenMissionInfoDb(s) // 每日任务

	addProgressDB := func(bin *sro.MissionInfo) {
		if bin == nil {
			return
		}
		if bin.Finish {
			rsp.MissionHistoryUniqueIds = append(rsp.MissionHistoryUniqueIds, bin.MissionId)
		} else {
			info := game.GetMissionProgressDB(s, bin.MissionId)
			if info != nil {
				rsp.ProgressDBs = append(rsp.ProgressDBs, info)
			}
		}
	}

	// rsp.ProgressDBs = append(rsp.ProgressDBs, &proto.MissionProgressDB{
	// 	MissionUniqueId:    2300,
	// 	Complete:           true,
	// 	StartTime:          time.Now(),
	// 	ProgressParameters: map[int64]int64{2: 5},
	// })
	// 添加每日任务
	if bin := game.GetDayMissionInfo(s); bin != nil {
		for _, info := range bin.MissionList {
			addProgressDB(info)
		}
	}
	// 添加每周任务

	// 添加成就任务
}

func MissionReward(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.MissionRewardRequest)
	rsp := response.(*proto.MissionRewardResponse)

	bin := game.GetMissionInfo(s, req.MissionUniqueId)
	conf := gdconf.GetMissionExcelTable(req.MissionUniqueId)
	if bin == nil || conf == nil || bin.Finish {
		return
	}
	s.MissionReward(bin)
	rsp.AddedHistoryDB = game.GetMissionHistoryDB(s, req.MissionUniqueId)
	// 获取奖励
	parcelResultList := game.GetParcelResultList(conf.MissionRewardParcelType,
		conf.MissionRewardParcelId, conf.MissionRewardAmount, false)
	rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResultList)
}

func GuideMissionSeasonList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.GuideMissionSeasonListResponse)

	rsp.MissionProgressDBs = make([]*proto.MissionProgressDB, 0)
	rsp.GuideMissionSeasonDBs = []*proto.GuideMissionSeasonDB{
		{
			SeasonId:                  3,
			LoginCount:                1,
			StartDate:                 mx.Now().Add(-time.Hour * 24),
			LoginDate:                 mx.Now(),
			IsComplete:                false,
			IsFinalMissionComplete:    true,
			CollectionItemReceiveDate: mx.MxTime{},
		},
		{
			SeasonId:                  1000,
			LoginCount:                1,
			StartDate:                 mx.Now().Add(-time.Hour * 24),
			LoginDate:                 mx.Now(),
			IsComplete:                false,
			IsFinalMissionComplete:    true,
			CollectionItemReceiveDate: mx.MxTime{},
		},
		{
			SeasonId:               1001,
			IsFinalMissionComplete: true,
		},
	}
}

func AccountGetTutorial(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.AccountGetTutorialResponse)

	if config.GetTutorial() {
		rsp.TutorialIds = game.GetTutorialList(s)
	} else {
		rsp.TutorialIds = make([]int64, 0)
		for i := 1; i < 28; i++ {
			rsp.TutorialIds = append(rsp.TutorialIds, int64(i))
		}
	}
}

func ScenarioSkip(s *enter.Session, request, response mx.Message) {

}

func AccountSetTutorial(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.AccountSetTutorialRequest)

	game.FinishTutorial(s, req.TutorialIds)
}

func MissionSync(s *enter.Session, request, response mx.Message) {
	if time.Now().After(alg.GetDayH(18)) {
		s.FinishMission(proto.MissionCompleteConditionType_Reset_LoginAtSpecificTime, 1, nil)
	}
}

func ScenarioList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.ScenarioListResponse)

	rsp.ScenarioGroupHistoryDBs = game.GetScenarioGroupHistoryDBs(s)
	rsp.ScenarioHistoryDBs = make([]*proto.ScenarioHistoryDB, 0)
	rsp.ScenarioCollectionDBs = make([]*proto.ScenarioCollectionDB, 0)

	for _, bin := range game.GetScenarioHistoryInfoList(s) {
		rsp.ScenarioHistoryDBs = append(rsp.ScenarioHistoryDBs, &proto.ScenarioHistoryDB{
			ScenarioUniqueId: bin.ScenarioUniqueId,
			ClearDateTime:    mx.Unix(bin.ClearDateTime, 0),
		})
	}
}

func ScenarioGroupHistoryUpdate(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ScenarioGroupHistoryUpdateRequest)
	rsp := response.(*proto.ScenarioGroupHistoryUpdateResponse)

	game.FinishScenarioGroupHistoryInfo(s, req.ScenarioGroupUniqueId, req.ScenarioType, 0)

	bin := game.GetScenarioGroupHistoryInfo(s, req.ScenarioGroupUniqueId)
	if bin == nil {
		return
	}
	rsp.ScenarioGroupHistoryDB = &proto.ScenarioGroupHistoryDB{
		AccountServerId:       s.AccountServerId,
		ScenarioGroupUqniueId: bin.ScenarioGroupUqniueId,
		ScenarioType:          bin.ScenarioType,
		EventContentId:        bin.EventContentId,
		ClearDateTime:         mx.Unix(bin.ClearDateTime, 0),
		IsReturn:              false,
	}
}

func ScenarioClear(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ScenarioClearRequest)
	rsp := response.(*proto.ScenarioClearResponse)

	game.BattleCheck(s, req.BattleSummary)
	parcelResultList := game.FinishScenarioHistoryInfo(s, req.ScenarioId)
	rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResultList)
	bin := game.GetScenarioHistoryInfo(s, req.ScenarioId)
	if bin != nil {
		rsp.ScenarioHistoryDB = &proto.ScenarioHistoryDB{
			ScenarioUniqueId: bin.ScenarioUniqueId,
			ClearDateTime:    mx.Unix(bin.ClearDateTime, 0),
		}
	}
}

func ScenarioSelect(s *enter.Session, request, response mx.Message) {

}
