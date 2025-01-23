package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func WeekDungeonList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.WeekDungeonListResponse)

	rsp.AdditionalStageIdList = make([]int64, 0)
	rsp.WeekDungeonStageHistoryDBList = make([]*proto.WeekDungeonStageHistoryDB, 0)
	for _, v := range game.GetWeekDungeonStageInfoList(s) {
		info := &proto.WeekDungeonStageHistoryDB{
			AccountServerId: s.AccountServerId,
			StageUniqueId:   v.StageId,
			StarGoalRecord:  make(map[proto.StarGoalType]int64),
			IsCleardEver:    false,
		}
		for starGoalType, status := range v.StarGoalRecord {
			info.StarGoalRecord[proto.StarGoalType(starGoalType)] = status
		}
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

	rsp.MissionProgressDBs = make([]*proto.MissionProgressDB, 0)
	battleSummary := req.Summary
	if battleSummary == nil {
		return
	}
	if battleSummary.EndType != proto.BattleEndType_Clear { // 战败返还
		rsp.ParcelResultDB = game.ParcelResultDB(s, []*game.ParcelResult{
			{
				ParcelType: proto.ParcelType_Currency,
				ParcelId:   proto.CurrencyTypes_SchoolDungeonTotalTicket,
				Amount:     1,
			},
		})
		return
	}
}

func SchoolDungeonList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.SchoolDungeonListResponse)

	rsp.SchoolDungeonStageHistoryDBList = make([]*proto.SchoolDungeonStageHistoryDB, 0)
	for _, bin := range game.GetSchoolDungeonStageInfoList(s) {
		info := &proto.SchoolDungeonStageHistoryDB{
			AccountServerId: s.AccountServerId,
			StarFlags:       make([]bool, 3),
		}
		info.StarFlags[0] = bin.IsWin
		info.StarFlags[1] = bin.IsSu
		info.StarFlags[2] = bin.IsTime
		rsp.SchoolDungeonStageHistoryDBList = append(rsp.SchoolDungeonStageHistoryDBList, info)
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

	rsp.ParcelResultDB = game.ParcelResultDB(s, []*game.ParcelResult{
		{
			ParcelType: proto.ParcelType_Currency,
			ParcelId:   proto.CurrencyTypes_SchoolDungeonTotalTicket,
			Amount:     -1,
		},
	})
}

func SchoolDungeonBattleResult(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.SchoolDungeonBattleResultRequest)
	rsp := response.(*proto.SchoolDungeonBattleResultResponse)

	rsp.MissionProgressDBs = make([]*proto.MissionProgressDB, 0)
	battleSummary := req.Summary
	if battleSummary == nil {
		return
	}
	if battleSummary.EndType != proto.BattleEndType_Clear { // 战败返还
		rsp.ParcelResultDB = game.ParcelResultDB(s, []*game.ParcelResult{
			{
				ParcelType: proto.ParcelType_Currency,
				ParcelId:   proto.CurrencyTypes_SchoolDungeonTotalTicket,
				Amount:     1,
			},
		})
		return
	}
	bin := game.GetSchoolDungeonStageInfo(s, req.StageUniqueId)
	if bin == nil {
		return
	}
	bin.IsWin = true
	bin.IsTime = alg.MaxBool(bin.IsTime, battleSummary.ElapsedRealtime < 120)
	isSu := true
	for _, heroes := range battleSummary.Group01Summary.Heroes {
		if heroes.HPRateAfter == 0 {
			isSu = false
		}
	}
	bin.IsSu = alg.MaxBool(bin.IsSu, isSu)
	// 发奖励！

	// 更新角色

	// 发固定

	// 发首通

	// 发三星

	rsp.SchoolDungeonStageHistoryDB = game.GetSchoolDungeonStageHistoryDB(s, req.StageUniqueId)
}
