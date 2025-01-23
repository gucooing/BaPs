package gdconf

import (
	"encoding/json"
	"os"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadSchoolDungeonStageExcel() {
	g.GetExcel().SchoolDungeonStageExcel = make([]*sro.SchoolDungeonStageExcel, 0)
	name := "SchoolDungeonStageExcel.json"
	file, err := os.ReadFile(g.excelDbPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().SchoolDungeonStageExcel); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetSchoolDungeonStageExcel()))
}

type SchoolDungeonStage struct {
	SchoolDungeonStageMap map[int64]*sro.SchoolDungeonStageExcel
}

func (g *GameConfig) gppSchoolDungeonStageExcel() {
	g.GetGPP().SchoolDungeonStage = &SchoolDungeonStage{
		SchoolDungeonStageMap: make(map[int64]*sro.SchoolDungeonStageExcel, 0),
	}

	for _, v := range g.GetExcel().GetSchoolDungeonStageExcel() {
		g.GetGPP().SchoolDungeonStage.SchoolDungeonStageMap[v.StageId] = v
	}

	logger.Info("学院交流会关卡信息数量完成:%v个", len(g.GetGPP().SchoolDungeonStage.SchoolDungeonStageMap))
}

func GetSchoolDungeonStageExcel(stageId int64) *sro.SchoolDungeonStageExcel {
	return GC.GetGPP().SchoolDungeonStage.SchoolDungeonStageMap[stageId]
}
