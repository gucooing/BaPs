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
		CharacterHash:     make(map[int64]int64),
	}
	for _, conf := range gdconf.GetDefaultCharacterExcelTable() {
		sid := GetServerId(s)
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
			ServerId:               sid,
			IsFavorite:             false,
		}
		bin.CharacterHash[sid] = conf.CharacterId
		bin.CharacterInfoList[conf.CharacterId] = infp
	}

	return bin
}

func NewCharacterEquipment(characterId int64) map[int32]int64 {
	conf := gdconf.GetCharacterExcel(characterId)
	if conf == nil {
		return nil
	}
	list := make(map[int32]int64)
	for i := 0; i < len(conf.EquipmentSlot); i++ {
		list[int32(i)] = 0
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

func GetCharacterCount(s *enter.Session) int64 {
	return int64(len(GetCharacterInfoList(s)))
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

func GetCharacterHash(s *enter.Session) map[int64]int64 {
	bin := GetCharacterBin(s)
	if bin == nil {
		return nil
	}
	if bin.CharacterHash == nil {
		bin.CharacterHash = make(map[int64]int64)
	}
	return bin.CharacterHash
}

func ServerIdToCharacterId(s *enter.Session, serverId int64) int64 {
	if id, ok := GetCharacterHash(s)[serverId]; ok {
		return id
	}
	return 0
}

func GetCharacterInfoByServerId(s *enter.Session, serverId int64) *sro.CharacterInfo {
	hash := GetCharacterHash(s)
	if hash == nil {
		return nil
	}
	return GetCharacterInfo(s, hash[serverId])
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
		IsFavorite:             bin.IsFavorite,
		EquipmentServerIds:     make([]int64, 0),
		PotentialStats:         make(map[int32]int32),
	}
	info.FavorRank = 50 // 由于excel中没有好感度配置所以默认满级
	for i := 0; i < 3; i++ {
		e, ok := bin.EquipmentList[int32(i)]
		if ok {
			if GetEquipmentInfo(s, e) == nil {
				bin.EquipmentList[int32(i)] = 0
				e = 0
			}
			info.EquipmentServerIds = append(info.EquipmentServerIds, e)
		} else {
			info.EquipmentServerIds = append(info.EquipmentServerIds, 0)
		}
	}
	for _, p := range []int32{1, 2, 3} {
		info.PotentialStats[p] = 0
	}

	return info
}

func AddCharacter(s *enter.Session, characterId int64) bool {
	bin := GetCharacterBin(s)
	if bin == nil {
		return false
	}
	if bin.CharacterInfoList == nil {
		bin.CharacterInfoList = make(map[int64]*sro.CharacterInfo)
	}
	if bin.CharacterHash == nil {
		bin.CharacterHash = make(map[int64]int64)
	}
	if _, ok := bin.CharacterInfoList[characterId]; ok {
		return false
	}
	conf := gdconf.GetCharacterExcel(characterId)
	if conf == nil {
		logger.Error("[UID:%v]未知的角色添加调用,characterId:%v", s.AccountServerId, characterId)
		return true
	}
	sId := GetServerId(s)
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
		ServerId:               sId,
	}
	bin.CharacterHash[sId] = characterId
	bin.CharacterInfoList[characterId] = info
	return true
}

func RepeatAddCharacter(s *enter.Session, characterId int64) []int64 {
	conf := gdconf.GetCharacterExcel(characterId)
	if conf == nil {
		return nil
	}
	list := make([]int64, 0)
	// 添加秘石
	AddItem(s, conf.SecretStoneItemId, conf.SecretStoneItemAmount)
	list = append(list, conf.SecretStoneItemId)
	// 添加碎片
	AddItem(s, conf.CharacterPieceItemId, conf.CharacterPieceItemAmount)
	list = append(list, conf.CharacterPieceItemId)
	return list
}

func ServerIdsToCharacterIds(s *enter.Session, serverIdList []int64) []int64 {
	list := make([]int64, 0)
	for _, serverId := range serverIdList {
		if db := GetCharacterInfoByServerId(s, serverId); db != nil {
			list = append(list, db.CharacterId)
		}
	}

	return list
}

var CharacterStarGradeMap = map[int32]int32{
	2: 30,  // 1->2
	3: 80,  // 2->3
	4: 100, // 3->4
	5: 120, // 4->5
}

func GetCharacterUpStarGradeNum(starGrade int32) int32 {
	return CharacterStarGradeMap[starGrade]
}

var WeaponStarGradeMap = map[int32]int32{
	1: 120, // 1->2
	2: 180, // 2->3
}

func GetWeaponUpStarGradeNum(starGrade int32) int32 {
	return WeaponStarGradeMap[starGrade]
}

func SetCharacterEquipment(s *enter.Session, characterServerId int64, equipmentIdServerId int64, index int32) bool {
	characterInfo := GetCharacterInfoByServerId(s, characterServerId)
	equipmentInfo := GetEquipmentInfo(s, equipmentIdServerId)
	if characterInfo == nil || equipmentInfo == nil {
		return false
	}
	if len(characterInfo.EquipmentList) < 3 {
		characterInfo.EquipmentList = NewCharacterEquipment(characterInfo.CharacterId)
	}
	// 不存在卸载的情况
	characterInfo.EquipmentList[index] = equipmentIdServerId
	equipmentInfo.CharacterServerId = characterServerId
	return true
}

func GetGearInfoList(s *enter.Session) map[int64]*sro.GearInfo {
	bin := GetCharacterBin(s)
	if bin == nil {
		return nil
	}
	if bin.GearInfoList == nil {
		bin.GearInfoList = make(map[int64]*sro.GearInfo)
	}
	return bin.GearInfoList
}

func GetGearInfo(s *enter.Session, serverId int64) *sro.GearInfo {
	bin := GetGearInfoList(s)
	if bin == nil {
		return nil
	}
	return bin[serverId]
}

func GetGearDB(s *enter.Session, serverId int64) *proto.GearDB {
	bin := GetGearInfo(s, serverId)
	if bin == nil {
		return nil
	}
	info := &proto.GearDB{
		Type:                   proto.ParcelType_CharacterGear,
		ServerId:               bin.ServerId,
		UniqueId:               bin.UniqueId,
		Level:                  bin.Level,
		Exp:                    bin.Exp,
		Tier:                   bin.Tier,
		SlotIndex:              bin.SlotIndex,
		BoundCharacterServerId: bin.CharacterServerId,
	}

	return info
}
