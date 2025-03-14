package pack

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/common/rank"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/proto"
)

func EliminateRaidLogin(s *enter.Session, request, response proto.Message) {
	rsp := response.(*proto.EliminateRaidLoginResponse)

	game.RaidEliminateCheck(s)
	rsp.SeasonType = game.GetEliminateRaidSeasonType()
	rsp.SweepPointByRaidUniqueId = make(map[int64]int64) // 扫荡信息
}

func EliminateRaidLobby(s *enter.Session, request, response proto.Message) {
	rsp := response.(*proto.EliminateRaidLobbyResponse)

	curBattle := game.GetCurRaidEliminateBattleInfo(s)
	// 超时了
	if curBattle != nil &&
		!curBattle.IsClose &&
		time.Now().After(time.Unix(curBattle.Begin, 0)) {
		parcelResult := game.RaidEliminateClose(s)
		rsp.RaidGiveUpDB = game.GetRaidEliminateGiveUpDB(s)
		rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResult)
	}
	game.RaidEliminateCheck(s)
	rsp.SeasonType = game.GetEliminateRaidSeasonType()
	rsp.RaidLobbyInfoDB = game.GetEliminateRaidLobbyInfoDB(s)
}

func EliminateRaidOpponentList(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.EliminateRaidOpponentListRequest)
	rsp := response.(*proto.EliminateRaidOpponentListResponse)

	rsp.OpponentUserDBs = make([]*proto.EliminateRaidUserDB, 0)
	cur := gdconf.GetCurRaidEliminateSchedule()
	if cur == nil {
		return
	}
	for i := int64(0); i < 15; i++ {
		ranking := req.Rank + i
		uid, _ := rank.GetUidByEliminateRank(cur.SeasonId, ranking)
		if uid == 0 {
			break
		}
		as := enter.GetSessionByUid(uid)
		if as != nil {
			rsp.OpponentUserDBs = append(rsp.OpponentUserDBs, game.GetEliminateRaidUserDB(as))
		}
	}
}

func EliminateRaidGetBestTeam(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.EliminateRaidGetBestTeamRequest)
	rsp := response.(*proto.EliminateRaidGetBestTeamResponse)

	rsp.RaidTeamSettingDBsDict = make(map[string][]*proto.RaidTeamSettingDB)
	as := enter.GetSessionByUid(req.SearchAccountId)
	if as == nil {
		return
	}
	for str, bin := range game.GetCurRaidEliminateInfo(as).GetBestScoreList() {
		if rsp.RaidTeamSettingDBsDict[str] == nil {
			rsp.RaidTeamSettingDBsDict[str] = make([]*proto.RaidTeamSettingDB, 0)
		}
		for _, teamBin := range bin.RaidTeamList {
			rsp.RaidTeamSettingDBsDict[str] = append(rsp.RaidTeamSettingDBsDict[str],
				game.GetRaidTeamSettingDB(as, teamBin))
		}
	}
}

func EliminateRaidCreateBattle(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.EliminateRaidCreateBattleRequest)
	rsp := response.(*proto.EliminateRaidCreateBattleResponse)

	defer func() {
		rsp.AccountCurrencyDB = game.GetAccountCurrencyDB(s)
	}()

	if game.GetEliminateRaidSeasonType() != proto.RaidSeasonType_Open {
		// 没开就请求,nt了
		return
	}
	game.NewCurRaidEliminateBattleInfo(s, req.RaidUniqueId, req.IsPractice)

	curBattle := game.GetCurRaidEliminateBattleInfo(s)
	if curBattle == nil {
		logger.Debug("大决战实例创建失败")
		return
	}
	if assist := req.AssistUseInfo; assist != nil && !curBattle.IsAssist {
		ac := enter.GetSessionByUid(assist.CharacterAccountId)
		assistInfo := game.GetAssistInfoByClanAssistUseInfo(ac, assist)
		rsp.AssistCharacterDB = game.GetAssistCharacterDB(ac, assistInfo, assist.AssistRelation)
	}

	rsp.RaidBattleDB = game.GetRaidBattleDB(s, curBattle)
	rsp.RaidDB = game.GetRaidDB(s, curBattle)
}

func EliminateRaidEndBattle(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.EliminateRaidEndBattleRequest)
	rsp := response.(*proto.EliminateRaidEndBattleResponse)

	curBattle := game.GetCurRaidEliminateBattleInfo(s)
	summary := req.Summary
	if summary == nil || curBattle == nil ||
		summary.RaidSummary == nil {
		return
	}
	raidSummary := summary.RaidSummary
	// 参战角色保存
	if !game.CheckRaidCharacter(s, int64(req.EchelonId), summary, curBattle) {
		return
	}
	// 记录boss情况
	for _, raidBossResult := range raidSummary.RaidBossResults {
		curBattle.AiPhase = raidBossResult.AIPhase
		curBattle.BossGroggyPoint += raidBossResult.RaidDamage.GivenGroggyPoint
		curBattle.GivenDamage += raidBossResult.RaidDamage.GivenDamage
		curBattle.IndexDamage = raidBossResult.RaidDamage.Index
	}
	curBattle.Frame += summary.EndFrame
	curBattle.ServerId++
	curBattle.IsClose = curBattle.MaxHp-curBattle.GivenDamage == 0
	// 判断是否结算
	if curBattle.IsClose {
		// 结算
		conf := gdconf.GetEliminateRaidStageExcelTable(curBattle.RaidUniqueId)
		if conf == nil {
			return
		}
		parcelResult := game.RaidEliminateClose(s)
		rsp.ClearTimePoint = curBattle.ClearTimePoint
		rsp.HPPercentScorePoint = curBattle.HpScorePoint
		rsp.DefaultClearPoint = curBattle.DefaultPoint
		rsp.RankingPoint = curBattle.ClearTimePoint + curBattle.HpScorePoint + curBattle.DefaultPoint
		rsp.BestRankingPoint = game.GetCurRaidEliminateInfo(s).GetBestScoreList()[conf.RaidBossGroup].GetScore()
		rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResult)
	}
}

