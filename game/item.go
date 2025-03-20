package game

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func NewItemList(s *enter.Session) map[int64]*sro.ItemInfo {
	bin := GetItemBin(s)
	if bin == nil {
		return nil
	}
	list := make(map[int64]*sro.ItemInfo)
	sid := GetServerId(s)
	info := &sro.ItemInfo{
		ServerId:   sid,
		UniqueId:   2,
		StackCount: 5,
	}
	list[2] = info
	s.AddPlayerHash(sid, info)

	return list
}

func GetItemBin(s *enter.Session) *sro.ItemBin {
	bin := GetPlayerBin(s)
	if bin == nil {
		return nil
	}
	if bin.ItemBin == nil {
		bin.ItemBin = &sro.ItemBin{}
	}
	return bin.ItemBin
}

func GetItemList(s *enter.Session) map[int64]*sro.ItemInfo {
	bin := GetItemBin(s)
	if bin == nil {
		return nil
	}
	if bin.ItemInfoList == nil {
		bin.ItemInfoList = NewItemList(s)
	}
	return bin.ItemInfoList
}

func GetItemInfo(s *enter.Session, itemId int64) *sro.ItemInfo {
	bin := GetItemList(s)
	if bin == nil {
		return nil
	}
	return bin[itemId]
}

func AddItem(s *enter.Session, id int64, num int32) int64 {
	if !gdconf.IsItem(id) {
		return 0
	}
	bin := GetItemBin(s)
	if bin == nil {
		return 0
	}
	if bin.ItemInfoList == nil {
		bin.ItemInfoList = NewItemList(s)
	}
	if info, ok := bin.ItemInfoList[id]; ok {
		info.StackCount += num
		return info.ServerId
	}
	sid := GetServerId(s)
	info := &sro.ItemInfo{
		ServerId:   sid,
		UniqueId:   id,
		StackCount: num,
	}
	bin.ItemInfoList[id] = info
	s.AddPlayerHash(sid, info)
	return info.ServerId
}

func RemoveItem(s *enter.Session, id int64, num int32) bool {
	bin := GetItemBin(s)
	if bin == nil {
		return false
	}
	if bin.ItemInfoList == nil {
		bin.ItemInfoList = NewItemList(s)
	}
	if info, ok := bin.ItemInfoList[id]; ok {
		if info.StackCount >= num {
			info.StackCount -= num
			return true
		}
	}
	return false
}

var DefaultCurrencyNum = map[int32]int64{
	proto.CurrencyTypes_Gem:                      600,   // 肝的砖
	proto.CurrencyTypes_GemPaid:                  0,     // 氪的砖
	proto.CurrencyTypes_GemBonus:                 600,   // 总砖
	proto.CurrencyTypes_Gold:                     10000, // 金币
	proto.CurrencyTypes_ActionPoint:              24,    // 体力
	proto.CurrencyTypes_AcademyTicket:            3,     // 课程表
	proto.CurrencyTypes_ArenaTicket:              5,     // 竞技场
	proto.CurrencyTypes_RaidTicket:               6,     // 总力战
	proto.CurrencyTypes_WeekDungeonChaserATicket: 0,
	proto.CurrencyTypes_WeekDungeonChaserBTicket: 0,
	proto.CurrencyTypes_WeekDungeonChaserCTicket: 0,
	proto.CurrencyTypes_SchoolDungeonATicket:     0,
	proto.CurrencyTypes_SchoolDungeonBTicket:     0,
	proto.CurrencyTypes_SchoolDungeonCTicket:     0,
	proto.CurrencyTypes_TimeAttackDungeonTicket:  3, // 综合战术考试
	proto.CurrencyTypes_MasterCoin:               0,
	proto.CurrencyTypes_WorldRaidTicketA:         40,
	proto.CurrencyTypes_WorldRaidTicketB:         40,
	proto.CurrencyTypes_WorldRaidTicketC:         40,
	proto.CurrencyTypes_ChaserTotalTicket:        6, // 悬赏通缉
	proto.CurrencyTypes_SchoolDungeonTotalTicket: 6, // 学院交流会
	proto.CurrencyTypes_EliminateTicketA:         7, // 大决战
	proto.CurrencyTypes_EliminateTicketB:         7, // 大决战
	proto.CurrencyTypes_EliminateTicketC:         7, // 大决战
	proto.CurrencyTypes_EliminateTicketD:         7, // 大决战
}

