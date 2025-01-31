package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func WeekDungeonList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.WeekDungeonListResponse)

	rsp.AdditionalStageIdList = make([]int64, 0)
	rsp.WeekDungeonStageHistoryDBList = make([]*proto.WeekDungeonStageHistoryDB, 0)
	for _, v := range game.GetWeekDungeonStageInfoList(s) {
		rsp.WeekDungeonStageHistoryDBList = append(
			rsp.WeekDungeonStageHistoryDBList, game.GetWeekDungeonStageHistoryDB(s, v.StageId))
	}
}

func WeekDungeonEnterBattle(s *enter.Session, request, response mx.Message) {
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

func WeekDungeonBattleResult(s *enter.Session, request, response mx.Message) {
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

func SchoolDungeonList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.SchoolDungeonListResponse)

	rsp.SchoolDungeonStageHistoryDBList = make([]*proto.SchoolDungeonStageHistoryDB, 0)
	for _, bin := range game.GetSchoolDungeonStageInfoList(s) {
		rsp.SchoolDungeonStageHistoryDBList = append(rsp.SchoolDungeonStageHistoryDBList,
			game.GetSchoolDungeonStageHistoryDB(s, bin.StageId))
	}
}

func SchoolDungeonEnterBattle(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.SchoolDungeonEnterBattleRequest)
	rsp := response.(*proto.SchoolDungeonEnterBattleResponse)

	conf := gdconf.GetSchoolDungeonStageExcel(req.StageUniqueId)
	if conf == nil {
		logger.Debug("不存在该学院交流会关卡 StageUniqueId:", req.StageUniqueId)
		return
	}

	rsp.ParcelResultDB = game.ParcelResultDB(s, game.GetSchoolDungeonCost(true, 1))
}

func SchoolDungeonBattleResult(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.SchoolDungeonBattleResultRequest)
	rsp := response.(*proto.SchoolDungeonBattleResultResponse)

	rsp.FirstClearReward = make([]*proto.ParcelInfo, 0)
	rsp.ThreeStarReward = make([]*proto.ParcelInfo, 0)
	rsp.MissionProgressDBs = make([]*proto.MissionProgressDB, 0)
	battleSummary := req.Summary
	if battleSummary == nil {
		return
	}
	if battleSummary.EndType != proto.BattleEndType_Clear { // 战败返还
		rsp.ParcelResultDB = game.ParcelResultDB(s, game.GetSchoolDungeonCost(false, 1))
		return
	}
	bin := game.GetSchoolDungeonStageInfo(s, req.StageUniqueId)
	if bin == nil {
		return
	}
	clearTime := game.BattleIsClearTimeInSec(req.Summary, 120)
	allAlive := game.BattleIsAllAlive(req.Summary)

	isFirstClear := !bin.IsWin // 是否第一次通过
	isThreeStar := (clearTime && allAlive) && (!bin.IsTime || !bin.IsSu)
	bin.IsWin = true
	bin.IsTime = bin.IsTime || clearTime
	bin.IsSu = bin.IsSu || allAlive
	conf := gdconf.GetSchoolDungeonStageExcel(req.StageUniqueId)
	if conf == nil {
		return
	}
	// 发奖励！
	parcelResultList, _ := game.ContentSweepSchoolDungeon(req.StageUniqueId, 1)
	// 星级判断
	for _, rewardConf := range gdconf.GetSchoolDungeonRewardExcelList(conf.StageRewardId) {
		if !rewardConf.IsDisplayed {
			continue
		}
		if (rewardConf.RewardTag == "ThreeStar" && isThreeStar) ||
			(rewardConf.RewardTag == "FirstClear" && isFirstClear) {
			parcelType := proto.GetParcelTypeValue(rewardConf.RewardParcelType)
			parcelResultList = append(parcelResultList, &game.ParcelResult{
				ParcelType: parcelType,
				ParcelId:   rewardConf.RewardParcelId,
				Amount:     rewardConf.RewardParcelAmount,
			})
			parcelInfo := game.GetParcelInfo(rewardConf.RewardParcelId,
				rewardConf.RewardParcelAmount, parcelType)
			if rewardConf.RewardTag == "ThreeStar" {
				rsp.ThreeStarReward = append(rsp.ThreeStarReward, parcelInfo)
			} else if rewardConf.RewardTag == "FirstClear" {
				rsp.FirstClearReward = append(rsp.FirstClearReward, parcelInfo)
			}
		}
	}

	rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResultList)
	rsp.SchoolDungeonStageHistoryDB = game.GetSchoolDungeonStageHistoryDB(s, req.StageUniqueId)
}
