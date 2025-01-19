package pack

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func MissionList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.MissionListResponse)

	rsp.ClearedOrignalMissionIds = make([]int64, 0) // 已清除任务
	rsp.MissionHistoryUniqueIds = make([]int64, 0)
	rsp.ProgressDBs = game.GetProgressDBs(s)                         // 进度数据
	rsp.DailySuddenMissionInfo = game.GetDailySuddenMissionInfoDb(s) // 每日任务
}

func GuideMissionSeasonList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.GuideMissionSeasonListResponse)

	rsp.MissionProgressDBs = make([]*proto.MissionProgressDB, 0)
	rsp.GuideMissionSeasonDBs = []*proto.GuideMissionSeasonDB{
		{
			SeasonId:                  3,
			LoginCount:                1,
			StartDate:                 time.Now().Add(-time.Hour * 24),
			LoginDate:                 time.Now(),
			IsComplete:                false,
			IsFinalMissionComplete:    false,
			CollectionItemReceiveDate: time.Time{},
		},
		{
			SeasonId:                  1000,
			LoginCount:                1,
			StartDate:                 time.Now().Add(-time.Hour * 24),
			LoginDate:                 time.Now(),
			IsComplete:                false,
			IsFinalMissionComplete:    false,
			CollectionItemReceiveDate: time.Time{},
		},
		{
			SeasonId: 1001,
		},
	}
}

func ToastList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.ToastListResponse)

	rsp.ToastDBs = make([]*proto.ToastDB, 0)
}

func AccountGetTutorial(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.AccountGetTutorialResponse)

	rsp.TutorialIds = make([]int64, 0)
	// for id, ok := range game.GetTutorialList(s) {
	// 	if ok {
	// 		rsp.TutorialIds = append(rsp.TutorialIds, id)
	// 	}
	// }
	for i := 1; i < 28; i++ {
		rsp.TutorialIds = append(rsp.TutorialIds, int64(i))
	}
}

func ScenarioSkip(s *enter.Session, request, response mx.Message) {

}

func AccountSetTutorial(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.AccountSetTutorialRequest)

	game.FinishTutorial(s, req.TutorialIds)
}

func MissionSync(s *enter.Session, request, response mx.Message) {

}

func ScenarioList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.ScenarioListResponse)

	rsp.ScenarioGroupHistoryDBs = make([]*proto.ScenarioGroupHistoryDB, 0)
	rsp.ScenarioHistoryDBs = make([]*proto.ScenarioHistoryDB, 0)

	rsp.ScenarioCollectionDBs = make([]*proto.ScenarioCollectionDB, 0)

	for _, bin := range game.GetScenarioGroupHistoryInfoList(s) {
		rsp.ScenarioGroupHistoryDBs = append(rsp.ScenarioGroupHistoryDBs, &proto.ScenarioGroupHistoryDB{
			AccountServerId:       s.AccountServerId,
			ScenarioGroupUqniueId: bin.ScenarioGroupUqniueId,
			ScenarioType:          bin.ScenarioType,
			EventContentId:        bin.EventContentId,
			ClearDateTime:         time.Unix(bin.ClearDateTime, 0),
			IsReturn:              false,
		})
	}

	for _, bin := range game.GetScenarioHistoryInfoList(s) {
		rsp.ScenarioHistoryDBs = append(rsp.ScenarioHistoryDBs, &proto.ScenarioHistoryDB{
			AccountServerId:  s.AccountServerId,
			ScenarioUniqueId: bin.ScenarioUniqueId,
			ClearDateTime:    time.Unix(bin.ClearDateTime, 0),
		})
	}
}

func ScenarioGroupHistoryUpdate(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.ScenarioGroupHistoryUpdateRequest)
	rsp := response.(*proto.ScenarioGroupHistoryUpdateResponse)

	game.FinishScenarioGroupHistoryInfo(s, req.ScenarioGroupUniqueId, req.ScenarioType)

	bin := game.GetScenarioGroupHistoryInfo(s, req.ScenarioGroupUniqueId)
	if bin == nil {
		return
	}
	rsp.ScenarioGroupHistoryDB = &proto.ScenarioGroupHistoryDB{
		AccountServerId:       s.AccountServerId,
		ScenarioGroupUqniueId: bin.ScenarioGroupUqniueId,
		ScenarioType:          bin.ScenarioType,
		EventContentId:        bin.EventContentId,
		ClearDateTime:         time.Unix(bin.ClearDateTime, 0),
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
			AccountServerId:  s.AccountServerId,
			ScenarioUniqueId: bin.ScenarioUniqueId,
			ClearDateTime:    time.Unix(bin.ClearDateTime, 0),
		}
	}
}

func ScenarioSelect(s *enter.Session, request, response mx.Message) {

}
