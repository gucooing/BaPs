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

var UpdateTimeDict = map[int32]bool{
	proto.CurrencyTypes_ActionPoint:              true,
	proto.CurrencyTypes_AcademyTicket:            true,
	proto.CurrencyTypes_ArenaTicket:              true,
	proto.CurrencyTypes_RaidTicket:               true,
	proto.CurrencyTypes_WeekDungeonChaserATicket: true,
	proto.CurrencyTypes_WeekDungeonChaserBTicket: true,
	proto.CurrencyTypes_WeekDungeonChaserCTicket: true,
	proto.CurrencyTypes_SchoolDungeonATicket:     true,
	proto.CurrencyTypes_SchoolDungeonBTicket:     true,
	proto.CurrencyTypes_SchoolDungeonCTicket:     true,
	proto.CurrencyTypes_TimeAttackDungeonTicket:  true,
	proto.CurrencyTypes_MasterCoin:               true,
	proto.CurrencyTypes_WorldRaidTicketA:         true,
	proto.CurrencyTypes_WorldRaidTicketB:         true,
	proto.CurrencyTypes_WorldRaidTicketC:         true,
	proto.CurrencyTypes_ChaserTotalTicket:        true,
	proto.CurrencyTypes_SchoolDungeonTotalTicket: true,
	proto.CurrencyTypes_EliminateTicketA:         true,
	proto.CurrencyTypes_EliminateTicketB:         true,
	proto.CurrencyTypes_EliminateTicketC:         true,
	proto.CurrencyTypes_EliminateTicketD:         true,
}

func NewCurrencyInfo(currencyId int32) *sro.CurrencyInfo {
	if currencyId == proto.CurrencyTypes_Max ||
		currencyId == proto.CurrencyTypes_Invalid {
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
	for id, _ := range DefaultCurrencyNum {
		db, ok := bin.CurrencyInfoList[id]
		if !ok {
			db = NewCurrencyInfo(id)
		}
		if db == nil {
			continue
		}
		bin.CurrencyInfoList[id] = db
		accountCurrencyDB.CurrencyDict[proto.CurrencyTypes(proto.CurrencyTypes_name[id])] = db.CurrencyNum
		if ok = UpdateTimeDict[id]; !ok {
			continue
		}
		accountCurrencyDB.UpdateTimeDict[proto.CurrencyTypes(proto.CurrencyTypes_name[id])] = time.Unix(db.UpdateTime, 0)
	}

	return accountCurrencyDB
}
