package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/protocol/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func CampaignList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.CampaignListResponse)

	rsp.CampaignChapterClearRewardHistoryDBs = make([]*proto.CampaignChapterClearRewardHistoryDB, 0)
	rsp.StageHistoryDBs = make([]*proto.CampaignStageHistoryDB, 0)

	cidList := make(map[int64]bool)
	for _, conf := range gdconf.GetCampaignStageExcelMap() {
		cid := conf.Id / 10000
		if _, ok := cidList[cid]; !ok {
			rsp.CampaignChapterClearRewardHistoryDBs = append(rsp.CampaignChapterClearRewardHistoryDBs, &proto.CampaignChapterClearRewardHistoryDB{
				ChapterUniqueId: cid,
				RewardType:      proto.StageDifficulty_Normal,
			})
			rsp.CampaignChapterClearRewardHistoryDBs = append(rsp.CampaignChapterClearRewardHistoryDBs, &proto.CampaignChapterClearRewardHistoryDB{
				ChapterUniqueId: cid,
				RewardType:      proto.StageDifficulty_Hard,
			})
			cidList[cid] = true
		}
		rsp.StageHistoryDBs = append(rsp.StageHistoryDBs, &proto.CampaignStageHistoryDB{
			ChapterUniqueId:                 cid,
			StageUniqueId:                   conf.Id,
			ClearTurnRecord:                 2,
			TacticClearCountWithRankSRecord: 2,
			Star1Flag:                       true,
			Star2Flag:                       true,
			Star3Flag:                       true,
			LastPlay:                        mx.Now(),
			FirstClearRewardReceive:         mx.Now(),
			StarRewardReceive:               mx.Now(),
			TodayPlayCount:                  1,

			TodayPurchasePlayCountHardStage: 0,
			TodayPlayCountForUI:             0,
		})
	}
}

func CampaignEnterMainStage(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.CampaignEnterMainStageRequest)
	rsp := response.(*proto.CampaignEnterMainStageResponse)

	bin := game.NewCampaignMainStageSaveDB(s, req.StageUniqueId)
	rsp.SaveDataDB = game.GetCampaignMainStageSaveDB(s, bin)
}

func CampaignChapterClearReward(s *enter.Session, request, response mx.Message) {

}

func CampaignEnterMainStageStrategySkip(s *enter.Session, request, response mx.Message) {
	// req := request.(*proto.CampaignEnterMainStageStrategySkipRequest)
	rsp := response.(*proto.CampaignEnterMainStageStrategySkipResponse)

	// 成本计算
	rsp.ParcelResultDB = new(proto.ParcelResultDB)
}

func CampaignMainStageStrategySkipResult(s *enter.Session, request, response mx.Message) {
	// req := request.(*proto.CampaignMainStageStrategySkipResultRequest)
	// rsp := response.(*proto.CampaignMainStageStrategySkipResultResponse)
	//
	// win := req.Summary.Winner == "Group01"

}