func NewCurrencyInfo() map[int32]*sro.CurrencyInfo {
	list := make(map[int32]*sro.CurrencyInfo)
	for k, v := range DefaultCurrencyNum {
		list[k] = &sro.CurrencyInfo{
			CurrencyId:  k,
			CurrencyNum: v,
			UpdateTime:  time.Now().Unix(),
		}
	}
	return list
}

func UpCurrency(s *enter.Session, parcelId int64, num int64) *sro.CurrencyInfo {
	bin := GetCurrencyList(s)
	if bin == nil {
		return nil
	}
	info := GetCurrencyInfo(s, int32(parcelId))
	if info == nil {
		return nil
	}
	if num < 0 && info.CurrencyNum < -(num) {
		return nil
	}
	info.CurrencyNum += num
	if num > 0 {
		info.UpdateTime = time.Now().Unix()
	}

	// 砖石特殊处理
	if parcelId == proto.CurrencyTypes_Gem || parcelId == proto.CurrencyTypes_GemPaid {
		gem := GetCurrencyInfo(s, int32(parcelId))
		gemBonus := GetCurrencyInfo(s, proto.CurrencyTypes_GemBonus)
		gemBonus.UpdateTime = gem.UpdateTime
		gemBonus.CurrencyNum += num
	}
	// 体力特殊处理
	if parcelId == proto.CurrencyTypes_ActionPoint {
		if info.CurrencyNum > 999 {
			AddMailBySystem(s, "MaxActionPoint", []*sro.ParcelInfo{
				{
					Type: proto.ParcelType_Currency,
					Id:   proto.CurrencyTypes_ActionPoint,
					Num:  info.CurrencyNum - 999,
				},
			})
			info.CurrencyNum = 999
		}
		if num < 0 {
			AddAccountExp(s, -num) // 如果是体力扣除,就触发账号经验处理
		}
		info.UpdateTime = time.Now().Add(1 * time.Hour).Unix()
	}

	return info
}

// SetCurrency 直接设置,如需要产出请勿使用此方法
func SetCurrency(s *enter.Session, parcelId int64, num int64) {
	bin := GetCurrencyList(s)
	if bin == nil {
		return
	}
	info := GetCurrencyInfo(s, int32(parcelId))
	if info == nil {
		return
	}
	info.UpdateTime = time.Now().Unix()
	info.CurrencyNum = num
}

func GetCurrencyList(s *enter.Session) map[int32]*sro.CurrencyInfo {
	bin := GetItemBin(s)
	if bin == nil {
		return nil
	}
	if bin.CurrencyInfoList == nil {
		bin.CurrencyInfoList = NewCurrencyInfo()
	}
	for id, db := range bin.CurrencyInfoList {
		// 特殊物品刷新查询
		if time.Unix(db.UpdateTime, 0).Before(alg.GetLastDay4()) {
			switch id {
			case proto.CurrencyTypes_ChaserTotalTicket,
				proto.CurrencyTypes_SchoolDungeonTotalTicket,
				proto.CurrencyTypes_RaidTicket: // 总力战
				db.CurrencyNum = alg.MaxInt64(db.CurrencyNum, 6)
				db.UpdateTime = time.Now().Unix()
			case proto.CurrencyTypes_ArenaTicket:
				db.CurrencyNum = alg.MaxInt64(db.CurrencyNum, 5)
				db.UpdateTime = time.Now().Unix()
			case proto.CurrencyTypes_TimeAttackDungeonTicket:
				db.CurrencyNum = alg.MaxInt64(db.CurrencyNum, 3)
				db.UpdateTime = time.Now().Unix()
			case proto.CurrencyTypes_EliminateTicketA,
				proto.CurrencyTypes_EliminateTicketB,
				proto.CurrencyTypes_EliminateTicketC,
				proto.CurrencyTypes_EliminateTicketD:
				if GetEliminateRaidSeasonType() == proto.RaidSeasonType_Open && db.CurrencyNum != 0 {
					db.CurrencyNum = alg.MinInt64(db.CurrencyNum+1, 7)
					db.UpdateTime = time.Now().Unix()
				} else {
					db.CurrencyNum = 0
					db.UpdateTime = time.Now().Unix()
				}
			case proto.CurrencyTypes_AcademyTicket:
				db.CurrencyNum = alg.MaxInt64(db.CurrencyNum, GetMaxAcademyTicket(s))
				db.UpdateTime = time.Now().Unix()
			}
		}
		if id == proto.CurrencyTypes_ActionPoint {
			RecoverActionPoint(s, db)
		}
	}

	return bin.CurrencyInfoList
}

