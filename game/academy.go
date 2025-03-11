package game

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func GetAcademyBin(s *enter.Session) *sro.AcademyBin {
	bin := GetPlayerBin(s)
	if bin == nil {
		return nil
	}
	if bin.AcademyBin == nil {
		bin.AcademyBin = &sro.AcademyBin{}
	}
	return bin.AcademyBin
}

func GetAcademyLocationInfoList(s *enter.Session) map[int64]*sro.AcademyLocationInfo {
	bin := GetAcademyBin(s)
	if bin == nil {
		return nil
	}
	if bin.AcademyLocationList == nil {
		bin.AcademyLocationList = make(map[int64]*sro.AcademyLocationInfo)
	}
	for _, conf := range gdconf.GetAcademyLocationExcelTableList() {
		if bin.AcademyLocationList[conf.Id] == nil {
			bin.AcademyLocationList[conf.Id] = &sro.AcademyLocationInfo{
				LocationId: conf.Id,
				Rank:       1,
				Exp:        0,
			}
		}
	}
	return bin.AcademyLocationList
}

func GetAcademyLocationInfo(s *enter.Session, locationId int64) *sro.AcademyLocationInfo {
	bin := GetAcademyLocationInfoList(s)
	if bin == nil {
		return nil
	}
	return bin[locationId]
}

func GetAcademyLocationRankSum(s *enter.Session) int64 {
	level := int64(0)
	for _, conf := range GetAcademyLocationInfoList(s) {
		level += conf.Rank
	}
	return level
}

func GetMaxAcademyTicket(s *enter.Session) int64 {
	return gdconf.GetScheduleTicktetMax(GetAcademyLocationRankSum(s))
}

func GetAcademyZoneInfoList(s *enter.Session) map[int64]*sro.AcademyZoneInfo {
	bin := GetAcademyBin(s)
	if bin == nil {
		return nil
	}
	if bin.AcademyZoneList == nil {
		bin.AcademyZoneList = make(map[int64]*sro.AcademyZoneInfo)
	}
	if !time.Unix(bin.LastUpData, 0).After(alg.GetDay4()) {
		bin.AcademyZoneList = UpAcademyZoneInfoList(s)
		bin.LastUpData = time.Now().Unix()
	}
	return bin.AcademyZoneList
}

func UpAcademyZoneInfoList(s *enter.Session) map[int64]*sro.AcademyZoneInfo {
	list := make(map[int64]*sro.AcademyZoneInfo)
	for _, conf := range gdconf.GetAcademyZoneExcelTableList() {
		info := &sro.AcademyZoneInfo{
			ZoneId:      conf.Id,
			StudentList: make([]int64, 0),
			IsUp:        false,
		}
		for i := 0; i < len(conf.StudentVisitProb); i++ {
			studentId := gdconf.RandCharacter()
			for _, cid := range info.StudentList {
				if cid == studentId {
					i--
					continue
				}
			}
			info.StudentList = append(info.StudentList, studentId)
		}

		list[conf.Id] = info
	}
	return list
}

func GetAcademyZoneInfo(s *enter.Session, zoneId int64) *sro.AcademyZoneInfo {
	bin := GetAcademyZoneInfoList(s)
	if bin == nil {
		return nil
	}
	return bin[zoneId]
}

func UpAcademyLocationExp(s *enter.Session, locationId int64, exp int64) {
	bin := GetAcademyLocationInfo(s, locationId)
	if bin == nil {
		return
	}
	bin.Exp += exp
	// 升级判断
	for {
		conf := gdconf.GetAcademyLocationRankExcelTable(bin.Rank)
		if conf == nil {
			bin.Rank--
			return
		}
		if bin.Exp < conf.RankExp {
			return
		}
		bin.Exp -= conf.RankExp
		bin.Rank++
	}
}

func GetAcademyDB(s *enter.Session) *proto.AcademyDB {
	info := &proto.AcademyDB{
		LastUpdate:               mx.MxTime{},
		ZoneVisitCharacterDBs:    make(map[int64][]*proto.VisitingCharacterDB),
		ZoneScheduleGroupRecords: make(map[int64][]int64),
	}
	for _, bin := range GetAcademyZoneInfoList(s) {
		if bin.IsUp {
			conf := gdconf.GetAcademyZoneExcelTable(bin.ZoneId)
			if conf == nil {
				continue
			}
			rewardGroupIdList := make([]int64, 0)
			rewardGroupIdList = append(rewardGroupIdList, conf.RewardGroupId)
			info.ZoneScheduleGroupRecords[bin.ZoneId] = rewardGroupIdList
		}
		visit := make([]*proto.VisitingCharacterDB, 0)
		for _, characterId := range bin.StudentList {
			visit = append(visit, &proto.VisitingCharacterDB{
				UniqueId: characterId,
				ServerId: GetCharacterServerId(s, characterId),
			})
		}
		info.ZoneVisitCharacterDBs[bin.ZoneId] = visit
	}
	info.LastUpdate = mx.Unix(GetAcademyBin(s).GetLastUpData(), 0)

	return info
}

func GetAcademyLocationDBs(s *enter.Session) []*proto.AcademyLocationDB {
	list := make([]*proto.AcademyLocationDB, 0)
	for _, info := range GetAcademyLocationInfoList(s) {
		list = append(list, &proto.AcademyLocationDB{
			LocationId: info.LocationId,
			Rank:       info.Rank,
			Exp:        info.Exp,
		})
	}
	return list
}

func GetAcademyLocationDB(s *enter.Session, locationId int64) *proto.AcademyLocationDB {
	bin := GetAcademyLocationInfo(s, locationId)
	if bin == nil {
		return nil
	}
	return &proto.AcademyLocationDB{
		LocationId: bin.LocationId,
		Rank:       bin.Rank,
		Exp:        bin.Exp,
	}
}
