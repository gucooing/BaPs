package game

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func NewItemList(s *enter.Session) map[int64]*sro.ItemInfo {
	list := make(map[int64]*sro.ItemInfo)
	list[2] = &sro.ItemInfo{
		ServerId:   GetServerId(s),
		UniqueId:   2,
		StackCount: 5,
	}
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
	info := &sro.ItemInfo{
		ServerId:   GetServerId(s),
		UniqueId:   id,
		StackCount: num,
	}
	bin.ItemInfoList[id] = info
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
	proto.CurrencyTypes_Gem:                      600,
	proto.CurrencyTypes_GemPaid:                  0,
	proto.CurrencyTypes_GemBonus:                 600,   // 砖石
	proto.CurrencyTypes_Gold:                     10000, // 金币
	proto.CurrencyTypes_ActionPoint:              24,    // 体力
	proto.CurrencyTypes_AcademyTicket:            3,
	proto.CurrencyTypes_ArenaTicket:              5,
	proto.CurrencyTypes_RaidTicket:               3,
	proto.CurrencyTypes_WeekDungeonChaserATicket: 0,
	proto.CurrencyTypes_WeekDungeonChaserBTicket: 0,
	proto.CurrencyTypes_WeekDungeonChaserCTicket: 0,
	proto.CurrencyTypes_SchoolDungeonATicket:     0,
	proto.CurrencyTypes_SchoolDungeonBTicket:     0,
	proto.CurrencyTypes_SchoolDungeonCTicket:     0,
	proto.CurrencyTypes_TimeAttackDungeonTicket:  3,
	proto.CurrencyTypes_MasterCoin:               0,
	proto.CurrencyTypes_WorldRaidTicketA:         40,
	proto.CurrencyTypes_WorldRaidTicketB:         40,
	proto.CurrencyTypes_WorldRaidTicketC:         40,
	proto.CurrencyTypes_ChaserTotalTicket:        6, // 悬赏通缉
	proto.CurrencyTypes_SchoolDungeonTotalTicket: 6, // 学院交流会
	proto.CurrencyTypes_EliminateTicketA:         3,
	proto.CurrencyTypes_EliminateTicketB:         3,
	proto.CurrencyTypes_EliminateTicketC:         3,
	proto.CurrencyTypes_EliminateTicketD:         3,
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
	info.UpdateTime = time.Now().Unix()

	gemBonus := GetCurrencyInfo(s, proto.CurrencyTypes_GemBonus)
	gem := GetCurrencyInfo(s, proto.CurrencyTypes_Gem)
	if gem != nil || gemBonus != nil {
		gem.CurrencyNum = gemBonus.CurrencyNum
	}
	if parcelId == proto.CurrencyTypes_ActionPoint &&
		num < 0 {
		AddAccountExp(s, -num) // 如果是体力扣除,就触发账号经验处理
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
	return bin.CurrencyInfoList
}

func GetCurrencyInfo(s *enter.Session, currencyId int32) *sro.CurrencyInfo {
	bin := GetItemBin(s)
	if bin == nil {
		return nil
	}
	if bin.CurrencyInfoList == nil {
		bin.CurrencyInfoList = NewCurrencyInfo()
	}
	if bin.CurrencyInfoList[currencyId] == nil {
		bin.CurrencyInfoList[currencyId] = &sro.CurrencyInfo{
			CurrencyId:  currencyId,
			CurrencyNum: DefaultCurrencyNum[currencyId],
			UpdateTime:  time.Now().Unix(),
		}
	}
	return bin.CurrencyInfoList[currencyId]
}

func GetAccountCurrencyDB(s *enter.Session) *proto.AccountCurrencyDB {
	accountCurrencyDB := &proto.AccountCurrencyDB{
		AccountLevel:           int64(GetAccountLevel(s)),
		AcademyLocationRankSum: 1,
		CurrencyDict:           make(map[proto.CurrencyTypes]int64),
		UpdateTimeDict:         make(map[proto.CurrencyTypes]mx.MxTime),
	}
	for id, db := range GetCurrencyList(s) {
		accountCurrencyDB.CurrencyDict[proto.CurrencyTypes(proto.CurrencyTypes_name[id])] = db.CurrencyNum
		accountCurrencyDB.UpdateTimeDict[proto.CurrencyTypes(proto.CurrencyTypes_name[id])] = mx.Unix(db.UpdateTime, 0)
	}

	return accountCurrencyDB
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
	bin.WeaponInfoList[characterId] = &sro.WeaponInfo{
		UniqueId:          characterId,
		CharacterServerId: characterInfo.ServerId,
		StarGrade:         1,
		Level:             1,
		Exp:               0,
		IsLocked:          false,
	}
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
				ParcelType: proto.ParcelType(proto.ParcelType_value[rewardType]),
				ParcelId:   idList[index],
				Amount:     num,
			})
		}
	}
	return list
}

func ParcelResultDB(s *enter.Session, parcelResultList []*ParcelResult) *proto.ParcelResultDB {
	info := &proto.ParcelResultDB{
		AccountDB:      GetAccountDB(s),
		MemoryLobbyDBs: make([]*proto.MemoryLobbyDB, 0),
		ItemDBs:        make(map[int64]*proto.ItemDB),
		EmblemDBs:      make([]*proto.EmblemDB, 0),
		CharacterDBs:   make([]*proto.CharacterDB, 0),
		WeaponDBs:      make([]*proto.WeaponDB, 0),

		AcademyLocationDBs:              nil,
		CostumeDBs:                      nil,
		TSSCharacterDBs:                 nil,
		EquipmentDBs:                    nil,
		RemovedEquipmentIds:             nil,
		RemovedItemIds:                  nil,
		FurnitureDBs:                    nil,
		RemovedFurnitureIds:             nil,
		IdCardBackgroundDBs:             nil,
		StickerDBs:                      nil,
		CharacterNewUniqueIds:           nil,
		SecretStoneCharacterIdAndCounts: nil,
		DisplaySequence:                 make([]*proto.ParcelInfo, 0),
		ParcelForMission:                nil,
		ParcelResultStepInfoList:        nil,
		BaseAccountExp:                  0,
		AdditionalAccountExp:            0,
		GachaResultCharacters:           nil,
	}

	for _, parcelResult := range parcelResultList {
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
			info.ItemDBs[serverId] = GetItemDB(s, parcelResult.ParcelId)
		case proto.ParcelType_Character: // 角色
			AddCharacter(s, parcelResult.ParcelId)
			info.CharacterDBs = append(info.CharacterDBs, GetCharacterDB(s, parcelResult.ParcelId))
		case proto.ParcelType_CharacterWeapon: // 角色武器 仅同步
			info.WeaponDBs = append(info.WeaponDBs, GetWeaponDB(s, parcelResult.ParcelId))
		default:
			logger.Warn("没有处理的奖励类型 Unknown ParcelType:%s", parcelResult.ParcelType.String())
		}
		info.DisplaySequence = append(info.DisplaySequence, &proto.ParcelInfo{
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

	return info
}