func EliminateRaidEnterBattle(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.EliminateRaidEnterBattleRequest)
	rsp := response.(*proto.EliminateRaidEnterBattleResponse)

	curBattle := game.GetCurRaidEliminateBattleInfo(s)
	if curBattle == nil || // 没有战斗
		curBattle.RaidUniqueId != req.RaidUniqueId || // 实例不对
		game.GetEliminateRaidSeasonType() != proto.RaidSeasonType_Open || // 没开启
		time.Now().After(time.Unix(curBattle.Begin, 0).Add(1*time.Hour)) { // 超时了
		return
	}

	defer func() {
		rsp.AccountCurrencyDB = game.GetAccountCurrencyDB(s)
	}()

	if assist := req.AssistUseInfo; assist != nil && !curBattle.IsAssist {
		ac := enter.GetSessionByUid(assist.CharacterAccountId)
		assistInfo := game.GetAssistInfo(ac, assist.EchelonType, assist.SlotNumber)
		rsp.AssistCharacterDB = game.GetAssistCharacterDB(ac, assistInfo, assist.AssistRelation)
	}

	rsp.RaidBattleDB = game.GetRaidBattleDB(s, curBattle)
	rsp.RaidDB = game.GetRaidDB(s, curBattle)
}

func EliminateRaidGiveUp(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.EliminateRaidGiveUpRequest)
	rsp := response.(*proto.EliminateRaidGiveUpResponse)

	curBattle := game.GetCurRaidEliminateBattleInfo(s)
	if curBattle == nil ||
		req.IsPractice != curBattle.IsPractice {
		return
	}
	curBattle.IsClose = true
	parcelResult := game.RaidEliminateClose(s)
	if !curBattle.IsPractice {
		rsp.RaidGiveUpDB = game.GetRaidEliminateGiveUpDB(s)
		rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResult)
	}
}

func EliminateRaidSeasonReward(s *enter.Session, request, response proto.Message) {
	rsp := response.(*proto.EliminateRaidSeasonRewardResponse)

	bin := game.GetCurRaidEliminateInfo(s)
	if bin == nil {
		return
	}
	defer func() {
		game.SetServerNotification(s, proto.ServerNotificationFlag_CanReceiveEliminateRaidReward, false)
		rsp.ReceiveRewardIds = game.GetReceiveRewardIds(bin.GetReceiveRewardIds())
	}()
	if bin.ReceiveRewardIds == nil {
		bin.ReceiveRewardIds = make(map[int64]bool)
	}
	seasonConf := gdconf.GetEliminateRaidSeasonManageExcelTable(bin.SeasonId)
	if seasonConf == nil ||
		len(seasonConf.StackedSeasonRewardGauge) != len(seasonConf.SeasonRewardId) {
		return
	}
	parcelResultList := make([]*game.ParcelResult, 0)
	for index, season := range seasonConf.StackedSeasonRewardGauge {
		rewardId := seasonConf.SeasonRewardId[index]
		if _, ok := bin.ReceiveRewardIds[rewardId]; !ok &&
			bin.TotalScore >= season {
			rewardConf := gdconf.GetEliminateRaidStageSeasonRewardExcelTable(rewardId)
			if rewardConf == nil {
				continue
			}
			parcelResultList = append(parcelResultList,
				game.GetParcelResultList(rewardConf.SeasonRewardParcelType,
					rewardConf.SeasonRewardParcelUniqueId, rewardConf.SeasonRewardAmount, false)...)

			bin.ReceiveRewardIds[rewardId] = true
		}
	}
	rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResultList)
}

func EliminateRaidRankingReward(s *enter.Session, request, response proto.Message) {
	rsp := response.(*proto.EliminateRaidRankingRewardResponse)

	bin := game.GetCurRaidEliminateInfo(s)
	if bin == nil || bin.IsRankingReward {
		return
	}
	conf := gdconf.GetEliminateRaidRankingRewardExcelTableBySeasonId(bin.SeasonId, bin.Ranking)
	if conf == nil {
		return
	}
	rsp.ParcelResultDB = game.ParcelResultDB(s, game.GetParcelResultList(conf.RewardParcelType,
		conf.RewardParcelUniqueId, conf.RewardParcelAmount, false))
	bin.IsRankingReward = true
	bin.RankingRewardId = conf.Id
	rsp.ReceivedRankingRewardId = bin.RankingRewardId
}
