package game

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func NewCampaignMainStageSaveDB(s *enter.Session, stageUniqueId int64) *sro.BattleCampaign {
	stageConf := gdconf.GetCampaignStageExcelTable(stageUniqueId)
	if stageConf == nil {
		logger.Debug("Unknown StageUniqueId:%v", stageUniqueId)
		return nil
	}
	info := &sro.BattleCampaign{
		StageId:           stageUniqueId,
		CreateTime:        time.Now().Unix(),
		LastEnemyEntityId: 0,
		EnemyInfos:        make(map[int64]*sro.EnemyInfo),
		StrategyObjects:   make(map[int64]*sro.EnemyInfo),
		CampaignState:     proto.CampaignState_BeforeStart,
	}
	if mapConf := gdconf.GetStrategyMap(stageConf.StrategyMap); mapConf != nil {
		info.LastEnemyEntityId = mapConf.LastEnemyEntityId
		// for k, v := range mapConf.EnemyInfos {
		//
		// }
	}
	return info
}

func GetCampaignMainStageSaveDB(s *enter.Session, bin *sro.BattleCampaign) *proto.CampaignMainStageSaveDB {
	stageConf := gdconf.GetCampaignStageExcelTable(bin.StageId)
	if stageConf == nil {
		return nil
	}

	info := &proto.CampaignMainStageSaveDB{
		ContentType:       proto.ContentType(proto.ContentType_value[stageConf.ContentType]),
		CampaignState:     proto.CampaignState(bin.CampaignState),
		CreateTime:        mx.Unix(bin.CreateTime, 0),
		StageUniqueId:     bin.StageId,
		AccountServerId:   s.AccountServerId,
		EnemyInfos:        make(map[int64]*proto.HexaUnit),
		StrategyObjects:   make(map[int64]*proto.Strategy),
		LastEnemyEntityId: int32(bin.LastEnemyEntityId),

		ActivatedHexaEventsAndConditions: make(map[int64][]int64),
		EnemyKillCountByUniqueId:         make(map[int64]int64),
		HexaEventDelayedExecutions:       make(map[int64][]int64),
		StageEntranceFee:                 make([]*proto.ParcelInfo, 0),
		TileMapStates:                    make(map[int32]*proto.HexaTileState),
		CurrentTurn:                      0,
		EnemyClearCount:                  0,
		TacticRankSCount:                 0,
		EchelonInfos:                     nil,
		WithdrawInfos:                    nil,
		StrategyObjectRewards:            nil,
		StrategyObjectHistory:            nil,
		DisplayInfos:                     nil,
		DeployedEchelonInfos:             nil,
	}
	for _, v := range bin.EnemyInfos {
		info.EchelonInfos[v.EntityId] = GetEnemyInfo(v)
	}
	for _, v := range bin.StrategyObjects {
		info.StrategyObjects[v.EntityId] = GetStrategyObject(v)
	}
	return info
}

func GetEnemyInfo(bin *sro.EnemyInfo) *proto.HexaUnit {
	info := &proto.HexaUnit{
		EntityId:      bin.EntityId,
		Id:            bin.Id,
		Rotate:        GetVector3(bin.Rotate),
		Location:      GetVector3(bin.Location),
		BuffGroupIds:  nil,
		SkillCardHand: nil,
		PlayAnimation: false,
		RewardItems:   nil,
	}

	return info
}

func GetStrategyObject(bin *sro.EnemyInfo) *proto.Strategy {
	info := &proto.Strategy{
		EntityId:      bin.EntityId,
		Rotate:        GetVector3(bin.Rotate),
		Location:      GetVector3(bin.Location),
		Id:            bin.EntityId,
		PlayAnimation: false,
		Activated:     false,
		Values:        nil,
		Index:         0,
		Movable:       false,
		NeedValueType: false,
	}

	return info
}

func GetVector3(x *sro.Vector3) *proto.Vector3 {
	if x == nil {
		return nil
	}
	return &proto.Vector3{
		X: x.X,
		Y: x.Y,
		Z: x.Z,
	}
}
