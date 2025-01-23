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
		CampaignState:                    proto.CampaignState_Win,
		CreateTime:                       mx.Now(),
		StageUniqueId:                    stageUniqueId,
		AccountServerId:                  s.AccountServerId,
		ActivatedHexaEventsAndConditions: make(map[int64][]int64),
		EnemyInfos:                       make(map[int64]*proto.HexaUnit),
		EnemyKillCountByUniqueId:         make(map[int64]int64),
		HexaEventDelayedExecutions:       make(map[int64][]int64),
		LastEnemyEntityId:                0,
		StageEntranceFee:                 make([]*proto.ParcelInfo, 0),
		StrategyObjects:                  make(map[int64]*proto.Strategy),
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
