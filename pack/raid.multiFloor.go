package pack

import (
	"github.com/gucooing/BaPs/protocol/mx"
	"time"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/protocol/proto"
)

func MultiFloorRaidSync(s *enter.Session, request, response mx.Message) {
	// req := request.(*proto.MultiFloorRaidSyncRequest)
	rsp := response.(*proto.MultiFloorRaidSyncResponse)

	rsp.MultiFloorRaidDBs = make([]*proto.MultiFloorRaidDB, 0)

	rsp.MultiFloorRaidDBs = game.GetMultiFloorRaidDBs(s)
}

func MultiFloorRaidEnterBattle(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.MultiFloorRaidEnterBattleRequest)
	rsp := response.(*proto.MultiFloorRaidEnterBattleResponse)

	rsp.AssistCharacterDBs = make([]*proto.AssistCharacterDB, 0)
	for _, assist := range req.AssistUseInfos {
		ac := enter.GetSessionByUid(assist.CharacterAccountId)
		assistInfo := game.GetAssistInfoByClanAssistUseInfo(ac, assist)
		if assistInfo != nil {
			rsp.AssistCharacterDBs = append(rsp.AssistCharacterDBs,
				game.GetAssistCharacterDB(ac, assistInfo, assist.AssistRelation))
		}
	}
}

func MultiFloorRaidEndBattle(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.MultiFloorRaidEndBattleRequest)
	rsp := response.(*proto.MultiFloorRaidEndBattleResponse)

	summary := req.Summary
	bin := game.GetCurRaidMultiFloorInfo(s)
	if summary == nil || bin == nil {
		return
	}
	// 不验是否合法了,随意
	isWin := false
	for _, raidBossResult := range summary.RaidSummary.RaidBossResults {
		if raidBossResult.EndHpRateRawValue == 0 {
			isWin = true
			break
		}
	}
	defer func() {
		rsp.MultiFloorRaidDB = game.GetMultiFloorRaidDB(s)
	}()
	if !isWin {
		return
	}
	// 更新数据
	if bin.ClearedDifficulty < req.Difficulty {
		bin.ClearedDifficulty = req.Difficulty
		bin.LastClearDate = time.Now().Unix()
		bin.Frame = summary.EndFrame
	} else if bin.ClearedDifficulty == req.Difficulty &&
		bin.Frame > summary.EndFrame {
		bin.LastClearDate = time.Now().Unix()
		bin.Frame = summary.EndFrame
	}
}

func MultiFloorRaidReceiveReward(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.MultiFloorRaidReceiveRewardRequest)
	rsp := response.(*proto.MultiFloorRaidReceiveRewardResponse)

	defer func() {
		rsp.MultiFloorRaidDB = game.GetMultiFloorRaidDB(s)
	}()

	bin := game.GetCurRaidMultiFloorInfo(s)
	if bin == nil {
		return
	}
	parcelResultList := make([]*game.ParcelResult, 0)
	for {
		bin.RewardDifficulty++
		confList := gdconf.GetMultiFloorRaidRewardExcelBySeasonId(req.SeasonId, bin.RewardDifficulty)
		if len(confList) == 0 {
			break
		}
		for _, conf := range confList {
			parcelResultList = append(parcelResultList, &game.ParcelResult{
				ParcelType: proto.ParcelType_None.Value(conf.ClearStageRewardParcelType),
				ParcelId:   conf.ClearStageRewardParcelUniqueID,
				Amount:     conf.ClearStageRewardAmount,
			})
		}
		if bin.RewardDifficulty == req.RewardDifficulty {
			break
		}
	}
	rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResultList)
	bin.LastRewardDate = time.Now().Unix()
}
