package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
)

func (g *GameConfig) loadGachaElementExcelTable() {
	g.GetExcel().GachaElementExcelTable = make([]*sro.GachaElementExcelTable, 0)
	name := "GachaElementExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().GachaElementExcelTable)
}

type GachaElementExcel struct {
	GachaElementExcelMap    map[int64]*sro.GachaElementExcelTable
	GachaElementGroupIdList map[int64]*GachaElementGroupId // GroupId
}

type GachaElementGroupId struct {
	GachaGroupId          int64
	Rarity                string
	GachaElementExcelList []*sro.GachaElementExcelTable
}

func (g *GameConfig) gppGachaElementExcelTable() {
	g.GetGPP().GachaElementExcel = &GachaElementExcel{
		GachaElementExcelMap:    make(map[int64]*sro.GachaElementExcelTable),
		GachaElementGroupIdList: make(map[int64]*GachaElementGroupId),
	}
	for _, v := range g.GetExcel().GetGachaElementExcelTable() {
		g.GetGPP().GachaElementExcel.GachaElementExcelMap[v.Id] = v

		if g.GetGPP().GachaElementExcel.GachaElementGroupIdList[v.GachaGroupId] == nil {
			g.GetGPP().GachaElementExcel.GachaElementGroupIdList[v.GachaGroupId] = &GachaElementGroupId{
				GachaGroupId:          v.GachaGroupId,
				Rarity:                v.Rarity,
				GachaElementExcelList: make([]*sro.GachaElementExcelTable, 0),
			}
		}
		g.GetGPP().GachaElementExcel.GachaElementGroupIdList[v.GachaGroupId].GachaElementExcelList = append(
			g.GetGPP().GachaElementExcel.GachaElementGroupIdList[v.GachaGroupId].GachaElementExcelList,
			v,
		)
	}

	logger.Info("处理随机组池配置完成,成就:%v个",
		len(g.GetGPP().GachaElementExcel.GachaElementExcelMap))
}

func GetGachaElementExcelTable(id int64) *sro.GachaElementExcelTable {
	return GC.GetGPP().GachaElementExcel.GachaElementExcelMap[id]
}

func GetGachaElementGroupIdByGachaGroupId(gachaGroupId int64) *GachaElementGroupId {
	return GC.GetGPP().GachaElementExcel.GachaElementGroupIdList[gachaGroupId]
}
