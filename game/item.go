package game

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/mx/proto"
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

var DefaultCurrencyNum = map[int32]int64{
	proto.CurrencyTypes_Gem:                      600,
	proto.CurrencyTypes_GemPaid:                  0,
	proto.CurrencyTypes_GemBonus:                 600,
	proto.CurrencyTypes_Gold:                     10000,
	proto.CurrencyTypes_ActionPoint:              24,
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
	proto.CurrencyTypes_ChaserTotalTicket:        6,
	proto.CurrencyTypes_SchoolDungeonTotalTicket: 6,
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

func UpCurrencyGem(s *enter.Session) *sro.CurrencyInfo {
	bin := GetCurrencyList(s)
	if bin == nil {
		return nil
	}
	gemBonus := GetCurrencyInfo(s, proto.CurrencyTypes_GemBonus)
	gem := GetCurrencyInfo(s, proto.CurrencyTypes_Gem)
	gemPaid := GetCurrencyInfo(s, proto.CurrencyTypes_GemPaid)

	if gemBonus == nil || gem == nil || gemPaid == nil ||
		gem.CurrencyNum < gemPaid.CurrencyNum {
		return gemBonus
	}

	gemBonus.CurrencyNum = gem.CurrencyNum - gemPaid.CurrencyNum
	gemPaid.CurrencyNum = 0
	gem.CurrencyNum = gemBonus.CurrencyNum

	gemBonus.UpdateTime = time.Now().Unix()
	gemPaid.UpdateTime = time.Now().Unix()
	gem.UpdateTime = time.Now().Unix()

	return gemBonus
}

func SetGemPaid(s *enter.Session, num int64) bool {
	gemPaid := GetCurrencyInfo(s, proto.CurrencyTypes_GemPaid)
	gem := GetCurrencyInfo(s, proto.CurrencyTypes_Gem)
	if gemPaid == nil || gem == nil {
		return false
	}
	if gem.CurrencyNum < (gemPaid.CurrencyNum + num) {
		return false
	}
	gemPaid.CurrencyNum += num
	return true
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
		UpdateTimeDict:         make(map[proto.CurrencyTypes]time.Time),
	}
	for id, db := range GetCurrencyList(s) {
		accountCurrencyDB.CurrencyDict[proto.CurrencyTypes(proto.CurrencyTypes_name[id])] = db.CurrencyNum
		accountCurrencyDB.UpdateTimeDict[proto.CurrencyTypes(proto.CurrencyTypes_name[id])] = time.Unix(db.UpdateTime, 0)
	}

	return accountCurrencyDB
}
