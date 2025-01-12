package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/proto"
)

func EliminateRaidLogin(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.EliminateRaidLoginResponse)

	rsp.SeasonType = proto.RaidSeasonType_Close
	rsp.ServerNotification = proto.ServerNotificationFlag_HasUnreadMail
}

func MultiFloorRaidSync(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.MultiFloorRaidSyncRequest)
	rsp := response.(*proto.MultiFloorRaidSyncResponse)

	rsp.MultiFloorRaidDBs = make([]*proto.MultiFloorRaidDB, 0)

	if req.SeasonId == 0 {
		rsp.MultiFloorRaidDBs = append(rsp.MultiFloorRaidDBs, &proto.MultiFloorRaidDB{
			SeasonId:          10,
			ClearedDifficulty: -1,
		})
	}
}

func ContentSweepMultiSweepPresetList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.ContentSweepMultiSweepPresetListResponse)

	rsp.MultiSweepPresetDBs = make([]*proto.MultiSweepPresetDB, 0)
}
