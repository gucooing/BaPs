package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/protocol/proto"
)

func WeekDungeonList(s *enter.Session, request, response proto.Message) {
	rsp := response.(*proto.WeekDungeonListResponse)

	rsp.AdditionalStageIdList = make([]int64, 0)
	rsp.WeekDungeonStageHistoryDBList = make([]*proto.WeekDungeonStageHistoryDB, 0)
	for _, v := range game.GetWeekDungeonStageInfoList(s) {
		rsp.WeekDungeonStageHistoryDBList = append(
			rsp.WeekDungeonStageHistoryDBList, game.GetWeekDungeonStageHistoryDB(s, v.StageId))
	}
}

func WeekDungeonEnterBattle(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.WeekDungeonEnterBattleRequest)
	rsp := response.(*proto.WeekDungeonEnterBattleResponse)

	conf := gdconf.GetWeekDungeonExcelTable(req.StageUniqueId)
	if conf == nil {
		return
	}
	parcelResult := game.GetParcelResultList(conf.StageEnterCostType,
		conf.StageEnterCostId, conf.StageEnterCostAmount, true)
	rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResult)
}

func WeekDungeonBattleResult(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.WeekDungeonBattleResultRequest)
	rsp := response.(*proto.WeekDungeonBattleResultResponse)

	defer func() {
		rsp.WeekDungeonStageHistoryDB = game.GetWeekDungeonStageHistoryDB(s, req.StageUniqueId)
	}()

	rsp.MissionProgressDBs = make([]*proto.MissionProgressDB, 0)
	battleSummary := req.Summary
	if battleSummary == nil {
		return
	}
	conf := gdconf.GetWeekDungeonExcelTable(req.StageUniqueId)
	if conf == nil {
		return
	}
	if battleSummary.EndType != proto.BattleEndType_Clear { // 战败返还 100%
		parcelResult := game.GetParcelResultList(conf.StageEnterCostType,
			conf.StageEnterCostId, conf.StageEnterCostAmount, false)
		rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResult)
		return
	}
	bin := game.GetWeekDungeonStageInfo(s, req.StageUniqueId)
	if bin == nil {
		return
	}
	if bin.StarGoalRecord == nil {
		bin.StarGoalRecord = make(map[string]int64)
	}
	// 计算得分
	for index, v := range conf.StarGoal {
		status := bin.StarGoalRecord[v]
		switch v {
		case "Clear":
			status = 1
		case "AllAlive":
			if game.BattleIsAllAlive(req.Summary) {
				status = 1
			}
		case "ClearTimeInSec":
			if game.BattleIsClearTimeInSec(req.Summary,
				conf.StarGoalAmount[index]) {
				status = 1
			}
		}
		bin.StarGoalRecord[v] = status
	}
	// 发奖励！
	parcelResultList, _ := game.ContentSweepWeekDungeon(req.StageUniqueId, 1)
	rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResultList)
}
