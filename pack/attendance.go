package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func AttendanceReward(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.AttendanceRewardRequest)
	rsp := response.(*proto.AttendanceRewardResponse)

	rsp.AttendanceBookRewards = make([]*proto.AttendanceBookReward, 0)
	rsp.AttendanceHistoryDBs = make([]*proto.AttendanceHistoryDB, 0)
	for id, day := range req.DayByBookUniqueId {
		bin := game.GetAttendanceInfo(s, id)
		if _, ok := bin.GetAttendedDay()[day]; ok {
			continue
		}
		conf := gdconf.GetAttendanceInfo(id)
		if conf == nil {
			continue
		}
		if rewardConf, ok := conf.AttendanceReward[day]; ok {
			parcelInfoList := []*sro.ParcelInfo{
				{
					Type: rewardConf.RewardParcelType,
					Id:   rewardConf.RewardId,
					Num:  rewardConf.RewardAmoun,
				},
			}
			game.AddMailBySystem(s, conf.MailType, parcelInfoList)
		}

		rsp.AttendanceBookRewards = append(rsp.AttendanceBookRewards,
			game.GetAttendanceBookReward(s, id))
		rsp.AttendanceHistoryDBs = append(rsp.AttendanceHistoryDBs,
			game.GetAttendanceHistoryDB(s, id))
	}
}
