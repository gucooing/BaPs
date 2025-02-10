package game

import (
	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func GetRaidMultiFloorBin(s *enter.Session) *sro.RaidMultiFloorBin {
	bin := GetPlayerBin(s)
	if bin == nil {
		return nil
	}
	if bin.RaidMultiFloorBin == nil {
		bin.RaidMultiFloorBin = &sro.RaidMultiFloorBin{}
	}
	return bin.RaidMultiFloorBin
}

func GetCurRaidMultiFloorInfo(s *enter.Session) *sro.RaidMultiFloorInfo {
	bin := GetRaidMultiFloorBin(s)
	if bin == nil {
		return nil
	}
	conf := gdconf.GetCurMultiFloorRaidSeasonManageExcel()
	if conf == nil {
		return nil
	}
	if bin.CurRaidMultiFloor == nil {
		bin.CurRaidMultiFloor = &sro.RaidMultiFloorInfo{
			SeasonId: conf.SeasonId,
		}
	}
	if bin.CurRaidMultiFloor.SeasonId != conf.SeasonId {
		bin.CurRaidMultiFloor = &sro.RaidMultiFloorInfo{
			SeasonId: conf.SeasonId,
		}
	}
	return bin.CurRaidMultiFloor
}

func GetMultiFloorRaidDBs(s *enter.Session) []*proto.MultiFloorRaidDB {
	list := make([]*proto.MultiFloorRaidDB, 0)
	list = append(list, GetMultiFloorRaidDB(s))

	return list
}

func GetMultiFloorRaidDB(s *enter.Session) *proto.MultiFloorRaidDB {
	bin := GetCurRaidMultiFloorInfo(s)
	if bin == nil {
		return nil
	}
	info := &proto.MultiFloorRaidDB{
		SeasonId:          bin.SeasonId,
		ClearedDifficulty: bin.ClearedDifficulty,
		LastClearDate:     mx.Unix(bin.LastClearDate, 0),
		ClearBattleFrame:  bin.Frame,
		RewardDifficulty:  bin.RewardDifficulty,
		LastRewardDate:    mx.Unix(bin.LastRewardDate, 0),
	}
	if info.ClearedDifficulty == 0 {
		info.ClearBattleFrame = -1
	}
	return info
}
