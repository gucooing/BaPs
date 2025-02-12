package gdconf

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gucooing/BaPs/pkg/logger"
)

type RaidEliminateSchedule struct {
	CurRaidEliminateSchedule *RaidEliminateScheduleInfo
	RaidEliminateScheduleMap map[int64]*RaidEliminateScheduleInfo
}

type RaidEliminateScheduleInfo struct {
	SeasonId        int64
	SeasonStartData string    `json:"SeasonStartData"`
	StartTime       time.Time `json:"-"`
	SeasonEndData   string    `json:"SeasonEndData"`
	EndTime         time.Time `json:"-"`
	NextSeasonId    int64
}

func (g *GameConfig) loadRaidEliminateSchedule() {
	g.GetGPP().RaidEliminateSchedule = &RaidEliminateSchedule{
		RaidEliminateScheduleMap: make(map[int64]*RaidEliminateScheduleInfo),
	}
	name := "RaidEliminateSchedule.json"
	file, err := os.ReadFile(g.dataPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetGPP().RaidEliminateSchedule.RaidEliminateScheduleMap); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}
	for _, v := range g.GetGPP().RaidEliminateSchedule.RaidEliminateScheduleMap {
		if v.NextSeasonId != 0 &&
			g.GetGPP().RaidEliminateSchedule.RaidEliminateScheduleMap[v.NextSeasonId] == nil {
			panic(fmt.Sprintf("缺少下一个大决战排期,NextSeasonId:%v", v.NextSeasonId))
		}
		v.StartTime, err = time.Parse("2006-01-02 15:04:05", v.SeasonStartData)
		v.EndTime, err = time.Parse("2006-01-02 15:04:05", v.SeasonEndData)
		if err != nil {
			panic(fmt.Sprintf("大决战排期时间格式错误,SeasonId:%v", v.SeasonId))
		}
		// 检查排期时间是否满足7天以上
		if v.StartTime.Add(7 * 24 * time.Hour).After(v.EndTime) {
			panic(fmt.Sprintf("大决战排期时间错误 排期不足7天,SeasonId:%v", v.SeasonId))
		}
	}
	logger.Info("大决战排期读取成功文件:%s 读取成功,解析数量:%v", name, len(g.GetGPP().RaidEliminateSchedule.RaidEliminateScheduleMap))
}

func GetCurRaidEliminateSchedule() *RaidEliminateScheduleInfo {
	cur := GC.GetGPP().RaidEliminateSchedule.CurRaidEliminateSchedule
	if cur != nil {
		if cur.EndTime.After(time.Now()) {
			return GC.GetGPP().RaidEliminateSchedule.CurRaidEliminateSchedule
		}
		next := GC.GetGPP().RaidEliminateSchedule.RaidEliminateScheduleMap[cur.NextSeasonId]
		if next != nil && next.StartTime.After(time.Now()) {
			return GC.GetGPP().RaidEliminateSchedule.CurRaidEliminateSchedule
		}
	}
	cur = nil
	for _, v := range GC.GetGPP().RaidEliminateSchedule.RaidEliminateScheduleMap {
		// 排期中
		if v.EndTime.After(time.Now()) && time.Now().After(v.StartTime) {
			cur = v
		} else { // 排期间隔中
			next := GC.GetGPP().RaidEliminateSchedule.RaidEliminateScheduleMap[v.NextSeasonId]
			if next != nil &&
				time.Now().After(v.EndTime) && next.StartTime.After(time.Now()) {
				cur = v
			}
		}
		if cur != nil {
			GC.GetGPP().RaidEliminateSchedule.CurRaidEliminateSchedule = cur
			return cur
		}
	}
	return cur
}

func GetNextRaidEliminateSchedule() *RaidEliminateScheduleInfo {
	cur := GetCurRaidEliminateSchedule()
	if cur == nil {
		return nil
	}
	return GC.GetGPP().RaidEliminateSchedule.RaidEliminateScheduleMap[cur.NextSeasonId]
}

func GetRaidEliminateScheduleMap() map[int64]*RaidEliminateScheduleInfo {
	return GC.GetGPP().RaidEliminateSchedule.RaidEliminateScheduleMap
}

func GetRaidEliminateScheduleInfo(seasonId int64) *RaidEliminateScheduleInfo {
	return GC.GetGPP().RaidEliminateSchedule.RaidEliminateScheduleMap[seasonId]
}