func GetCurrencyInfo(s *enter.Session, currencyId int32) *sro.CurrencyInfo {
	bin := GetCurrencyList(s)
	if bin == nil {
		return nil
	}
	if bin[currencyId] == nil {
		bin[currencyId] = &sro.CurrencyInfo{
			CurrencyId:  currencyId,
			CurrencyNum: DefaultCurrencyNum[currencyId],
			UpdateTime:  time.Now().Unix(),
		}
		if currencyId == proto.CurrencyTypes_ActionPoint {
			bin[currencyId].UpdateTime = time.Now().Add(1 * time.Hour).Unix()
		}
	}
	return bin[currencyId]
}

func GetAccountCurrencyDB(s *enter.Session) *proto.AccountCurrencyDB {
	accountCurrencyDB := &proto.AccountCurrencyDB{
		AccountLevel:           int64(GetAccountLevel(s)),
		AcademyLocationRankSum: GetAcademyLocationRankSum(s),
		CurrencyDict:           make(map[proto.CurrencyTypes]int64),
		UpdateTimeDict:         make(map[proto.CurrencyTypes]mx.MxTime),
	}
	for id, db := range GetCurrencyList(s) {
		accountCurrencyDB.CurrencyDict[proto.CurrencyTypes(proto.CurrencyTypes_name[id])] = db.CurrencyNum
		accountCurrencyDB.UpdateTimeDict[proto.CurrencyTypes(proto.CurrencyTypes_name[id])] = mx.Unix(db.UpdateTime, 0)
	}

	return accountCurrencyDB
}

func RecoverActionPoint(s *enter.Session, db *sro.CurrencyInfo) {
	if db == nil {
		return
	}
	maxAp := gdconf.GetAPAutoChargeMax(GetAccountLevel(s))
	if db.CurrencyNum >= maxAp {
		db.UpdateTime = time.Now().Add(1 * time.Hour).Unix()
		return
	}
	num := int64(time.Now().Sub(time.Unix(db.UpdateTime, 0).Add(-1*time.Hour)).Minutes() / 6)
	db.CurrencyNum = alg.MinInt64(db.CurrencyNum+num, maxAp)
	if num > 0 { // 这样处理精度不高 但是方便
		db.UpdateTime = time.Now().Add(1 * time.Hour).Unix()
	}
}

func GetItemDB(s *enter.Session, id int64) *proto.ItemDB {
	bin := GetItemInfo(s, id)
	if bin == nil {
		return nil
	}
	return &proto.ItemDB{
		Type:       proto.ParcelType_Item,
		ServerId:   bin.ServerId,
		UniqueId:   bin.UniqueId,
		StackCount: bin.StackCount,
	}
}

func GetWeaponInfoList(s *enter.Session) map[int64]*sro.WeaponInfo {
	bin := GetItemBin(s)
	if bin == nil {
		return nil
	}
	if bin.WeaponInfoList == nil {
		bin.WeaponInfoList = make(map[int64]*sro.WeaponInfo)
	}
	return bin.WeaponInfoList
}

func GetWeaponInfo(s *enter.Session, characterId int64) *sro.WeaponInfo {
	bin := GetWeaponInfoList(s)
	if bin == nil {
		return nil
	}
	return bin[characterId]
}

