package game

import (
	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

const (
	MaxMainEchelonNum    int32 = 4 // 最大队伍主角色数量
	MaxSupportEchelonNum int32 = 2 // 最大队伍支援角色数量
	MaxSkillEchelonNum   int32 = 3 // 最大队伍技能优先释放数量
	MaxCombatStyleNum    int32 = 6 // 战斗风格指数?
)

var MainEchelonNum = map[proto.EchelonType]int32{
	proto.EchelonType_MultiFloorRaid: 6,
}

var SupportEchelonNum = map[proto.EchelonType]int32{
	proto.EchelonType_MultiFloorRaid: 4,
}

var SkillEchelonNum = map[proto.EchelonType]int32{
	proto.EchelonType_MultiFloorRaid: 5,
}

var CombatStyleNum = map[proto.EchelonType]int32{
	proto.EchelonType_MultiFloorRaid: 10,
}

func GetMaxMainEchelonNum(echelonType proto.EchelonType) int32 {
	if v, ok := MainEchelonNum[echelonType]; ok {
		return v
	}
	return MaxMainEchelonNum
}

func GetSupportEchelonNum(echelonType proto.EchelonType) int32 {
	if v, ok := SupportEchelonNum[echelonType]; ok {
		return v
	}
	return MaxSupportEchelonNum
}

func GetSkillEchelonNum(skillType proto.EchelonType) int32 {
	if v, ok := SkillEchelonNum[skillType]; ok {
		return v
	}
	return MaxSkillEchelonNum
}

func GetCombatStyleNum(skillType proto.EchelonType) int32 {
	if v, ok := CombatStyleNum[skillType]; ok {
		return v
	}
	return MaxCombatStyleNum
}

func NewEchelonTypeInfoList() map[int32]*sro.EchelonTypeInfo {
	list := make(map[int32]*sro.EchelonTypeInfo)
	for _, conf := range gdconf.GetDefaultEchelonExcelList() {
		if list[conf.EchlonId] == nil {
			list[conf.EchlonId] = &sro.EchelonTypeInfo{
				EchelonNum:      1,
				EchelonInfoList: make(map[int64]*sro.EchelonInfo),
			}
		}
		if list[conf.EchlonId].EchelonInfoList == nil {
			list[conf.EchlonId].EchelonInfoList = make(map[int64]*sro.EchelonInfo)
		}
		UpEchelonInfo(list[conf.EchlonId], conf, list[conf.EchlonId].EchelonNum)
	}
	return list
}

// 深度更新
func UpEchelonInfo(typeInfo *sro.EchelonTypeInfo, conf *sro.DefaultEchelonExcelTable, num int64) *sro.EchelonInfo {
	if typeInfo == nil || conf == nil {
		return nil
	}
	if typeInfo.EchelonInfoList == nil {
		typeInfo.EchelonInfoList = make(map[int64]*sro.EchelonInfo)
	}

	info := &sro.EchelonInfo{
		EchelonType:          conf.EchlonId,
		ExtensionType:        conf.ExtensionType,
		EchelonNum:           num,
		LeaderCharacter:      conf.LeaderId,
		TssId:                conf.TssId,
		MainCharacterList:    make(map[int32]int64),
		SupportCharacterList: make(map[int32]int64),
		SkillCharacterList:   make(map[int32]int64),
	}
	typeInfo.EchelonInfoList[num] = info
	if typeInfo.EchelonNum <= num {
		typeInfo.EchelonNum = num + 1
	}

	var i int32 = 1
	for ; i <= GetMaxMainEchelonNum(proto.EchelonType(conf.EchlonId)); i++ {
		if info.MainCharacterList == nil {
			info.MainCharacterList = make(map[int32]int64)
		}
		if len(conf.MainId) < int(i) {
			info.MainCharacterList[i] = 0
		} else {
			info.MainCharacterList[i] = conf.MainId[i-1]
		}
	}
	i = 1
	for ; i <= GetSupportEchelonNum(proto.EchelonType(conf.EchlonId)); i++ {
		if info.SupportCharacterList == nil {
			info.SupportCharacterList = make(map[int32]int64)
		}
		if len(conf.SupportId) < int(i) {
			info.SupportCharacterList[i] = 0
		} else {
			info.SupportCharacterList[i] = conf.SupportId[i-1]
		}
	}
	i = 1
	for ; i <= GetSkillEchelonNum(proto.EchelonType(conf.EchlonId)); i++ {
		if info.SkillCharacterList == nil {
			info.SkillCharacterList = make(map[int32]int64)
		}
		if len(conf.SkillId) < int(i) {
			info.SkillCharacterList[i] = 0
		} else {
			info.SkillCharacterList[i] = conf.SkillId[i-1]
		}
	}
	return info
}

func GetEchelonBin(s *enter.Session) *sro.EchelonBin {
	bin := GetPlayerBin(s)
	if bin == nil {
		return nil
	}
	if bin.EchelonBin == nil {
		bin.EchelonBin = &sro.EchelonBin{}
	}
	return bin.EchelonBin
}

func GetEchelonTypeInfoList(s *enter.Session) map[int32]*sro.EchelonTypeInfo {
	bin := GetEchelonBin(s)
	if bin == nil {
		return nil
	}
	if bin.EchelonTypeInfoList == nil {
		bin.EchelonTypeInfoList = NewEchelonTypeInfoList()
	}
	return bin.EchelonTypeInfoList
}

func GetEchelonTypeInfo(s *enter.Session, echelonType int32) *sro.EchelonTypeInfo {
	bin := GetEchelonTypeInfoList(s)
	if bin == nil {
		return nil
	}
	if bin[echelonType] == nil {
		bin[echelonType] = &sro.EchelonTypeInfo{
			EchelonInfoList: make(map[int64]*sro.EchelonInfo),
			EchelonNum:      1,
		}
	}
	return bin[echelonType]
}

func GetEchelonInfo(s *enter.Session, echelonType int32, num int64) *sro.EchelonInfo {
	bin := GetEchelonTypeInfo(s, echelonType)
	if bin == nil {
		return nil
	}
	if bin.EchelonInfoList == nil {
		bin.EchelonInfoList = make(map[int64]*sro.EchelonInfo)
	}
	if bin.EchelonInfoList[num] == nil {
		bin.EchelonInfoList[num] = &sro.EchelonInfo{
			EchelonType:          echelonType,
			ExtensionType:        proto.EchelonExtensionType_Base,
			EchelonNum:           num,
			LeaderCharacter:      0,
			MainCharacterList:    make(map[int32]int64),
			SupportCharacterList: make(map[int32]int64),
			SkillCharacterList:   make(map[int32]int64),
			TssId:                0,
		}
	}
	if bin.EchelonNum <= num {
		bin.EchelonNum++
	}
	return bin.EchelonInfoList[num]
}

func NewEchelonPresetGuidList() map[int32]*sro.EchelonTypeInfo {
	list := make(map[int32]*sro.EchelonTypeInfo)
	for _, gid := range []int32{0, 1, 2, 3} {
		typeInfo := &sro.EchelonTypeInfo{
			EchelonInfoList: make(map[int64]*sro.EchelonInfo),
			EchelonNum:      4,
		}
		for _, index := range []int64{0, 1, 2, 3, 4} {
			typeInfo.EchelonInfoList[index] = &sro.EchelonInfo{
				EchelonType:          gid,
				ExtensionType:        0,
				EchelonNum:           index,
				LeaderCharacter:      0,
				MainCharacterList:    make(map[int32]int64),
				SupportCharacterList: make(map[int32]int64),
				SkillCharacterList:   make(map[int32]int64),
				TssId:                0,
			}
		}
		list[gid] = typeInfo
	}
	return list
}

func GetEchelonPresetGuidList(s *enter.Session) map[int32]*sro.EchelonTypeInfo {
	bin := GetEchelonBin(s)
	if bin == nil {
		return nil
	}
	if bin.EchelonPresetGuidList == nil {
		bin.EchelonPresetGuidList = NewEchelonPresetGuidList()
	}
	return bin.EchelonPresetGuidList
}

func GetEchelonDB(s *enter.Session, db *sro.EchelonInfo) *proto.EchelonDB {
	if db == nil {
		return nil
	}
	info := &proto.EchelonDB{
		AccountServerId:               s.AccountServerId,
		EchelonType:                   proto.EchelonType(db.EchelonType),
		EchelonNumber:                 db.EchelonNum,
		ExtensionType:                 proto.EchelonExtensionType(db.ExtensionType),
		LeaderServerId:                0,
		MainSlotServerIds:             make([]int64, 0),
		SupportSlotServerIds:          make([]int64, 0),
		TSSInteractionServerId:        db.TssId,
		UsingFlag:                     0,
		SkillCardMulliganCharacterIds: make([]int64, 0),
		CombatStyleIndex:              make([]int, 0),
	}
	if characterInfo := GetCharacterInfo(s, db.LeaderCharacter); characterInfo != nil {
		info.LeaderServerId = characterInfo.ServerId
	} else {
		logger.Debug("[UID:%v]玩家队伍队长为空")
	}
	var i int32 = 1
	for ; i <= GetMaxMainEchelonNum(proto.EchelonType(db.EchelonType)); i++ {
		var serverId int64 = 0
		characterId, ok := db.MainCharacterList[i]
		if ok {
			characterInfo := GetCharacterInfo(s, characterId)
			if characterInfo != nil {
				serverId = characterInfo.ServerId
			}
		}
		info.MainSlotServerIds = append(info.MainSlotServerIds, serverId)
	}
	i = 1
	for ; i <= GetSupportEchelonNum(proto.EchelonType(db.EchelonType)); i++ {
		var serverId int64 = 0
		characterId, ok := db.SupportCharacterList[i]
		if ok {
			characterInfo := GetCharacterInfo(s, characterId)
			if characterInfo != nil {
				serverId = characterInfo.ServerId
			}
		}
		info.SupportSlotServerIds = append(info.SupportSlotServerIds, serverId)
	}
	i = 1
	for ; i <= GetSkillEchelonNum(proto.EchelonType(db.EchelonType)); i++ {
		var characterId int64 = 0
		characterId, ok := db.SkillCharacterList[i]
		if ok {
			characterInfo := GetCharacterInfo(s, characterId)
			if characterInfo == nil {
				characterId = 0
			}
		}
		info.SkillCardMulliganCharacterIds = append(info.SkillCardMulliganCharacterIds, characterId)
	}
	i = 1
	for ; i <= GetCombatStyleNum(proto.EchelonType(db.EchelonType)); i++ {
		info.CombatStyleIndex = append(info.CombatStyleIndex, 0)
	}

	return info
}

func GetEchelonPresetGroupDB(db *sro.EchelonInfo) *proto.EchelonPresetDB {
	if db == nil {
		return nil
	}
	info := &proto.EchelonPresetDB{
		GroupIndex:             db.EchelonType,
		Index:                  int32(db.EchelonNum),
		Label:                  "",
		LeaderUniqueId:         db.LeaderCharacter,
		TSSInteractionUniqueId: db.TssId,
		StrikerUniqueIds:       make([]int64, 0), // 主角色
		SpecialUniqueIds:       make([]int64, 0), // 支援角色
		CombatStyleIndex:       make([]int64, 0),
		MulliganUniqueIds:      make([]int64, 0),
		ExtensionType:          proto.EchelonExtensionType_Base,
		StrikerSlotCount:       0,
		SpecialSlotCount:       0,
	}
	var i int32 = 1
	for ; i <= 6; i++ {
		id, _ := db.MainCharacterList[i]
		info.StrikerUniqueIds = append(info.StrikerUniqueIds, id)
	}
	i = 1
	for ; i <= 2; i++ {
		id, _ := db.SupportCharacterList[i]
		info.SpecialUniqueIds = append(info.SpecialUniqueIds, id)
	}
	i = 1
	for ; i <= 6; i++ {
		info.CombatStyleIndex = append(info.CombatStyleIndex, 0)
	}

	return info
}

const (
	AssistTermRewardPeriodFromSec = 20
	AssistRewardLimit             = 1000000000
	AssistRentRewardDailyMaxCount = 20
	AssistRentalFeeAmount         = 40000
)

func GetAssistList(s *enter.Session) map[int32]*sro.AssistList {
	bin := GetEchelonBin(s)
	if bin == nil {
		return nil
	}
	if bin.AssistList == nil {
		bin.AssistList = make(map[int32]*sro.AssistList)
	}
	return bin.AssistList
}

func GetAssistListByEchelonType(s *enter.Session, echelonType proto.EchelonType) *sro.AssistList {
	bin := GetEchelonBin(s)
	if bin == nil {
		return nil
	}
	if bin.AssistList == nil {
		bin.AssistList = make(map[int32]*sro.AssistList)
	}
	return bin.AssistList[int32(echelonType)]
}

func GetAssistInfo(s *enter.Session, echelonType proto.EchelonType, slot int32) *sro.AssistInfo {
	bin := GetAssistListByEchelonType(s, echelonType)
	if bin == nil {
		return nil
	}
	return bin.AssistInfoList[slot]
}

func GetAssistInfoByClanAssistUseInfo(s *enter.Session, assist *proto.ClanAssistUseInfo) *sro.AssistInfo {
	if assist == nil {
		return nil
	}
	bin := GetAssistInfo(s, assist.EchelonType, assist.SlotNumber)
	if bin == nil {
		return nil
	}
	bin.IsMulligan = assist.IsMulligan
	bin.CombatStyleIndex = assist.CombatStyleIndex
	bin.TotalRentCount++
	return bin
}

func GetClanAssistSlotDBs(s *enter.Session) []*proto.ClanAssistSlotDB {
	list := make([]*proto.ClanAssistSlotDB, 0)
	for _, assist := range GetAssistList(s) {
		if assist.AssistInfoList == nil {
			assist.AssistInfoList = make(map[int32]*sro.AssistInfo)
		}
		for slot, info := range assist.AssistInfoList {
			clanAssistSlotDB := GetClanAssistSlotDB(s, info)
			if clanAssistSlotDB == nil {
				delete(assist.AssistInfoList, slot)
				continue
			}
			list = append(list, clanAssistSlotDB)
		}
	}

	return list
}

func GetClanAssistSlotDB(s *enter.Session, info *sro.AssistInfo) *proto.ClanAssistSlotDB {
	characterInfo := GetCharacterInfo(s, info.CharacterId)
	if characterInfo == nil {
		return nil
	}
	return &proto.ClanAssistSlotDB{
		EchelonType:      proto.EchelonType(info.EchelonType),
		SlotNumber:       info.SlotNumber,
		CharacterDBId:    characterInfo.ServerId,
		DeployDate:       mx.Unix(info.DeployDate, 0),
		TotalRentCount:   info.TotalRentCount,
		CombatStyleIndex: 0,
	}
}

func GetAssistCharacterDBs(s *enter.Session, assistRelation proto.AssistRelation) []*proto.AssistCharacterDB {
	list := make([]*proto.AssistCharacterDB, 0)
	for _, assist := range GetAssistList(s) {
		for slot, info := range assist.AssistInfoList {
			if assist.AssistInfoList == nil {
				assist.AssistInfoList = make(map[int32]*sro.AssistInfo)
			}
			assistCharacterDB := GetAssistCharacterDB(s, info, assistRelation)
			if assistCharacterDB == nil {
				delete(assist.AssistInfoList, slot)
				continue
			}
			list = append(list)
		}
	}

	return list
}

func GetAssistCharacterDB(s *enter.Session, info *sro.AssistInfo, assistRelation proto.AssistRelation) *proto.AssistCharacterDB {
	if info == nil {
		return nil
	}
	characterInfo := GetCharacterInfo(s, info.CharacterId)
	if characterInfo == nil {
		return nil
	}
	assistCharacterDB := &proto.AssistCharacterDB{
		EchelonType:             proto.EchelonType(info.EchelonType),
		AccountId:               s.AccountServerId,
		AssistRelation:          assistRelation,
		AssistCharacterServerId: characterInfo.ServerId,
		EquipmentDBs:            make([]*proto.EquipmentDB, 0),
		GearDB:                  GetGearDB(s, characterInfo.GearServerId),
		NickName:                GetNickname(s),
		SlotNumber:              int32(info.SlotNumber),
		WeaponDB:                GetWeaponDB(s, characterInfo.CharacterId),
		CharacterDB:             GetCharacterDB(s, characterInfo.CharacterId),
		IsMulligan:              info.IsMulligan,
		CombatStyleIndex:        info.CombatStyleIndex,

		CostumeId:        0,
		CostumeDB:        nil,
		IsTSAInteraction: false,
	}
	for _, serverId := range characterInfo.EquipmentList {
		if equipmentDB := GetEquipmentDB(s, serverId); equipmentDB != nil {
			assistCharacterDB.EquipmentDBs = append(assistCharacterDB.EquipmentDBs,
				equipmentDB)
		}
	}
	return assistCharacterDB
}
