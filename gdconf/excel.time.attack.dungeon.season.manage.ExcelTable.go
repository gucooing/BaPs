package gdconf

import (
	"encoding/json"
	"os"
	"time"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
)

func (g *GameConfig) loadTimeAttackDungeonSeasonManageExcelTable() {
	g.GetExcel().TimeAttackDungeonSeasonManageExcelTable = make([]*sro.TimeAttackDungeonSeasonManageExcelTable, 0)
	name := "TimeAttackDungeonSeasonManageExcelTable.json"
	file, err := os.ReadFile(g.excelPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetExcel().TimeAttackDungeonSeasonManageExcelTable); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	logger.Info("文件:%s 读取成功,解析数量:%v", name, len(g.GetExcel().GetTimeAttackDungeonSeasonManageExcelTable()))
}

type TimeAttackDungeonSeasonManageExcel struct {
	Cur                                   *sro.TimeAttackDungeonSeasonManageExcelTable
	TimeAttackDungeonSeasonManageExcelMap map[int64]*sro.TimeAttackDungeonSeasonManageExcelTable
}

func (g *GameConfig) gppTimeAttackDungeonSeasonManageExcelTable() {
	g.GetGPP().TimeAttackDungeonSeasonManageExcel = &TimeAttackDungeonSeasonManageExcel{
		TimeAttackDungeonSeasonManageExcelMap: make(map[int64]*sro.TimeAttackDungeonSeasonManageExcelTable, 0),
	}

	for _, v := range g.GetExcel().GetTimeAttackDungeonSeasonManageExcelTable() {
		g.GetGPP().TimeAttackDungeonSeasonManageExcel.TimeAttackDungeonSeasonManageExcelMap[v.Id] = v
	}

	logger.Info("处理综合战术考试排期表完成,数量:%v个", len(g.GetGPP().TimeAttackDungeonSeasonManageExcel.TimeAttackDungeonSeasonManageExcelMap))
}

func GetCurTimeAttackDungeonSeasonManageExcelTable() *sro.TimeAttackDungeonSeasonManageExcelTable {
	conf := GC.GetGPP().TimeAttackDungeonSeasonManageExcel
	if conf == nil {
		return nil
	}

	getCur := func() {
		for _, v := range conf.TimeAttackDungeonSeasonManageExcelMap {
			startTime, err := time.Parse("2006-01-02 15:04:05", v.StartDate)
			endTime, err := time.Parse("2006-01-02 15:04:05", v.EndDate)
			if err != nil {
				logger.Error("综合战术考试排期表时间格式错误")
				return
			}
			if time.Now().After(startTime) && time.Now().Before(endTime) {
				conf.Cur = v
				return
			}
		}
		logger.Warn("找不到当前综合战术考试排期")
	}

	if conf.Cur == nil {
		getCur()
	}
	if conf.Cur != nil {
		startTime, err := time.Parse("2006-01-02 15:04:05", conf.Cur.StartDate)
		endTime, err := time.Parse("2006-01-02 15:04:05", conf.Cur.EndDate)
		if err != nil {
			logger.Error("综合战术考试排期表时间格式错误")
			return nil
		}
		if !time.Now().After(startTime) || !time.Now().Before(endTime) {
			getCur()
		}
	}
	return conf.Cur
}
