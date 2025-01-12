package pack

import (
	"github.com/gucooing/BaPs/common/enter"
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

}

func ToastList(s *enter.Session, request, response mx.Message) {

}