func AddWeapon(s *enter.Session, characterId int64) {
	bin := GetItemBin(s)
	if bin == nil {
		return
	}
	conf := gdconf.GetCharacterWeaponExcelTable(characterId)
	if bin.WeaponInfoList == nil {
		bin.WeaponInfoList = make(map[int64]*sro.WeaponInfo)
	}
	characterInfo := GetCharacterInfo(s, characterId)
	if conf == nil || characterInfo == nil ||
		characterInfo.StarGrade < 5 {
		return
	}
	info := &sro.WeaponInfo{
		UniqueId:          characterId,
		CharacterServerId: characterInfo.ServerId,
		StarGrade:         1,
		Level:             1,
		Exp:               0,
		IsLocked:          false,
	}
	info.Level = 50 // TODO 没有表
	info.StarGrade = 3
	bin.WeaponInfoList[characterId] = info
}

func GetWeaponDBs(s *enter.Session) []*proto.WeaponDB {
	list := make([]*proto.WeaponDB, 0)
	for _, bin := range GetWeaponInfoList(s) {
		list = append(list, &proto.WeaponDB{
			Type:                   proto.ParcelType_CharacterWeapon,
			UniqueId:               bin.UniqueId,
			Level:                  bin.Level,
			Exp:                    bin.Exp,
			StarGrade:              bin.StarGrade,
			BoundCharacterServerId: bin.CharacterServerId,
			IsLocked:               bin.IsLocked,
		})
	}

	return list
}

func GetWeaponDB(s *enter.Session, characterId int64) *proto.WeaponDB {
	bin := GetWeaponInfo(s, characterId)
	if bin == nil {
		return nil
	}
	return &proto.WeaponDB{
		Type:                   proto.ParcelType_CharacterWeapon,
		UniqueId:               bin.UniqueId,
		Level:                  bin.Level,
		Exp:                    bin.Exp,
		StarGrade:              bin.StarGrade,
		BoundCharacterServerId: bin.CharacterServerId,
		IsLocked:               bin.IsLocked,
	}
}

func GetEquipmentInfoList(s *enter.Session) map[int64]*sro.EquipmentInfo {
	bin := GetItemBin(s)
	if bin == nil {
		return nil
	}
	if bin.EquipmentInfoList == nil {
		bin.EquipmentInfoList = make(map[int64]*sro.EquipmentInfo)
	}
	return bin.EquipmentInfoList
}

func GetEquipmentInfo(s *enter.Session, serverId int64) *sro.EquipmentInfo {
	bin := GetEquipmentInfoList(s)
	if bin == nil {
		return nil
	}
	return bin[serverId]
}

// AddEquipment 传入装备id
// 如果是非佩戴装备设置id为k，并返回装备id
// 如果是佩戴装备，则设置唯一id为k，并返回唯一id/*
func AddEquipment(s *enter.Session, equipmentId int64, num int64) []int64 {
	bin := GetItemBin(s)
	if bin == nil {
		return nil
	}
	conf := gdconf.GetEquipmentExcelTable(equipmentId)
	if conf == nil {
		return nil
	}
	if bin.EquipmentInfoList == nil {
		bin.EquipmentInfoList = make(map[int64]*sro.EquipmentInfo)
	}
	if conf.MaxLevel < 10 { // 装备材料
		if info := s.GetEquipmentByKeyId(equipmentId); info != nil {
			info.StackCount += num
			return []int64{info.ServerId}
		} else {
			sid := GetServerId(s)
			equipmentInfo := &sro.EquipmentInfo{
				UniqueId:   equipmentId,
				ServerId:   sid,
				StackCount: num,
			}
			bin.EquipmentInfoList[sid] = equipmentInfo
			s.AddPlayerHash(equipmentId, equipmentInfo)
			return []int64{sid}
		}
	}
	sIdLi := make([]int64, 0)
	for i := 0; int64(i) < num; i++ {
		sid := GetServerId(s)
		bin.EquipmentInfoList[sid] = &sro.EquipmentInfo{
			UniqueId:          equipmentId,
			CharacterServerId: 0,
			Level:             1,
			Exp:               0,
			ServerId:          sid,
			Tier:              conf.TierInit,
			StackCount:        1,
			IsLocked:          false,
		}
		sIdLi = append(sIdLi, sid)
	}
	return sIdLi
}

