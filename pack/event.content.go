package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/protocol/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func EventContentAdventureList(s *enter.Session, request, response mx.Message) {
	// req := request.(*proto.EventContentAdventureListRequest)
	rsp := response.(*proto.EventContentAdventureListResponse)

	rsp.AlreadyReceiveRewardId = make([]int64, 0)
	rsp.EventContentBonusRewardDBs = make([]*proto.EventContentBonusRewardDB, 0)
	rsp.StageHistoryDBs = make([]*proto.CampaignStageHistoryDB, 0)
	rsp.StrategyObjecthistoryDBs = make([]*proto.StrategyObjectHistoryDB, 0)
	rsp.StagePoint = 0
}

func EventContentBoxGachaShopList(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.EventContentBoxGachaShopListRequest)
	rsp := response.(*proto.EventContentBoxGachaShopListResponse)

	rsp.BoxGachaDB = &proto.EventContentBoxGachaDB{
		EventContentId: req.EventContentId,
		Seed:           0,
		Round:          1, // 回合数 默认1
		PurchaseCount:  0, // 购买数量
	}
	rsp.BoxGachaGroupIdByCount = make(map[int64]int64)
}

func EventContentScenarioGroupHistoryUpdate(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.EventContentScenarioGroupHistoryUpdateRequest)
	rsp := response.(*proto.EventContentScenarioGroupHistoryUpdateResponse)

	defer func() {
		rsp.ScenarioGroupHistoryDBs = game.GetScenarioGroupHistoryDBs(s)
	}()

	game.FinishScenarioGroupHistoryInfo(s, req.ScenarioGroupUniqueId, req.ScenarioType, req.EventContentId)
	// 奖励

}
