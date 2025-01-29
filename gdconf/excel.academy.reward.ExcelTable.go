package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadAcademyRewardExcelTable() {
	g.GetExcel().AcademyRewardExcelTable = make([]*sro.AcademyRewardExcelTable, 0)
	name := "AcademyRewardExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().AcademyRewardExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetAcademyRewardExcelTable()))
}

type AcademyRewardExcel struct {
	AcademyRewardExcelMap map[int64]map[int64]*sro.AcademyRewardExcelTable
}

func (g *GameConfig) gppAcademyRewardExcelTable() {
	g.GetGPP().AcademyRewardExcel = &AcademyRewardExcel{
		AcademyRewardExcelMap: make(map[int64]map[int64]*sro.AcademyRewardExcelTable),
	}
	for _, v := range g.GetExcel().GetAcademyRewardExcelTable() {
		if g.GetGPP().AcademyRewardExcel.AcademyRewardExcelMap[v.ScheduleGroupId] == nil {
			g.GetGPP().AcademyRewardExcel.AcademyRewardExcelMap[v.ScheduleGroupId] = make(map[int64]*sro.AcademyRewardExcelTable)
		}
		g.GetGPP().AcademyRewardExcel.AcademyRewardExcelMap[v.ScheduleGroupId][v.LocationRank] = v
	}

	logger.Info("处理课程表奖励配置完成,数量:%v个",
		len(g.GetGPP().AcademyRewardExcel.AcademyRewardExcelMap))
}

func GetAcademyRewardExcelTable(gId int64, rank int64) *sro.AcademyRewardExcelTable {
	confList := GC.GetGPP().AcademyRewardExcel.AcademyRewardExcelMap[gId]
	if confList == nil {
		return nil
	}
	loadRank := int64(0)
	for _, conf := range confList {
		if conf.LocationRank <= rank &&
			loadRank < conf.LocationRank {
			loadRank = conf.LocationRank
		}
	}
	return confList[loadRank]
}
