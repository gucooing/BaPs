package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/proto"
)

func MissionList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.MissionListResponse)

	rsp.ClearedOrignalMissionIds = make([]int64, 0) // 已清除任务
	rsp.MissionHistoryUniqueIds = make([]int64, 0)
	rsp.ProgressDBs = make([]*proto.MissionProgressDB, 0) // 进度数据
	rsp.DailySuddenMissionInfo = &proto.MissionInfo{}     // 每日任务
}

func GuideMissionSeasonList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.GuideMissionSeasonListResponse)

	rsp.GuideMissionSeasonDBs = make([]*proto.GuideMissionSeasonDB, 0)
}

func ToastList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.ToastListResponse)

	rsp.ToastDBs = make([]*proto.ToastDB, 0)
}

func AccountGetTutorial(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.AccountGetTutorialResponse)

	rsp.TutorialIds = make([]int64, 0)
	for id, ok := range game.GetTutorialList(s) {
		if ok {
			rsp.TutorialIds = append(rsp.TutorialIds, id)
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

}
