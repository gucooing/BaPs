package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/proto"
)

func MomoTalkOutLine(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.MomoTalkOutLineResponse)

	rsp.MomoTalkOutLineDBs = make([]*proto.MomoTalkOutLineDB, 0)
	rsp.FavorScheduleRecords = make(map[int64][]int64)
}

func ScenarioList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.ScenarioListResponse)

	rsp.ScenarioGroupHistoryDBs = make([]*proto.ScenarioGroupHistoryDB, 0)
}
