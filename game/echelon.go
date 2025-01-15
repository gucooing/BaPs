package game

import (
	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/mx/proto"
	"github.com/gucooing/BaPs/pkg/logger"
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
		info := NewEchelonInfo(list[conf.EchlonId], conf)
		if info == nil {
			continue
		}
		list[conf.EchlonId].EchelonInfoList[info.EchelonNum] = info
	}
	return list
}

func NewEchelonInfo(typeInfo *sro.EchelonTypeInfo, conf *sro.DefaultEchelonExcelTable) *sro.EchelonInfo {
	if typeInfo == nil || conf == nil {
		return nil
	}
	info := &sro.EchelonInfo{
		EchelonType:          conf.EchlonId,
		ExtensionType:        proto.EchelonExtensionType_Base,
		EchelonNum:           typeInfo.EchelonNum,
		LeaderCharacter:      conf.LeaderId,
		TssId:                conf.TssId,
		MainCharacterList:    make(map[int32]int64),
		SupportCharacterList: make(map[int32]int64),
		SkillCharacterList:   make(map[int32]int64),
	}
	var i int32 = 1
	for ; i <= GetMaxMainEchelonNum(proto.EchelonType(conf.EchlonId)); i++ {
		if len(conf.MainId) < int(i) {
			info.MainCharacterList[i] = 0
		} else {
			info.MainCharacterList[i] = conf.MainId[i-1]
		}
	}
	i = 1
	for ; i <= GetSupportEchelonNum(proto.EchelonType(conf.EchlonId)); i++ {
		if len(conf.SupportId) < int(i) {
			info.SupportCharacterList[i] = 0
		} else {
			info.SupportCharacterList[i] = conf.SupportId[i-1]
		}
	}
	// i  = 1
	// for ;i<= GetSkillEchelonNum(proto.EchelonType(conf.EchlonId));i++ {
	// 	info.MainCharacterList[i] = conf.MainId[i-1]
	// }
	typeInfo.EchelonNum++
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
		var serverId int64 = 0
		characterId, ok := db.SkillCharacterList[i]
		if ok {
			characterInfo := GetCharacterInfo(s, characterId)
			if characterInfo != nil {
				serverId = characterInfo.ServerId
			}
		}
		info.SkillCardMulliganCharacterIds = append(info.SkillCardMulliganCharacterIds, serverId)
	}
	i = 1
	for ; i <= GetCombatStyleNum(proto.EchelonType(db.EchelonType)); i++ {
		info.SkillCardMulliganCharacterIds = append(info.SkillCardMulliganCharacterIds, 0)
	}

	return info
}
