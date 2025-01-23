package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadCampaignUnitExcelTable() {
	g.GetExcel().CampaignUnitExcelTable = make([]*sro.CampaignUnitExcelTable, 0)
	name := "CampaignUnitExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().CampaignUnitExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetCampaignUnitExcelTable()))
}

type CampaignUnitExcel struct {
	CampaignUnitExcelMap      map[int64]*sro.CampaignUnitExcelTable
	CampaignUnitExcelStageMap map[int64]*CampaignUnitExcelGrade
}

type CampaignUnitExcelGrade struct {
	Boss      *sro.CampaignUnitExcelTable
	GradeList map[string]*sro.CampaignUnitExcelTable
}

func (g *GameConfig) gppCampaignUnitExcelTable() {
	info := &CampaignUnitExcel{
		CampaignUnitExcelMap:      make(map[int64]*sro.CampaignUnitExcelTable, 0),
		CampaignUnitExcelStageMap: make(map[int64]*CampaignUnitExcelGrade),
	}

	for _, v := range g.GetExcel().GetCampaignUnitExcelTable() {
		info.CampaignUnitExcelMap[v.Id] = v
		stageId := v.Id / 100
		if info.CampaignUnitExcelStageMap[stageId] == nil {
			info.CampaignUnitExcelStageMap[stageId] = &CampaignUnitExcelGrade{
				GradeList: make(map[string]*sro.CampaignUnitExcelTable, 0),
			}
		}
		if v.Grade == "Boss" {
			info.CampaignUnitExcelStageMap[stageId].Boss = v
		} else {
			info.CampaignUnitExcelStageMap[stageId].GradeList[v.Grade] = v
		}
	}

	g.GetGPP().CampaignUnitExcel = info

	logger.Info("任务关卡怪物信息关卡数量完成:%v个", len(g.GetGPP().CampaignUnitExcel.CampaignUnitExcelStageMap))
}