func DelEquipment(s *enter.Session, serverId int64, num int64) (bool, int64) {
	bin := GetItemBin(s)
	if bin == nil {
		return false, 0
	}
	info, ok := bin.EquipmentInfoList[serverId]
	if !ok {
		return false, 0
	}
	conf := gdconf.GetEquipmentExcelTable(info.UniqueId)
	statConf := gdconf.GetEquipmentStatExcelTable(info.UniqueId)
	if conf == nil || statConf == nil {
		return false, 0
	}
	// 扣钱
	UpCurrency(s, int64(proto.CurrencyTypes_value[statConf.LevelUpFeedCostCurrency]),
		-(statConf.LevelUpFeedCostAmount * num))
	if conf.MaxLevel < 10 { // 装备材料
		info.StackCount -= num
	} else {
		delete(bin.EquipmentInfoList, info.ServerId)
	}
	return true, statConf.LevelUpFeedExp * num
}

func GetEquipmentDBs(s *enter.Session) []*proto.EquipmentDB {
	list := make([]*proto.EquipmentDB, 0)
	for index, bin := range GetEquipmentInfoList(s) {
		if conf := gdconf.GetEquipmentExcelTable(bin.UniqueId); conf == nil {
			delete(GetEquipmentInfoList(s), index)
			continue
		}
		list = append(list, &proto.EquipmentDB{
			Type:                   proto.ParcelType_Equipment,
			UniqueId:               bin.UniqueId,
			Level:                  bin.Level,
			Exp:                    bin.Exp,
			StackCount:             bin.StackCount,
			BoundCharacterServerId: bin.CharacterServerId,
			Tier:                   bin.Tier,
			ServerId:               bin.ServerId,
			IsLocked:               bin.IsLocked,
		})
	}

	return list
}

func GetEquipmentDB(s *enter.Session, serverId int64) *proto.EquipmentDB {
	bin := GetEquipmentInfo(s, serverId)
	if bin == nil {
		return nil
	}
	return &proto.EquipmentDB{
		Type:                   proto.ParcelType_Equipment,
		UniqueId:               bin.UniqueId,
		Level:                  bin.Level,
		Exp:                    bin.Exp,
		StackCount:             bin.StackCount,
		BoundCharacterServerId: bin.CharacterServerId,
		Tier:                   bin.Tier,
		ServerId:               bin.ServerId,
		IsLocked:               bin.IsLocked,
	}
}

func NewIdCardBackgroundList(s *enter.Session) map[int64]*sro.IdCardBackgroundInfo {
	list := make(map[int64]*sro.IdCardBackgroundInfo)
	for _, conf := range gdconf.GetIdCardBackgroundExcelList() {
		if conf.IsDefault {
			list[conf.Id] = &sro.IdCardBackgroundInfo{
				UniqueId: conf.Id,
				ServerId: GetServerId(s),
			}
		}
	}

	return list
}

func GetCardBackgroundIdInfoList(s *enter.Session) map[int64]*sro.IdCardBackgroundInfo {
	bin := GetItemBin(s)
	if bin == nil {
		return nil
	}
	if bin.IdCardBackgroundList == nil {
		bin.IdCardBackgroundList = NewIdCardBackgroundList(s)
	}
	return bin.IdCardBackgroundList
}

func GetCardBackgroundIdInfo(s *enter.Session, backgroundId int64) *sro.IdCardBackgroundInfo {
	bin := GetItemBin(s)
	if bin == nil {
		return nil
	}
	if bin.IdCardBackgroundList == nil {
		bin.IdCardBackgroundList = NewIdCardBackgroundList(s)
	}
	return bin.IdCardBackgroundList[backgroundId]
}

func AddCardBackgroundId(s *enter.Session, backgroundId int64) int64 {
	conf := gdconf.GetIdCardBackgroundExcel(backgroundId)
	bin := GetItemBin(s)
	if bin == nil || conf == nil {
		return 0
	}
	if bin.IdCardBackgroundList == nil {
		bin.IdCardBackgroundList = NewIdCardBackgroundList(s)
	}
	if v, ok := bin.IdCardBackgroundList[backgroundId]; ok {
		return v.ServerId
	}
	bin.IdCardBackgroundList[backgroundId] = &sro.IdCardBackgroundInfo{
		UniqueId: backgroundId,
		ServerId: GetServerId(s),
	}
	return bin.IdCardBackgroundList[backgroundId].ServerId
}

