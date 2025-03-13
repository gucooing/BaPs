package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadTimeAttackDungeonRewardExcelTable() {
	g.GetExcel().TimeAttackDungeonRewardExcelTable = make([]*sro.TimeAttackDungeonRewardExcelTable, 0)
	name := "TimeAttackDungeonRewardExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().TimeAttackDungeonRewardExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetTimeAttackDungeonRewardExcelTable()))
}

type TimeAttackDungeonRewardExcel struct {
	TimeAttackDungeonRewardExcelMap map[int64]*sro.TimeAttackDungeonRewardExcelTable
}

func (g *GameConfig) gppTimeAttackDungeonRewardExcelTable() {
	g.GetGPP().TimeAttackDungeonRewardExcel = &TimeAttackDungeonRewardExcel{
		TimeAttackDungeonRewardExcelMap: make(map[int64]*sro.TimeAttackDungeonRewardExcelTable),
	}

	for _, v := range g.GetExcel().GetTimeAttackDungeonRewardExcelTable() {
		g.GetGPP().TimeAttackDungeonRewardExcel.TimeAttackDungeonRewardExcelMap[v.Id] = v
	}

	logger.Info("处理综合战术考试关卡配置完成,技能配置:%v个",
		len(g.GetGPP().TimeAttackDungeonRewardExcel.TimeAttackDungeonRewardExcelMap))
}

func GetTimeAttackDungeonRewardExcelTable(id int64) *sro.TimeAttackDungeonRewardExcelTable {
	return GC.GetGPP().TimeAttackDungeonRewardExcel.TimeAttackDungeonRewardExcelMap[id]
}
