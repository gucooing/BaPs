package game

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/mx/proto"
	"github.com/gucooing/BaPs/pkg/logger"
)

func NewCurrency() *sro.CurrencyBin {
	bin := &sro.CurrencyBin{
		CurrencyInfoList: make(map[int32]*sro.CurrencyInfo),
	}

	return bin
}

var DefaultCurrencyNum = map[int32]int64{
	proto.CurrencyTypes_Gold:                      1000,
	proto.CurrencyTypes_GemPaid:                   0,
	proto.CurrencyTypes_GemBonus:                  600,
	proto.CurrencyTypes_Gem:                       600,
	proto.CurrencyTypes_ActionPoint:               24,
	proto.CurrencyTypes_AcademyTicket:             3,
	proto.CurrencyTypes_ArenaTicket:               5,
	proto.CurrencyTypes_RaidTicket:                3,
	proto.CurrencyTypes_WeekDungeonChaserATicket:  0,
	proto.CurrencyTypes_WeekDungeonFindGiftTicket: 0,
	proto.CurrencyTypes_WeekDungeonBloodTicket:    0,
	proto.CurrencyTypes_WeekDungeonChaserBTicket:  0,
	proto.CurrencyTypes_WeekDungeonChaserCTicket:  0,
	proto.CurrencyTypes_SchoolDungeonATicket:      0,
	proto.CurrencyTypes_SchoolDungeonBTicket:      0,
	proto.CurrencyTypes_SchoolDungeonCTicket:      0,
	proto.CurrencyTypes_TimeAttackDungeonTicket:   3,
	proto.CurrencyTypes_MasterCoin:                0,
	proto.CurrencyTypes_WorldRaidTicketA:          40,
	proto.CurrencyTypes_WorldRaidTicketB:          40,
	proto.CurrencyTypes_WorldRaidTicketC:          40,
	proto.CurrencyTypes_ChaserTotalTicket:         6,
	proto.CurrencyTypes_SchoolDungeonTotalTicket:  6,
	proto.CurrencyTypes_EliminateTicketA:          3,
	proto.CurrencyTypes_EliminateTicketB:          3,
	proto.CurrencyTypes_EliminateTicketC:          3,
	proto.CurrencyTypes_EliminateTicketD:          3,
	proto.CurrencyTypes_Max:                       0,
}

func NewCurrencyInfo(currencyId int32) *sro.CurrencyInfo {
	if currencyId == 0 {
		return nil
	}
	currencyNum, ok := DefaultCurrencyNum[currencyId]
	if !ok {
		logger.Error("未初始化的货币:%v", currencyId)
		return nil
	}
	return &sro.CurrencyInfo{
		CurrencyId:  currencyId,
		UpdateTime:  time.Now().Unix(),
		CurrencyNum: currencyNum,
	}
}

func GetCurrencyBin(s *enter.Session) *sro.CurrencyBin {
	bin := GetPlayerBin(s)
	if bin == nil {
		return nil
	}
	if bin.CurrencyBin == nil {
		bin.CurrencyBin = NewCurrency()
	}
	return bin.CurrencyBin
}

func GetAccountCurrencyDB(s *enter.Session) *proto.AccountCurrencyDB {
	accountCurrencyDB := &proto.AccountCurrencyDB{
		AccountLevel:           1,
		AcademyLocationRankSum: 1,
		CurrencyDict:           make(map[proto.CurrencyTypes]int64),
		UpdateTimeDict:         make(map[proto.CurrencyTypes]time.Time),
	}
	bin := GetCurrencyBin(s)
	if bin == nil {
		return accountCurrencyDB
	}
	if bin.CurrencyInfoList == nil {
		bin.CurrencyInfoList = make(map[int32]*sro.CurrencyInfo)
	}
	for v, id := range proto.CurrencyTypes_value {
		db, ok := bin.CurrencyInfoList[id]
		if !ok {
			db = NewCurrencyInfo(id)
		}
		if db == nil {
			continue
		}
		accountCurrencyDB.CurrencyDict[proto.CurrencyTypes(v)] = db.CurrencyNum
		accountCurrencyDB.UpdateTimeDict[proto.CurrencyTypes(v)] = time.Unix(db.UpdateTime, 0)
	}

	return accountCurrencyDB
}
