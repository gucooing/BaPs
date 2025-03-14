package game

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/common/rank"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func GetEliminateRaidSeasonType() proto.RaidSeasonType {
	cur := gdconf.GetCurRaidEliminateSchedule()
	if cur == nil {
		return proto.RaidSeasonType_Close
	}
	if cur.StartTime.Before(time.Now()) &&
		cur.EndTime.After(time.Now()) {
		return proto.RaidSeasonType_Open
	}
	return proto.RaidSeasonType_Close // proto.RaidSeasonType_Settlement
}

func GetRaidEliminateBin(s *enter.Session) *sro.RaidEliminateBin {
	bin := GetPlayerBin(s)
	if bin == nil {
		return nil
	}
	if bin.RaidEliminateBin == nil {
		bin.RaidEliminateBin = &sro.RaidEliminateBin{}
	}
	return bin.RaidEliminateBin
}

func GetCurRaidEliminateInfo(s *enter.Session) *sro.RaidEliminateInfo {
	bin := GetRaidEliminateBin(s)
	conf := gdconf.GetCurRaidEliminateSchedule()
	if bin == nil || conf == nil {
		return nil
	}
	if bin.CurRaidEliminate == nil ||
		bin.GetCurRaidEliminate().GetSeasonId() != conf.SeasonId {
		bin.CurRaidEliminate = &sro.RaidEliminateInfo{
			SeasonId: conf.SeasonId,
		}

		bin.LastRaidEliminate = bin.CurRaidEliminate
	}

	// 如果赛季已经进入结算期
	if time.Now().After(conf.EndTime) {
		// 上面已经抛弃无效数据了,这里不需要处理
		bin.CurRaidEliminate.Ranking = rank.GetRaidEliminateRank(conf.SeasonId, s.AccountServerId)
		bin.CurRaidEliminate.Tier = gdconf.GetEliminateRaidTier(
			bin.CurRaidEliminate.SeasonId, bin.CurRaidEliminate.Ranking)
	}

	return bin.CurRaidEliminate
}

// GetLastRaidEliminateInfo 此处返回的可能是nil
func GetLastRaidEliminateInfo(s *enter.Session) *sro.RaidEliminateInfo {
	GetCurRaidEliminateInfo(s)

	return GetRaidEliminateBin(s).GetLastRaidEliminate()
}

// NewCurRaidEliminateBattleInfo 创建新的大决战
func NewCurRaidEliminateBattleInfo(s *enter.Session, raidUniqueId int64, isPractice bool) {
	bin := GetRaidEliminateBin(s)
	conf := gdconf.GetEliminateRaidStageExcelTable(raidUniqueId)
	if bin == nil || conf == nil {
		logger.Debug("玩家实例不存在或大决战关卡不存在RaidUniqueId:%v", raidUniqueId)
		return
	}
	chConf := gdconf.GetCharacterStatExcelTable(conf.RaidCharacterId)
	if chConf == nil {
		logger.Error("大决战boss实例不存在RaidCharacterId:%v", conf.RaidCharacterId)
		return
	}
	if !isPractice {
		// 扣票
		UpCurrency(s, conf.RaidEnterCostId, -conf.RaidEnterCostAmount)
	}
	bin.CurRaidBattleInfo = &sro.CurRaidBattleInfo{
		RaidUniqueId: raidUniqueId,
		IsPractice:   isPractice,
		RaidTeamList: make(map[int32]*sro.RaidTeamInfo),
		Frame:        0,
		Begin:        time.Now().Add(1 * time.Hour).Unix(),
		MaxHp:        chConf.MaxHP100,
		SeasonId:     GetCurRaidEliminateInfo(s).SeasonId,
		ServerId:     1,
		ContentType:  proto.ContentType_EliminateRaid,
	}
}

func RaidEliminateCheck(s *enter.Session) {
	bin := GetCurRaidEliminateInfo(s)
	if bin == nil {
		return
	}
	// 检查总分奖励领取
	seasonConf := gdconf.GetEliminateRaidSeasonManageExcelTable(bin.SeasonId)
	if seasonConf == nil ||
		len(seasonConf.StackedSeasonRewardGauge) != len(seasonConf.SeasonRewardId) {
		return
	}
	for index, season := range seasonConf.StackedSeasonRewardGauge {
		if _, ok := bin.ReceiveRewardIds[seasonConf.SeasonRewardId[index]]; !ok &&
			bin.TotalScore >= season {
			SetServerNotification(s, proto.ServerNotificationFlag_CanReceiveEliminateRaidReward, true)
			break
		}
	}
}

