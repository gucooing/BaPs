package game

import (
	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/mx/proto"
)

func NewCharacter() *sro.CharacterBin {
	bin := &sro.CharacterBin{
		CharacterInfoList: make(map[int64]*sro.CharacterInfo),
	}
	for _, conf := range gdconf.GetDefaultCharacterExcelTable() {
		infp := &sro.CharacterInfo{
			CharacterId:            conf.CharacterId,
			Level:                  conf.Level,
			Exp:                    conf.Exp,
			FavorRank:              conf.FavorRank,
			FavorExp:               conf.FavorExp,
			StarGrade:              conf.StarGrade,
			ExSkillLevel:           conf.ExSkillLevel,
			PassiveSkillLevel:      conf.PassiveSkillLevel,
			ExtraPassiveSkillLevel: conf.ExtraPassiveSkillLevel,
			CommonSkillLevel:       conf.CommonSkillLevel,
			LeaderSkillLevel:       conf.LeaderSkillLevel,
			EquipmentList:          NewCharacterEquipment(conf.CharacterId),
		}
		bin.CharacterInfoList[conf.CharacterId] = infp
	}

	return bin
}

func NewCharacterEquipment(characterId int64) map[string]int64 {
	conf := gdconf.GetCharacterExcel(characterId)
	if conf == nil {
		return nil
	}
	list := make(map[string]int64)
	for _, e := range conf.EquipmentSlot {
		list[e] = 0
	}
	return list
}

func GetCharacterBin(s *enter.Session) *sro.CharacterBin {
	bin := GetPlayerBin(s)
	if bin == nil {
		return nil
	}
	if bin.CharacterBin == nil {
		bin.CharacterBin = NewCharacter()
	}
	return bin.CharacterBin
}

func GetCharacterInfoList(s *enter.Session) map[int64]*sro.CharacterInfo {
	bin := GetCharacterBin(s)
	return bin.GetCharacterInfoList()
}

func GetCharacterDBs(s *enter.Session) []*proto.CharacterDB {
	list := make([]*proto.CharacterDB, 0)
	for _, bin := range GetCharacterInfoList(s) {
		info := &proto.CharacterDB{
			Type:                   proto.ParcelType_Character,
			ServerId:               GetServerId(),
			UniqueId:               bin.CharacterId,
			StarGrade:              bin.StarGrade,
			Level:                  bin.Level,
			Exp:                    bin.Exp,
			FavorRank:              bin.FavorRank,
			FavorExp:               bin.FavorExp,
			PublicSkillLevel:       bin.CommonSkillLevel,
			ExSkillLevel:           bin.ExSkillLevel,
			PassiveSkillLevel:      bin.PassiveSkillLevel,
			ExtraPassiveSkillLevel: bin.ExtraPassiveSkillLevel,
			LeaderSkillLevel:       bin.LeaderSkillLevel,
			IsNew:                  false,
			IsLocked:               false,
			IsFavorite:             false,
			EquipmentServerIds:     make([]int64, 0),
			PotentialStats:         make(map[int32]int32),
		}
		for _, e := range bin.EquipmentList {
			info.EquipmentServerIds = append(info.EquipmentServerIds, e)
		}
		for _, p := range []int32{1, 2, 3} {
			info.PotentialStats[p] = 0
		}

		list = append(list, info)
	}

	return list
}
