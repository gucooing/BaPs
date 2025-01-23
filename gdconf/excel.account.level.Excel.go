package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadAccountLevelExcel() {
	g.GetExcel().AccountLevelExcel = make([]*sro.AccountLevelExcel, 0)
	name := "AccountLevelExcel.json"
	file, err := os.ReadFile(g.excelDbPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().AccountLevelExcel); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetAccountLevelExcel()))
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
			return newLevel, exp
		}
		if exp >= conf.Exp {
			exp -= conf.Exp
		} else {
			return newLevel, exp
		}
	}
}

func GetAPAutoChargeMax(level int32) int64 {
	conf := GC.GetGPP().AccountLevel.AccountLevelExcelMap[level]
	if conf == nil {
		return 0
	}
	return conf.APAutoChargeMax
}
