// Code auto-generated by gucooing & Hiro420 DO NOT EDIT.

package proto

import (
	"github.com/bytedance/sonic"
	"github.com/gucooing/BaPs/protocol/mx"
)

type CafeDB struct {
	CafeDBId                int64
	CafeId                  int64
	AccountId               int64
	CafeRank                int32
	LastUpdate              mx.MxTime
	LastSummonDate          mx.MxTime
	IsNew                   bool
	CafeVisitCharacterDBs   map[int64]*CafeCharacterDB
	FurnitureDBs            []*FurnitureDB
	ProductionAppliedTime   mx.MxTime
	ProductionDB            *CafeProductionDB
	CurrencyDict_Obsolete   map[CurrencyTypes]int64
	UpdateTimeDict_Obsolete map[CurrencyTypes]mx.MxTime
}

func (x *CafeDB) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}
