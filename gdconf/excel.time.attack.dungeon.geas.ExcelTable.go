package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadTimeAttackDungeonGeasExcelTable() {
	g.GetExcel().TimeAttackDungeonGeasExcelTable = make([]*sro.TimeAttackDungeonGeasExcelTable, 0)
	name := "TimeAttackDungeonGeasExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().TimeAttackDungeonGeasExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetTimeAttackDungeonGeasExcelTable()))
}

type TimeAttackDungeonGeasExcel struct {
	TimeAttackDungeonGeasExcelMap map[int64]*sro.TimeAttackDungeonGeasExcelTable
}

func (g *GameConfig) gppTimeAttackDungeonGeasExcelTable() {
	g.GetGPP().TimeAttackDungeonGeasExcel = &TimeAttackDungeonGeasExcel{
		TimeAttackDungeonGeasExcelMap: make(map[int64]*sro.TimeAttackDungeonGeasExcelTable),
	}

	for _, v := range g.GetExcel().GetTimeAttackDungeonGeasExcelTable() {
		g.GetGPP().TimeAttackDungeonGeasExcel.TimeAttackDungeonGeasExcelMap[v.Id] = v
	}

	logger.Info("处理综合战术考试关卡配置完成,技能配置:%v个",
		len(g.GetGPP().TimeAttackDungeonGeasExcel.TimeAttackDungeonGeasExcelMap))
}

func GetTimeAttackDungeonGeasExcelTable(id int64) *sro.TimeAttackDungeonGeasExcelTable {
	return GC.GetGPP().TimeAttackDungeonGeasExcel.TimeAttackDungeonGeasExcelMap[id]
}
