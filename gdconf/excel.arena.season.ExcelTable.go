package gdconf

import (
	"time"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/mx"
)

func (g *GameConfig) loadArenaSeasonExcelTable() {
	g.GetExcel().ArenaSeasonExcelTable = make([]*sro.ArenaSeasonExcelTable, 0)
	name := "ArenaSeasonExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().ArenaSeasonExcelTable)
}

type ArenaSeasonExcel struct {
	ArenaSeasonExcelMap map[int64]*sro.ArenaSeasonExcelTable
	CurArenaSeason      *sro.ArenaSeasonExcelTable
}

func (g *GameConfig) gppArenaSeasonExcelTable() {
	g.GetGPP().ArenaSeasonExcel = &ArenaSeasonExcel{
		ArenaSeasonExcelMap: make(map[int64]*sro.ArenaSeasonExcelTable),
	}

	for _, v := range g.GetExcel().GetArenaSeasonExcelTable() {
		g.GetGPP().ArenaSeasonExcel.ArenaSeasonExcelMap[v.UniqueId] = v
	}

	logger.Info("处理竞技场赛季排期配置表完成,大决战总分奖励配置:%v个",
		len(g.GetGPP().ArenaSeasonExcel.ArenaSeasonExcelMap))
}

func genCurArenaSeason() *sro.ArenaSeasonExcelTable {
	for _, v := range GC.GetExcel().GetArenaSeasonExcelTable() { // 读取原始文件,保证顺序
		startTime, err := time.Parse("2006-01-02 15:04:05", v.SeasonStartDate)
		endTime, err := time.Parse("2006-01-02 15:04:05", v.SeasonEndDate)
		if err != nil {
			logger.Error("竞技场排期表时间格式错误")
			return nil
		}
		if time.Now().After(startTime) && time.Now().Before(endTime) { // 上期结束且下期未开启
			return v
		}
	}
	logger.Warn("找不到当前竞技场排期")
	return nil
}

func GetCurArenaSeason() *sro.ArenaSeasonExcelTable {
	info := GC.GetGPP().ArenaSeasonExcel
	if info.CurArenaSeason == nil {
		info.CurArenaSeason = genCurArenaSeason()
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", info.CurArenaSeason.SeasonEndDate)
	if err != nil && time.Now().After(endTime) {
		info.CurArenaSeason = genCurArenaSeason()
	}

	return info.CurArenaSeason
}

func GetArenaSeasonExcelTable(id int64) *sro.ArenaSeasonExcelTable {
	return GC.GetGPP().ArenaSeasonExcel.ArenaSeasonExcelMap[id]
}
