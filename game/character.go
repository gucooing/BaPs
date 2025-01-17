package game

import (
	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/proto"
)

func NewCharacter(s *enter.Session) *sro.CharacterBin {
	if s == nil {
		return nil
	}
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
			ServerId:               GetServerId(s),
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
		bin.CharacterBin = NewCharacter(s)
	}
	return bin.CharacterBin
}

func GetCharacterInfoList(s *enter.Session) map[int64]*sro.CharacterInfo {
	bin := GetCharacterBin(s)
	if bin.CharacterInfoList == nil {
		bin.CharacterInfoList = make(map[int64]*sro.CharacterInfo)
	}
	return bin.GetCharacterInfoList()
}

func GetCharacterInfoListByServerId(s *enter.Session) map[int64]*sro.CharacterInfo {
	list := make(map[int64]*sro.CharacterInfo)
	for _, v := range GetCharacterInfoList(s) {
		list[v.ServerId] = v
	}
	return list
}

func GetCharacterInfo(s *enter.Session, characterId int64) *sro.CharacterInfo {
	bin := GetCharacterInfoList(s)
	return bin[characterId]
}

func GetCharacterServerId(s *enter.Session, characterId int64) int64 {
	bin := GetCharacterInfo(s, characterId)
	if bin == nil {
		return 0
	}
	return bin.ServerId
}

func GetCharacterDBs(s *enter.Session) []*proto.CharacterDB {
	list := make([]*proto.CharacterDB, 0)
	for _, bin := range GetCharacterInfoList(s) {
		list = append(list, GetCharacterDB(s, bin.CharacterId))
	}

	return list
}

func GetCharacterDB(s *enter.Session, characterId int64) *proto.CharacterDB {
	bin := GetCharacterInfo(s, characterId)
	if bin == nil {
		return nil
	}
	info := &proto.CharacterDB{
		Type:                   proto.ParcelType_Character,
		ServerId:               bin.ServerId,
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

	return info
}

func AddCharacter(s *enter.Session, characterId int64) bool {
	list := GetCharacterInfoList(s)
	if _, ok := list[characterId]; ok {
		return false
	}
	conf := gdconf.GetCharacterExcel(characterId)
	if conf == nil {
		logger.Error("[UID:%v]未知的角色添加调用,characterId:%v", s.AccountServerId, characterId)
		return true
	}
	info := &sro.CharacterInfo{
		CharacterId:            characterId,
		Level:                  1,
		Exp:                    0,
		FavorRank:              1,
		FavorExp:               0,
		StarGrade:              conf.DefaultStarGrade,
		ExSkillLevel:           1,
		PassiveSkillLevel:      1,
		ExtraPassiveSkillLevel: 1,
		CommonSkillLevel:       1,
		LeaderSkillLevel:       1,
		EquipmentList:          NewCharacterEquipment(characterId),
		ServerId:               GetServerId(s),
	}
	list[characterId] = info
	return true
}

func ServerIdsToCharacterIds(s *enter.Session, serverIdList []int64) []int64 {
	list := make([]int64, 0)
	bin := GetCharacterInfoListByServerId(s)
	for _, serverId := range serverIdList {
		if db, ok := bin[serverId]; ok {
			list = append(list, db.CharacterId)
		}
	}

	return list
}

func ServerIdToCharacterId(s *enter.Session, serverId int64) int64 {
	bin := GetCharacterInfoListByServerId(s)
	if db, ok := bin[serverId]; ok {
		return db.CharacterId
	}

	return 0
}
