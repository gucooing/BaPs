package game

import (
	"github.com/gucooing/BaPs/protocol/mx"
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/protocol/proto"
)

func GetAttendanceBin(s *enter.Session) *sro.AttendanceBin {
	bin := GetPlayerBin(s)
	if bin == nil {
		return nil
	}
	if bin.AttendanceBin == nil {
		bin.AttendanceBin = &sro.AttendanceBin{}
	}
	return bin.AttendanceBin
}

func GetAttendanceList(s *enter.Session) map[int64]*sro.AttendanceInfo {
	bin := GetAttendanceBin(s)
	if bin == nil {
		return nil
	}
	if bin.AttendanceList == nil {
		bin.AttendanceList = make(map[int64]*sro.AttendanceInfo)
	}
	for id, conf := range gdconf.GetAttendanceMap() {
		// 结束开始 不添加新的
		if time.Now().After(conf.StartableEndTime) {
			continue
		}
		// 开始,添加数据
		if time.Now().After(conf.StartTime) &&
			bin.AttendanceList[id] == nil {
			bin.AttendanceList[id] = &sro.AttendanceInfo{
				AttendanceId: id,
				ServerId:     GetServerId(s),
				AttendedDay:  make(map[int64]int64),
				LastReward:   0,
			}
		}
	}
	return bin.AttendanceList
}

func GetAttendanceInfo(s *enter.Session, attendanceId int64) *sro.AttendanceInfo {
	bin := GetAttendanceList(s)
	if bin == nil {
		return nil
	}
	info := bin[attendanceId]
	conf := gdconf.GetAttendanceInfo(attendanceId)
	if info == nil || conf == nil {
		delete(bin, attendanceId)
		return nil
	}
	if int64(len(info.AttendedDay)) >= conf.BookSize || time.Now().After(conf.EndTime) {
		// stop
		switch conf.Type {
		case "Basic":
			info.AttendedDay = make(map[int64]int64)
			info.LastReward = 0
		}
		info.Expired = true
	} else if time.Unix(info.LastReward, 0).After(alg.GetLastDayH(4)) {
		// day
		info.Expired = true
	} else {
		info.Expired = false
	}

	if info.Expired {
		return nil
	}

	return info
}

func GetAttendanceListDay(s *enter.Session) []*sro.AttendanceInfo {
	list := make([]*sro.AttendanceInfo, 0)
	for id, _ := range GetAttendanceList(s) {
		if bin := GetAttendanceInfo(s, id); bin != nil {
			list = append(list, bin)
		}
	}
	return list
}

func GetAttendanceBookRewards(s *enter.Session) []*proto.AttendanceBookReward {
	list := make([]*proto.AttendanceBookReward, 0)
	for _, bin := range GetAttendanceListDay(s) {
		info := GetAttendanceBookReward(s, bin.GetAttendanceId())
		if info == nil {
			continue
		}
		list = append(list, info)
	}
	return list
}

func GetAttendanceBookReward(s *enter.Session, attendanceId int64) *proto.AttendanceBookReward {
	bin := GetAttendanceInfo(s, attendanceId)
	conf := gdconf.GetAttendanceInfo(attendanceId)
	if bin == nil || conf == nil {
		return nil
	}
	info := &proto.AttendanceBookReward{
		UniqueId:          bin.AttendanceId,
		AccountType:       proto.AccountState_Normal.Value(conf.AccountType),
		BookSize:          conf.BookSize,
		MailType:          proto.MailType_System.Value(conf.MailType),
		Title:             conf.Title,
		TitleImagePath:    conf.TitleImagePath,
		AccountLevelLimit: conf.AccountLevelLimit,
		Type:              proto.AttendanceType_Basic.Value(conf.Type),
		StartDate:         mx.MxTime(conf.StartTime),
		StartableEndDate:  mx.MxTime(conf.StartableEndTime),
		EndDate:           mx.MxTime(conf.EndTime),
		DailyRewards:      make(map[int64][]*proto.ParcelInfo),
		DailyRewardIcons:  make(map[int64]string),
	}

	for index, rewardConf := range conf.AttendanceReward {
		reward := []*proto.ParcelInfo{
			{
				Key: &proto.ParcelKeyPair{
					Type: proto.ParcelType(rewardConf.RewardParcelType),
					Id:   rewardConf.RewardId,
				},
				Amount: rewardConf.RewardAmoun,
				Multiplier: &proto.BasisPoint{
					RawValue: 10000,
				},
				Probability: &proto.BasisPoint{
					RawValue: 10000,
				},
			},
		}
		info.DailyRewards[index] = reward

		info.DailyRewardIcons[index] = ""
	}
	return info
}

func GetAttendanceHistoryDBs(s *enter.Session) []*proto.AttendanceHistoryDB {
	list := make([]*proto.AttendanceHistoryDB, 0)
	for _, bin := range GetAttendanceListDay(s) {
		info := GetAttendanceHistoryDB(s, bin.GetAttendanceId())
		if info == nil {
			continue
		}
		list = append(list, info)
	}

	return list
}

func GetAttendanceHistoryDB(s *enter.Session, attendanceId int64) *proto.AttendanceHistoryDB {
	bin := GetAttendanceInfo(s, attendanceId)
	conf := gdconf.GetAttendanceInfo(attendanceId)
	if bin == nil || conf == nil {
		return nil
	}
	info := &proto.AttendanceHistoryDB{
		ServerId:               bin.ServerId,
		AttendanceBookUniqueId: bin.AttendanceId,
		AttendedDay:            make(map[int64]mx.MxTime),
		Expired:                bin.Expired,
	}
	for day, data := range bin.AttendedDay {
		info.AttendedDay[day] = mx.Unix(data, 0)
	}
	return info
}
