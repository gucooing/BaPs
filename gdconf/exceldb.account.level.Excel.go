package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadAccountLevelExcel() {
	g.GetExcel().AccountLevelExcel = make([]*sro.AccountLevelExcel, 0)
	name := "AccountLevelExcel.json"
	loadExcelJson(g.excelDbPath+name, &g.GetExcel().AccountLevelExcel)
}

type AccountLevel struct {
	AccountLevelExcelMap map[int32]*sro.AccountLevelExcel
}

func (g *GameConfig) gppAccountLevelExcel() {
	g.GetGPP().AccountLevel = &AccountLevel{
		AccountLevelExcelMap: make(map[int32]*sro.AccountLevelExcel, 0),
	}

	for _, v := range g.GetExcel().GetAccountLevelExcel() {
		g.GetGPP().AccountLevel.AccountLevelExcelMap[v.Level] = v
	}

	logger.Info("处理账号等级配置完成数量:%v个", len(g.GetGPP().AccountLevel.AccountLevelExcelMap))
}

func UpAccountLevel(level int32, exp int64) (int32, int64) {
	i := int32(0)
	for ; ; i++ {
		newLevel := level + i
		conf := GC.GetGPP().AccountLevel.AccountLevelExcelMap[newLevel]
		if conf == nil {
			return newLevel - 1, exp
		}
		if exp >= conf.Exp {
			exp -= conf.Exp
			if conf.Exp == 0 {
				exp = 0 // 特殊处理，避免客户端出现超升级动画的神奇现象
			}
		} else {
			return newLevel, exp
		}
	}
}

func GetAccountLevelExcel(level int32) *sro.AccountLevelExcel {
	return GC.GetGPP().AccountLevel.AccountLevelExcelMap[level]
}

func GetAPAutoChargeMax(level int32) int64 {
	conf := GC.GetGPP().AccountLevel.AccountLevelExcelMap[level]
	if conf == nil {
		return 0
	}
	return conf.APAutoChargeMax
}
