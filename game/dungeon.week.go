package game

import (
	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/protocol/proto"
)

func NewDungeonBin() *sro.DungeonBin {
	return &sro.DungeonBin{}
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
		info.StarGoalRecord[proto.StarGoalType_None.Value(starGoalType)] = status
	}
	return info
}