func GetCurRaidEliminateBattleInfo(s *enter.Session) *sro.CurRaidBattleInfo {
	return GetRaidEliminateBin(s).GetCurRaidBattleInfo()
}

func GetEliminateRaidLobbyInfoDB(s *enter.Session) *proto.EliminateRaidLobbyInfoDB {
	bin := GetCurRaidEliminateInfo(s)
	conf := gdconf.GetEliminateRaidSeasonManageExcelTable(bin.GetSeasonId())
	if conf == nil || bin == nil {
		return nil
	}
	info := &proto.EliminateRaidLobbyInfoDB{
		RaidLobbyInfoDB: &proto.RaidLobbyInfoDB{
			PlayingRaidDB: GetRaidDB(s, GetCurRaidEliminateBattleInfo(s)),
			PlayableHighestDifficulty: map[string]proto.Difficulty{
				conf.OpenRaidBossGroup01: proto.Difficulty_Lunatic,
				conf.OpenRaidBossGroup02: proto.Difficulty_Lunatic,
				conf.OpenRaidBossGroup03: proto.Difficulty_Lunatic,
			},
			ReceivedRankingRewardId: bin.GetRankingRewardId(),
			ReceiveRewardIds:        GetReceiveRewardIds(bin.GetReceiveRewardIds()),

			ReceiveLimitedRewardIds:  nil,
			ClanAssistUseInfo:        nil,
			SweepPointByRaidUniqueId: make(map[int64]int64),
			RemainFailCompensation: map[int32]bool{
				0: false,
				1: false,
				2: false,
			},
		},
		OpenedBossGroups: nil,
	}
	if cur := gdconf.GetCurRaidEliminateSchedule(); cur != nil {
		info.SeasonId = cur.SeasonId
		info.SeasonStartDate = mx.MxTime(cur.StartTime)
		info.SeasonEndDate = mx.MxTime(cur.EndTime)
		info.SettlementEndDate = mx.MxTime(cur.EndTime)
		info.Ranking = rank.GetRaidEliminateRank(cur.SeasonId, s.AccountServerId)
		info.Tier = gdconf.GetEliminateRaidTier(cur.SeasonId, info.Ranking)
		info.TotalRankingPoint = bin.GetTotalScore()
		info.CanReceiveRankingReward = GetCanReceiveRankingReward(
			time.Now().After(cur.EndTime), bin.GetIsRankingReward())
		info.BestRankingPointPerBossGroup = GetBestRankingPointPerBossGroup(s)
	}
	if next := gdconf.GetNextRaidEliminateSchedule(); next != nil {
		info.NextSeasonId = next.SeasonId
		info.NextSeasonStartDate = mx.MxTime(next.StartTime)
		info.NextSeasonEndDate = mx.MxTime(next.EndTime)
		info.SettlementEndDate = mx.MxTime(next.EndTime)
	}
	// 如果有进行中的战斗
	if curBattle := GetCurRaidEliminateBattleInfo(s); curBattle != nil &&
		!curBattle.IsClose {
		for _, teamInfo := range GetCurRaidEliminateBattleInfo(s).GetRaidTeamList() {
			for _, rcInfo := range teamInfo.MainCharacterList {
				info.ParticipateCharacterServerIds =
					append(info.ParticipateCharacterServerIds, rcInfo.ServerId)
			}
			for _, rcInfo := range teamInfo.SupportCharacterList {
				info.ParticipateCharacterServerIds =
					append(info.ParticipateCharacterServerIds, rcInfo.ServerId)
			}
		}
	}

	return info
}

func GetBestRankingPointPerBossGroup(s *enter.Session) map[string]int64 {
	list := make(map[string]int64)
	bin := GetCurRaidEliminateInfo(s)
	for k, v := range bin.GetBestScoreList() {
		list[k] = v.Score
	}

	return list
}

