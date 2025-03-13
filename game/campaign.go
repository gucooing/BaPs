package game

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func NewCampaignMainStageSaveDB(s *enter.Session, stageUniqueId int64) *proto.CampaignMainStageSaveDB {
	stageConf := gdconf.GetCampaignStageExcelTable(stageUniqueId)
	if stageConf == nil {
		logger.Debug("Unknown StageUniqueId:%v", stageUniqueId)
		return nil
	}
	info := &proto.CampaignMainStageSaveDB{
		ContentType:                      proto.ContentType(proto.ContentType_value[stageConf.ContentType]),
		CampaignState:                    proto.CampaignState_BeforeStart, //  proto.CampaignState_Win,
		CreateTime:                       mx.Now(),
		StageUniqueId:                    stageUniqueId,
		AccountServerId:                  s.AccountServerId,
		ActivatedHexaEventsAndConditions: make(map[int64][]int64),
		EnemyInfos:                       GetEnemyInfos(stageUniqueId),
		EnemyKillCountByUniqueId:         make(map[int64]int64),
		HexaEventDelayedExecutions:       make(map[int64][]int64),
		LastEnemyEntityId:                10010,
		StageEntranceFee:                 make([]*proto.ParcelInfo, 0),
		StrategyObjects:                  GetStrategyObjects(stageUniqueId),
		TileMapStates:                    make(map[int32]*proto.HexaTileState),

		CurrentTurn:           0,
		EnemyClearCount:       0,
		TacticRankSCount:      0,
		EchelonInfos:          nil,
		WithdrawInfos:         nil,
		StrategyObjectRewards: nil,
		StrategyObjectHistory: nil,
		DisplayInfos:          nil,
		DeployedEchelonInfos:  nil,
	}

	return info
}

func GetEnemyInfos(StageUniqueId int64) map[int64]*proto.HexaUnit {
	list := make(map[int64]*proto.HexaUnit)
	// list[10008] = &proto.HexaUnit{
	// 	EntityId: 10008,
	// 	Id:       101110101,
	// 	Rotate: &proto.Vector3{
	// 		X: 0,
	// 		Y: 240,
	// 		Z: 0,
	// 	},
	// 	Location: &proto.Vector3{
	// 		X: -1,
	// 		Y: 1,
	// 		Z: 0,
	// 	},
	// }
	// list[10009] = &proto.HexaUnit{
	// 	EntityId: 10009,
	// 	Id:       101110102,
	// 	Rotate: &proto.Vector3{
	// 		X: 0,
	// 		Y: 240,
	// 		Z: 0,
	// 	},
	// 	Location: &proto.Vector3{
	// 		X: 0,
	// 		Y: 0,
	// 		Z: 0,
	// 	},
	// }

	return list
}

func GetStrategyObjects(StageUniqueId int64) map[int64]*proto.Strategy {
	list := make(map[int64]*proto.Strategy)
	// list[10010] = &proto.Strategy{
	// 	EntityId: 10010,
	// 	Id:       101101,
	// 	Rotate: &proto.Vector3{
	// 		X: 0,
	// 		Y: 0,
	// 		Z: 0,
	// 	},
	// 	Location: &proto.Vector3{
	// 		X: -2,
	// 		Y: 2,
	// 		Z: 0,
	// 	},
	// }

	return list
}