func GetIdCardBackgroundDBs(s *enter.Session) []*proto.IdCardBackgroundDB {
	list := make([]*proto.IdCardBackgroundDB, 0)
	for _, v := range GetCardBackgroundIdInfoList(s) {
		list = append(list, &proto.IdCardBackgroundDB{
			Type:     proto.ParcelType_IdCardBackground,
			ServerId: v.ServerId,
			UniqueId: v.UniqueId,
		})
	}

	return list
}

func GetIdCardBackgroundDB(s *enter.Session, backgroundId int64) *proto.IdCardBackgroundDB {
	bin := GetCardBackgroundIdInfo(s, backgroundId)
	if bin == nil {
		return nil
	}
	return &proto.IdCardBackgroundDB{
		Type:     proto.ParcelType_IdCardBackground,
		ServerId: bin.ServerId,
		UniqueId: bin.UniqueId,
	}
}

type ParcelResult struct {
	ParcelType proto.ParcelType
	ParcelId   int64
	Amount     int64
}

func GetParcelResultList(typeList []string, idList, numList []int64, isDel bool) []*ParcelResult {
	//  不验有没有那么多了,随意了
	list := make([]*ParcelResult, 0)
	if len(typeList) == len(idList) &&
		len(idList) == len(numList) {
		for index, rewardType := range typeList {
			num := numList[index]
			if isDel {
				num = -numList[index]
			}
			list = append(list, &ParcelResult{
				ParcelType: proto.GetParcelTypeValue(rewardType),
				ParcelId:   idList[index],
				Amount:     num,
			})
		}
	}
	return list
}

