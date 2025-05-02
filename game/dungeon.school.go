package game

import (
	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/protocol/proto"
)

func GetSchoolDungeonCost(isDel bool, count int64) []*ParcelResult {
	if isDel {
		return []*ParcelResult{
			{
				ParcelType: proto.ParcelType_Currency,
				ParcelId:   int64(proto.CurrencyTypes_SchoolDungeonTotalTicket),
				Amount:     -1 * count,
			},
			{
				ParcelType: proto.ParcelType_Currency,
				ParcelId:   int64(proto.CurrencyTypes_ActionPoint),
				Amount:     -10 * count,
			},
		}
	} else {
		return []*ParcelResult{
			{
				ParcelType: proto.ParcelType_Currency,
				ParcelId:   int64(proto.CurrencyTypes_SchoolDungeonTotalTicket),
				Amount:     1 * count,
			},
			{
				ParcelType: proto.ParcelType_Currency,
				ParcelId:   int64(proto.CurrencyTypes_ActionPoint),
				Amount:     8 * count,
			},
		}
	}
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
		StageUniqueId: stageId,
		StarFlags:     make([]bool, 3),
		Star1Flag:     bin.IsWin,
		Star2Flag:     bin.IsSu,
		Star3Flag:     bin.IsTime,
		IsClearedEver: bin.IsWin,
	}
	info.StarFlags[0] = bin.IsWin
	info.StarFlags[1] = bin.IsSu
	info.StarFlags[2] = bin.IsTime
	return info
}
