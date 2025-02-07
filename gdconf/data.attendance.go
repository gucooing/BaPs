package gdconf

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gucooing/BaPs/pkg/logger"
)

type Attendance struct {
	AttendanceMap map[int64]*AttendanceInfo
}

type AttendanceInfo struct {
	Id                int64                       `json:"Id"`
	Type              string                      `json:"Type"`
	AccountType       string                      `json:"AccountType"`
	AccountLevelLimit int64                       `json:"AccountLevelLimit"`
	Title             string                      `json:"Title"`
	TitleImagePath    string                      `json:"TitleImagePath"`
	CountRule         string                      `json:"CountRule"`
	CountReset        string                      `json:"CountReset"`
	BookSize          int64                       `json:"BookSize"`
	StartTime         time.Time                   `json:"-"`
	StartDate         string                      `json:"StartDate"`
	StartableEndTime  time.Time                   `json:"-"`
	StartableEndDate  string                      `json:"StartableEndDate"`
	EndTime           time.Time                   `json:"-"`
	EndDate           string                      `json:"EndDate"`
	MailType          string                      `json:"MailType"`
	AttendanceReward  map[int64]*AttendanceReward `json:"AttendanceReward"`
}

type AttendanceReward struct {
	RewardParcelType int32 `json:"RewardParcelType"`
	RewardId         int64 `json:"RewardId"`
	RewardAmoun      int64 `json:"RewardAmoun"`
}

func (g *GameConfig) loadAttendance() {
	g.GetGPP().Attendance = &Attendance{
		AttendanceMap: make(map[int64]*AttendanceInfo),
	}
	name := "Attendance.json"
	file, err := os.ReadFile(g.dataPath + name)
	if err != nil {
		logger.Error("文件:%s 读取失败,err:%s", name, err)
		return
	}
	if err := json.Unmarshal(file, &g.GetGPP().Attendance.AttendanceMap); err != nil {
		logger.Error("文件:%s 解析失败,err:%s", name, err)
		return
	}

	for _, v := range g.GetGPP().Attendance.AttendanceMap {
		v.StartTime, err = time.Parse("2006-01-02 15:04:05", v.StartDate)
		v.EndTime, err = time.Parse("2006-01-02 15:04:05", v.EndDate)
		v.StartableEndTime, err = time.Parse("2006-01-02 15:04:05", v.StartableEndDate)
		if err != nil {
			panic(fmt.Sprintf("登录奖励排期时间格式错误,Id:%v", v.Id))
		}
		// 检查排期奖励是否合法
		if v.BookSize != int64(len(v.AttendanceReward)) {
			panic(fmt.Sprintf("登录奖励错误 排期奖励不足,Id:%v", v.Id))
		}
	}
	logger.Info("登录奖励排期读取成功文件:%s 读取成功,解析数量:%v", name, len(g.GetGPP().Attendance.AttendanceMap))
}

func GetAttendanceMap() map[int64]*AttendanceInfo {
	return GC.GetGPP().Attendance.AttendanceMap
}

func GetAttendanceInfo(id int64) *AttendanceInfo {
	return GC.GetGPP().Attendance.AttendanceMap[id]
}
