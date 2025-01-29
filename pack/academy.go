package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func AcademyGetInfo(s *enter.Session, request, response mx.Message) {
	// req := request.(*proto.AcademyGetInfoRequest)
	rsp := response.(*proto.AcademyGetInfoResponse)

	rsp.AcademyDB = game.GetAcademyDB(s)
	rsp.AcademyLocationDBs = game.GetAcademyLocationDBs(s)
}

func AcademyAttendSchedule(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.AcademyAttendScheduleRequest)
	rsp := response.(*proto.AcademyAttendScheduleResponse)

	parcelResultList := make([]*game.ParcelResult, 0)
	defer func() {
		rsp.AcademyDB = game.GetAcademyDB(s)
		rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResultList)
	}()
	conf := gdconf.GetAcademyZoneExcelTable(req.ZoneId)
	if conf == nil {
		return
	}
	locationInfo := game.GetAcademyLocationInfo(s, conf.LocationId)
	zoneInfo := game.GetAcademyZoneInfo(s, req.ZoneId)
	if locationInfo == nil ||
		locationInfo.Rank < conf.LocationRankForUnlock ||
		zoneInfo == nil || zoneInfo.IsUp {
		return
	}
	rewardConf := gdconf.GetAcademyRewardExcelTable(conf.RewardGroupId, locationInfo.Rank)
	if rewardConf == nil {
		return
	}
	// 添加学院经验
	parcelResultList = append(parcelResultList, &game.ParcelResult{
		ParcelType: proto.ParcelType_LocationExp,
		ParcelId:   conf.LocationId,
		Amount:     100,
	})
	parcelResultList = append(parcelResultList, &game.ParcelResult{
		ParcelType: proto.ParcelType_Currency,
		ParcelId:   proto.CurrencyTypes_AcademyTicket,
		Amount:     -1,
	})
	// 添加学生好感度
	for studentId, _ := range zoneInfo.StudentList {
		if game.GetCharacterInfo(s, studentId) != nil {
			parcelResultList = append(parcelResultList, &game.ParcelResult{
				ParcelType: proto.ParcelType_FavorExp,
				ParcelId:   studentId,
				Amount:     rewardConf.FavorExp,
			})
		}
	}
	// 获取可能的奖励
	parcelResultList = append(parcelResultList, game.GetParcelResultList(
		rewardConf.ExtraRewardParcelType, rewardConf.ExtraRewardParcelId,
		rewardConf.ExtraRewardAmount, false)...)
	zoneInfo.IsUp = true
}