func GetEliminateRaidUserDB(s *enter.Session) *proto.EliminateRaidUserDB {
	info := &proto.EliminateRaidUserDB{
		RaidUserDB:              GetRaidEliminateUserDB(s),
		BossGroupToRankingPoint: GetBestRankingPointPerBossGroup(s),
	}

	return info
}

func GetRaidEliminateUserDB(s *enter.Session) *proto.RaidUserDB {
	curInfo := GetCurRaidEliminateInfo(s)
	if curInfo == nil {
		return nil
	}
	ranking, score := rank.GetRaidEliminateRankAndScore(curInfo.SeasonId, s.AccountServerId)
	info := &proto.RaidUserDB{
		AccountId:                   s.AccountServerId,
		RepresentCharacterUniqueId:  GetRepresentCharacterUniqueId(s),
		RepresentCharacterCostumeId: 0,
		Level:                       int64(GetAccountLevel(s)),
		Nickname:                    GetNickname(s),
		Tier:                        gdconf.GetEliminateRaidTier(curInfo.SeasonId, ranking),
		Rank:                        ranking,
		BestRankingPoint:            int64(score),
		BestRankingPointDetail:      score,
		AccountAttachmentDB:         GetAccountAttachmentDB(s),
	}

	return info
}

func GetRaidEliminateGiveUpDB(s *enter.Session) *proto.RaidGiveUpDB {
	curBattle := GetCurRaidEliminateBattleInfo(s)
	if curBattle == nil {
		return nil
	}
	ranking, bast := rank.GetRaidEliminateRankAndScore(curBattle.SeasonId, s.AccountServerId)
	info := &proto.RaidGiveUpDB{
		Ranking:          ranking,
		RankingPoint:     curBattle.ClearTimePoint + curBattle.HpScorePoint + curBattle.DefaultPoint,
		BestRankingPoint: int64(bast),
	}
	return info
}

func RaidEliminateClose(s *enter.Session) []*ParcelResult {
	curBattle := GetCurRaidEliminateBattleInfo(s)
	cur := GetCurRaidEliminateInfo(s)
	if curBattle == nil {
		return nil
	}
	conf := gdconf.GetEliminateRaidStageExcelTable(curBattle.RaidUniqueId)
	if conf == nil {
		return nil
	}
	list := make([]*ParcelResult, 0)
	// 计算分数
	curBattle.DefaultPoint = conf.DefaultClearScore
	curBattle.HpScorePoint = conf.HppercentScore * curBattle.GivenDamage / curBattle.MaxHp
	curBattle.ClearTimePoint = alg.MaxInt64(conf.MaximumScore-conf.PerSecondMinusScore/300*int64(curBattle.Frame), 0)

	// 如果不是模拟,且战斗结束
	if !curBattle.IsPractice && curBattle.IsClose && len(curBattle.RaidTeamList) > 0 {
		rankingPoint := curBattle.ClearTimePoint + curBattle.HpScorePoint + curBattle.DefaultPoint
		cur.TotalScore += rankingPoint // 累积分数
		if cur.BestScoreList == nil {
			cur.BestScoreList = make(map[string]*sro.RaidEliminateBest)
		}
		if cur.BestScoreList[conf.RaidBossGroup] == nil {
			cur.BestScoreList[conf.RaidBossGroup] = &sro.RaidEliminateBest{}
		}
		best := cur.BestScoreList[conf.RaidBossGroup]
		if best.Score < rankingPoint {
			best.Score = rankingPoint                  // 标记最高分
			best.RaidTeamList = curBattle.RaidTeamList // 标记最高分队伍
		}
		// 更新排名
		allBestScore := int64(0)
		for _, k := range cur.BestScoreList {
			allBestScore += k.Score
		}
		rank.SetRaidEliminateScore(curBattle.SeasonId, s.AccountServerId, float64(allBestScore))
		// 计算奖励
		for _, rewardConf := range gdconf.GetEliminateRaidStageRewardExcelTable(conf.RaidRewardGroupId) {
			list = append(list, &ParcelResult{
				ParcelType: proto.GetParcelTypeValue(rewardConf.ClearStageRewardParcelType),
				ParcelId:   rewardConf.ClearStageRewardParcelUniqueId,
				Amount:     rewardConf.ClearStageRewardAmount,
			})
		}
	}
	curBattle.IsClose = true
	return list
}
