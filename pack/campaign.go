package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/proto"
)

func CampaignList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.CampaignListResponse)

	rsp.StageHistoryDBs = make([]*proto.CampaignStageHistoryDB, 0)
}
