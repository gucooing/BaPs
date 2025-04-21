package gdconf

import (
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
)

func (g *GameConfig) loadCampaignStageExcelTable() {
	g.GetExcel().CampaignStageExcelTable = make([]*sro.CampaignStageExcelTable, 0)
	name := "CampaignStageExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().CampaignStageExcelTable)
}

type CampaignStageExcel struct {
	CampaignStageExcelMap map[int64]*sro.CampaignStageExcelTable
}

func (g *GameConfig) gppCampaignStageExcelTable() {
	g.GetGPP().CampaignStageExcel = &CampaignStageExcel{
		CampaignStageExcelMap: make(map[int64]*sro.CampaignStageExcelTable, 0),
	}

	for _, v := range g.GetExcel().GetCampaignStageExcelTable() {
		g.GetGPP().CampaignStageExcel.CampaignStageExcelMap[v.Id] = v
	}

	logger.Info("处理任务关卡信息数量完成:%v个", len(g.GetGPP().CampaignStageExcel.CampaignStageExcelMap))
}

func GetCampaignStageExcelTable(id int64) *sro.CampaignStageExcelTable {
	return GC.GetGPP().CampaignStageExcel.CampaignStageExcelMap[id]
}

func GetCampaignStageExcelMap() map[int64]*sro.CampaignStageExcelTable {
	return GC.GetGPP().CampaignStageExcel.CampaignStageExcelMap
}
