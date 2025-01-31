package game

import (
	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/protocol/proto"
)

func NewDungeonBin() *sro.DungeonBin {
	return &sro.DungeonBin{}
}

func GetSchoolDungeonCost(isDel bool, count int64) []*ParcelResult {
	if isDel {
		return []*ParcelResult{
			{
				ParcelType: proto.ParcelType_Currency,
				ParcelId:   proto.CurrencyTypes_SchoolDungeonTotalTicket,
				Amount:     -1 * count,
			},
			{
				ParcelType: proto.ParcelType_Currency,
				ParcelId:   proto.CurrencyTypes_ActionPoint,
				Amount:     -10 * count,
			},
		}
	} else {
		return []*ParcelResult{
			{
				ParcelType: proto.ParcelType_Currency,
				ParcelId:   proto.CurrencyTypes_SchoolDungeonTotalTicket,
				Amount:     1 * count,
			},
			{
				ParcelType: proto.ParcelType_Currency,
				ParcelId:   proto.CurrencyTypes_ActionPoint,
				Amount:     8 * count,
			},
		}
	}
}

func GetDungeonBin(s *enter.Session) *sro.DungeonBin {
	bin := GetPlayerBin(s)
	if bin == nil {
		return nil
	}
	if bin.DungeonBin == nil {
		bin.DungeonBin = NewDungeonBin()
	}
	return bin.DungeonBin
}

func GetWeekDungeonStageInfoList(s *enter.Session) map[int64]*sro.WeekDungeonStageInfo {
	bin := GetDungeonBin(s)
	if bin == nil {
		return nil
	}
	if bin.WeekDungeonStageHistory == nil {
		bin.WeekDungeonStageHistory = make(map[int64]*sro.WeekDungeonStageInfo)
	}
	return bin.WeekDungeonStageHistory
}

func GetWeekDungeonStageInfo(s *enter.Session, stageId int64) *sro.WeekDungeonStageInfo {
	bin := GetWeekDungeonStageInfoList(s)
	if bin == nil {
		return nil
	}
	if bin[stageId] == nil {
		bin[stageId] = &sro.WeekDungeonStageInfo{
			StageId:        stageId,
			StarGoalRecord: make(map[string]int64),
		}
	}
	return bin[stageId]
}

func GetWeekDungeonStageHistoryDB(s *enter.Session, stageId int64) *proto.WeekDungeonStageHistoryDB {
	bin := GetWeekDungeonStageInfo(s, stageId)
	if bin == nil {
		return nil
	}
	info := &proto.WeekDungeonStageHistoryDB{
		AccountServerId: s.AccountServerId,
		StageUniqueId:   bin.StageId,
		StarGoalRecord:  make(map[proto.StarGoalType]int64),
		IsCleardEver:    false,
	}
	for starGoalType, status := range bin.StarGoalRecord {
		info.StarGoalRecord[proto.StarGoalType(starGoalType)] = status
	}
	return info
}

func GetSchoolDungeonStageInfoList(s *enter.Session) map[int64]*sro.SchoolDungeonStageInfo {
	bin := GetDungeonBin(s)
	if bin == nil {
		return nil
	}
	if bin.SchoolDungeonStageHistory == nil {
		bin.SchoolDungeonStageHistory = make(map[int64]*sro.SchoolDungeonStageInfo)
	}
	return bin.SchoolDungeonStageHistory
}

func GetSchoolDungeonStageInfo(s *enter.Session, stageId int64) *sro.SchoolDungeonStageInfo {
	bin := GetSchoolDungeonStageInfoList(s)
	if bin == nil {
		return nil
	}
	if bin[stageId] == nil {
		bin[stageId] = &sro.SchoolDungeonStageInfo{
			StageId: stageId,
		}
	}
	return bin[stageId]
}

func GetSchoolDungeonStageHistoryDB(s *enter.Session, stageId int64) *proto.SchoolDungeonStageHistoryDB {
	bin := GetSchoolDungeonStageInfo(s, stageId)
	if bin == nil {
		return nil
	}
	info := &proto.SchoolDungeonStageHistoryDB{
		AccountServerId: s.AccountServerId,
		StageUniqueId:   stageId,
		StarFlags:       make([]bool, 3),
		Star1Flag:       bin.IsWin,
		Star2Flag:       bin.IsSu,
		Star3Flag:       bin.IsTime,
		IsClearedEver:   bin.IsWin,
	}
	info.StarFlags[0] = bin.IsWin
	info.StarFlags[1] = bin.IsSu
	info.StarFlags[2] = bin.IsTime
	return info
}

func BattleIsAllAlive(battleSummary *proto.BattleSummary) bool {
	if battleSummary == nil {
		return false
	}
	isSu := true
	for _, heroes := range battleSummary.Group01Summary.Heroes {
		if heroes.HPRateAfter == 0 {
			isSu = false
		}
	}
	return isSu
}

func BattleIsClearTimeInSec(battleSummary *proto.BattleSummary, realtime float64) bool {
	if battleSummary == nil {
		return false
	}
	return battleSummary.ElapsedRealtime < realtime
}