// ParcelResultDB TODO 没有验证是否能被消耗
func ParcelResultDB(s *enter.Session, parcelResultList []*ParcelResult) *proto.ParcelResultDB {
	info := &proto.ParcelResultDB{
		MemoryLobbyDBs:      make([]*proto.MemoryLobbyDB, 0),
		ItemDBs:             make(map[int64]*proto.ItemDB),
		EmblemDBs:           make([]*proto.EmblemDB, 0),
		CharacterDBs:        make([]*proto.CharacterDB, 0),
		WeaponDBs:           make([]*proto.WeaponDB, 0),
		EquipmentDBs:        make(map[int64]*proto.EquipmentDB),
		FurnitureDBs:        make(map[int64]*proto.FurnitureDB),
		IdCardBackgroundDBs: make(map[int64]*proto.IdCardBackgroundDB),
		AcademyLocationDBs:  make([]*proto.AcademyLocationDB, 0),

		CostumeDBs:                      nil,
		TSSCharacterDBs:                 nil,
		RemovedEquipmentIds:             nil,
		RemovedItemIds:                  nil,
		RemovedFurnitureIds:             nil,
		StickerDBs:                      nil,
		CharacterNewUniqueIds:           nil,
		SecretStoneCharacterIdAndCounts: nil,
		ParcelResultStepInfoList:        nil,
		BaseAccountExp:                  0,
		AdditionalAccountExp:            0,
		GachaResultCharacters:           nil,
	}
	defer func() {
		info.AccountDB = GetAccountDB(s)
		info.AccountCurrencyDB = GetAccountCurrencyDB(s)
	}()
	parcelInfoList := make([]*proto.ParcelInfo, 0)

	for _, parcelResult := range parcelResultList {
		isParcelResult := true
		switch parcelResult.ParcelType {
		case proto.ParcelType_Currency: // 货币
			UpCurrency(s, parcelResult.ParcelId, parcelResult.Amount)
			info.AccountCurrencyDB = GetAccountCurrencyDB(s)
		case proto.ParcelType_MemoryLobby: // 记忆大厅
			UpMemoryLobbyInfo(s, parcelResult.ParcelId)
			info.MemoryLobbyDBs = append(info.MemoryLobbyDBs,
				GetMemoryLobbyDB(s, parcelResult.ParcelId))
		case proto.ParcelType_Emblem: // 称号
			UpEmblemInfoList(s, []int64{parcelResult.ParcelId})
			info.EmblemDBs = append(info.EmblemDBs,
				GetEmblemDB(s, parcelResult.ParcelId))
		case proto.ParcelType_Item: // 背包物品
			serverId := AddItem(s, parcelResult.ParcelId, int32(parcelResult.Amount))
			if serverId == 0 {
				continue
			}
			info.ItemDBs[serverId] = GetItemDB(s, parcelResult.ParcelId)
		case proto.ParcelType_Character: // 角色
			if !AddCharacter(s, parcelResult.ParcelId) { // 重复添加处理
				for _, itemId := range RepeatAddCharacter(s, parcelResult.ParcelId) {
					if itemInfo := GetItemDB(s, itemId); itemInfo != nil {
						info.ItemDBs[itemInfo.ServerId] = itemInfo
					}
				}
			} else {
				info.CharacterDBs = append(info.CharacterDBs, GetCharacterDB(s, parcelResult.ParcelId))
			}
		case proto.ParcelType_FavorExp: // 角色好感度
			isParcelResult = false
			info.CharacterDBs = append(info.CharacterDBs, GetCharacterDB(s, parcelResult.ParcelId))
		case proto.ParcelType_CharacterWeapon: // 角色武器 仅同步
			isParcelResult = false
			info.WeaponDBs = append(info.WeaponDBs, GetWeaponDB(s, parcelResult.ParcelId))
		case proto.ParcelType_Equipment: // 装备 仅添加
			for _, serverId := range AddEquipment(s, parcelResult.ParcelId, parcelResult.Amount) {
				info.EquipmentDBs[serverId] = GetEquipmentDB(s, serverId)
			}
		case proto.ParcelType_Furniture: // 家具 仅添加
			for _, serverId := range AddFurnitureInfo(s, parcelResult.ParcelId, parcelResult.Amount) {
				info.FurnitureDBs[serverId] = GetFurnitureDB(s, serverId)
			}
		case proto.ParcelType_IdCardBackground: // 账号背景页 仅添加
			serverid := AddCardBackgroundId(s, parcelResult.ParcelId)
			info.IdCardBackgroundDBs[serverid] = GetIdCardBackgroundDB(s, parcelResult.ParcelId)
		case proto.ParcelType_LocationExp: // 课程表经验更改
			UpAcademyLocationExp(s, parcelResult.ParcelId, parcelResult.Amount)
			info.AcademyLocationDBs = append(info.AcademyLocationDBs, GetAcademyLocationDB(s, parcelResult.ParcelId))
		default:
			logger.Warn("没有处理的奖励类型 Unknown ParcelType:%s", parcelResult.ParcelType.String())
		}
		if parcelResult.Amount >= 0 && isParcelResult {
			parcelInfoList = append(parcelInfoList, &proto.ParcelInfo{
				Key: &proto.ParcelKeyPair{
					Type: parcelResult.ParcelType,
					Id:   parcelResult.ParcelId,
				},
				Amount: parcelResult.Amount,
				Multiplier: &proto.BasisPoint{
					RawValue: 10000,
				},
				Probability: &proto.BasisPoint{
					RawValue: 10000,
				},
			})
		}
	}
	info.DisplaySequence = parcelInfoList
	info.ParcelForMission = parcelInfoList

	return info
}

func NoGMFack(parcelType proto.ParcelType, id int64) bool {
	switch parcelType {
	case proto.ParcelType_Character:
		return gdconf.GetCharacterExcel(id) == nil
	case proto.ParcelType_Currency:
		_, ok := DefaultCurrencyNum[int32(id)]
		return !ok
	case proto.ParcelType_Equipment:
		return gdconf.GetEquipmentExcelTable(id) == nil
	case proto.ParcelType_Item:
		return gdconf.GetItemExcelTable(id) == nil
	case proto.ParcelType_Furniture:
		return gdconf.GetFurnitureExcelTable(id) == nil
	default:
		return true
	}
}
