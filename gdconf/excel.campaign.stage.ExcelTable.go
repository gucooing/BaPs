package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadCampaignStageExcelTable() {
	g.GetExcel().CampaignStageExcelTable = make([]*sro.CampaignStageExcelTable, 0)
	name := "CampaignStageExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().CampaignStageExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetCampaignStageExcelTable()))
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

func GetCampaignStageExcelTable(stageId int64) *sro.CampaignStageExcelTable {
	return GC.GetGPP().CampaignStageExcel.CampaignStageExcelMap[stageId]
}

func GetCampaignStageExcelMap() map[int64]*sro.CampaignStageExcelTable {
	return GC.GetGPP().CampaignStageExcel.CampaignStageExcelMap
}
