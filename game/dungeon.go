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
		StarFlags:       make([]bool, 3),
	}
	info.StarFlags[0] = bin.IsWin
	info.StarFlags[1] = bin.IsSu
	info.StarFlags[2] = bin.IsTime
	return info
}
