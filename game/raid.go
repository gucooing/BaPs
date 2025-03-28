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

func GetRaidSeasonType() proto.RaidSeasonType {
	cur := gdconf.GetCurRaidSchedule()
	if cur == nil {
		return proto.RaidSeasonType_Close
	}
	if cur.StartTime.Before(time.Now()) &&
		cur.EndTime.After(time.Now()) {
		return proto.RaidSeasonType_Open
	}
	return proto.RaidSeasonType_Settlement
}

func GetRaidBin(s *enter.Session) *sro.RaidBin {
	bin := GetPlayerBin(s)
	if bin == nil {
		return nil
	}
	if bin.RaidBin == nil {
		bin.RaidBin = &sro.RaidBin{}
	}
	return bin.RaidBin
}

func GetCurRaidInfo(s *enter.Session) *sro.RaidInfo {
	bin := GetRaidBin(s)
	conf := gdconf.GetCurRaidSchedule()
	if bin == nil || conf == nil {
		return nil
	}
	if bin.CurRaidInfo == nil ||
		bin.GetCurRaidInfo().GetSeasonId() != conf.SeasonId {
		bin.CurRaidInfo = &sro.RaidInfo{
			SeasonId: conf.SeasonId,
		}
	}
	// 如果赛季已经进入结算期
	if time.Now().After(conf.EndTime) {
		bin.CurRaidInfo.Ranking = rank.GetRaidRank(conf.SeasonId, s.AccountServerId)
		bin.CurRaidInfo.Tier = gdconf.GetRaidTier(bin.CurRaidInfo.SeasonId,
			bin.CurRaidInfo.Ranking)
	}

	return bin.CurRaidInfo
}

// GetLastRaidInfo 此处返回的可能是nil
func GetLastRaidInfo(s *enter.Session) *sro.RaidInfo {
	GetCurRaidInfo(s)

	return GetRaidBin(s).GetLastRaidInfo()
}

func GetCurRaidTeamList(s *enter.Session) map[int32]*sro.RaidTeamInfo {
	bin := GetCurRaidInfo(s)
	if bin == nil {
		return nil
	}
	if bin.RaidTeamList == nil {
		bin.RaidTeamList = make(map[int32]*sro.RaidTeamInfo)
	}
	return bin.RaidTeamList
}

// NewCurRaidBattleInfo 创建新的总力战
func NewCurRaidBattleInfo(s *enter.Session, raidUniqueId int64, isPractice bool) {
	bin := GetRaidBin(s)
	conf := gdconf.GetRaidStageExcelTable(raidUniqueId)
	if bin == nil || conf == nil {
		logger.Debug("玩家实例不存在或总力战关卡不存在RaidUniqueId:%v", raidUniqueId)
		return
	}
	bin.CurRaidBattleInfo = &sro.CurRaidBattleInfo{
		RaidUniqueId: raidUniqueId,
		IsPractice:   isPractice,
		RaidTeamList: make(map[int32]*sro.RaidTeamInfo),
		Frame:        0,
		Begin:        time.Now().Unix(),
		// MaxHp:        ,
		SeasonId:     GetCurRaidInfo(s).SeasonId,
		ServerId:     1,
		ContentType:  proto.ContentType_Raid,
		RaidBoosList: make([]*sro.RaidBoosInfo, 0),
	}
	for index, bosscid := range conf.BossCharacterId {
		chConf := gdconf.GetCharacterStatExcelTable(bosscid)
		if chConf == nil {
			logger.Error("总力战boss实例不存在BossCharacterId:%v", bosscid)
			return
		}
		raidBoosInfo := &sro.RaidBoosInfo{
			Index: int32(index),
			MaxHp: chConf.MaxHP100,
		}
		bin.CurRaidBattleInfo.RaidBoosList = append(bin.CurRaidBattleInfo.RaidBoosList, raidBoosInfo)
	}
}

func GetCurRaidBattleInfo(s *enter.Session) *sro.CurRaidBattleInfo {
	return GetRaidBin(s).GetCurRaidBattleInfo()
}

func RaidCheck(s *enter.Session) {
	bin := GetCurRaidInfo(s)
	if bin == nil {
		return
	}
	// 检查总分奖励领取
	seasonConf := gdconf.GetRaidSeasonManageExcelTable(bin.SeasonId)
	if seasonConf == nil ||
		len(seasonConf.StackedSeasonRewardGauge) != len(seasonConf.SeasonRewardId) {
		return
	}
	for index, season := range seasonConf.StackedSeasonRewardGauge {
		if _, ok := bin.ReceiveRewardIds[seasonConf.SeasonRewardId[index]]; !ok &&
			bin.TotalScore >= season {
			SetServerNotification(s, proto.ServerNotificationFlag_CanReceiveRaidReward, true)
			break
		}
	}
}

func GetClearDifficulty(s *enter.Session) []proto.Difficulty {
	list := []proto.Difficulty{
		proto.Difficulty_Normal,
		proto.Difficulty_Hard,
		proto.Difficulty_VeryHard,
		proto.Difficulty_Hardcore,
		proto.Difficulty_Extreme,
		proto.Difficulty_Insane,
		proto.Difficulty_Torment,
		proto.Difficulty_Lunatic,
	}

	return list
}

func GetPlayableHighestDifficulty(s *enter.Session) map[string]proto.Difficulty {
	list := make(map[string]proto.Difficulty, 0)
	if cur := gdconf.GetCurRaidSchedule(); cur != nil {
		conf := gdconf.GetRaidSeasonManageExcelTable(cur.SeasonId)
		for _, name := range conf.OpenRaidBossGroup {
			list[name] = proto.Difficulty(alg.MinInt32(GetCurRaidInfo(s).GetDifficulty()+1, proto.Difficulty_Lunatic))
		}
	}
	return list
}

func GetReceiveRewardIds(bin map[int64]bool) []int64 {
	list := make([]int64, 0)
	for id, ok := range bin {
		if ok {
			list = append(list, id)
		}
	}
	return list
}

func GetCanReceiveRankingReward(isTime, isReward bool) bool {
	return isTime && isReward == false
}

func GetRaidLobbyInfoDB(s *enter.Session) *proto.RaidLobbyInfoDB {
	bin := GetCurRaidInfo(s)
	info := &proto.RaidLobbyInfoDB{
		PlayableHighestDifficulty: GetPlayableHighestDifficulty(s),
		PlayingRaidDB:             GetRaidDB(s, GetCurRaidBattleInfo(s)),
		ReceiveRewardIds:          GetReceiveRewardIds(bin.GetReceiveRewardIds()),
		ReceivedRankingRewardId:   bin.GetRankingRewardId(),

		RemainFailCompensation: map[int32]bool{
			0: false,
		},
		SweepPointByRaidUniqueId: make(map[int64]int64),
		ReceiveLimitedRewardIds:  make([]int64, 0),
		ClanAssistUseInfo:        nil,
	}
	if cur := gdconf.GetCurRaidSchedule(); cur != nil {
		info.SeasonId = cur.SeasonId
		info.SeasonStartDate = mx.MxTime(cur.StartTime)
		info.SeasonEndDate = mx.MxTime(cur.EndTime)
		info.SettlementEndDate = mx.MxTime(cur.EndTime)
		info.Ranking = rank.GetRaidRank(cur.SeasonId, s.AccountServerId)
		info.Tier = gdconf.GetRaidTier(cur.SeasonId, info.Ranking)
		info.BestRankingPoint = bin.GetBestScore()   // 最高分
		info.TotalRankingPoint = bin.GetTotalScore() // 总分
		info.CanReceiveRankingReward = GetCanReceiveRankingReward(
			time.Now().After(cur.EndTime), bin.GetIsRankingReward())
	}
	if next := gdconf.GetNextRaidSchedule(); next != nil {
		info.NextSeasonId = next.SeasonId
		info.NextSeasonStartDate = mx.MxTime(next.StartTime)
		info.NextSeasonEndDate = mx.MxTime(next.EndTime)
		info.SettlementEndDate = mx.MxTime(next.EndTime)
	}
	// 如果有进行中的战斗
	if curBattle := GetCurRaidBattleInfo(s); curBattle != nil &&
		!curBattle.IsClose {
		for _, teamInfo := range GetCurRaidBattleInfo(s).GetRaidTeamList() {
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

// GetSingleRaidUserDB 拉取的一定是本次数据
func GetSingleRaidUserDB(s *enter.Session) *proto.SingleRaidUserDB {
	info := &proto.SingleRaidUserDB{
		RaidUserDB:        GetRaidUserDB(s),
		RaidTeamSettingDB: GetRaidTeamSettingDB(s, GetCurRaidTeamList(s)[1]),
	}

	return info
}

func GetRaidUserDB(s *enter.Session) *proto.RaidUserDB {
	curInfo := GetCurRaidInfo(s)
	if curInfo == nil {
		return nil
	}
	ranking, score := rank.GetRaidRankAndScore(curInfo.SeasonId, s.AccountServerId)
	info := &proto.RaidUserDB{
		AccountId:                   s.AccountServerId,
		RepresentCharacterUniqueId:  GetRepresentCharacterUniqueId(s),
		RepresentCharacterCostumeId: 0,
		Level:                       int64(GetAccountLevel(s)),
		Nickname:                    GetNickname(s),
		Tier:                        gdconf.GetRaidTier(curInfo.SeasonId, ranking),
		Rank:                        ranking,
		BestRankingPoint:            int64(score),
		BestRankingPointDetail:      score,
		AccountAttachmentDB:         GetAccountAttachmentDB(s),
	}

	return info
}

func GetRaidTeamSettingDB(s *enter.Session, teamInfo *sro.RaidTeamInfo) *proto.RaidTeamSettingDB {
	if teamInfo == nil {
		return new(proto.RaidTeamSettingDB)
	}
	info := &proto.RaidTeamSettingDB{
		AccountId:                     s.AccountServerId,
		EchelonType:                   proto.EchelonType(teamInfo.EchelonType),
		EchelonExtensionType:          proto.EchelonExtensionType_Base,
		SkillCardMulliganCharacterIds: make([]int64, 0),
		LeaderCharacterUniqueId:       teamInfo.LeaderCharacter,
		MainCharacterDBs:              make([]*proto.RaidCharacterDB, 0),
		SupportCharacterDBs:           make([]*proto.RaidCharacterDB, 0),
		TryNumber:                     teamInfo.TryNumber,
		TSSInteractionUniqueId:        0,
	}
	for _, cid := range teamInfo.SkillCharacterList {
		info.SkillCardMulliganCharacterIds =
			append(info.SkillCardMulliganCharacterIds, cid)
	}
	for slot, bin := range teamInfo.MainCharacterList {
		info.MainCharacterDBs = append(info.MainCharacterDBs, GetRaidCharacterDB(bin, slot))
	}
	for slot, bin := range teamInfo.SupportCharacterList {
		info.SupportCharacterDBs = append(info.SupportCharacterDBs, GetRaidCharacterDB(bin, slot))
	}

	return info
}

func GetRaidCharacterDB(bin *sro.RaidCharacterInfo, slot int32) *proto.RaidCharacterDB {
	if bin == nil {
		return nil
	}
	info := &proto.RaidCharacterDB{
		Level:            bin.Level,
		HasWeapon:        bin.HasWeapon,
		ServerId:         bin.ServerId,
		UniqueId:         bin.CharacterId,
		StarGrade:        bin.StarGrade,
		SlotIndex:        slot,
		AccountId:        bin.AccountId,
		IsAssist:         bin.IsAssist,
		WeaponStarGrade:  bin.WeaponStarGrade,
		CostumeId:        0,
		CombatStyleIndex: 0,
	}

	return info
}

func GetRaidBattleDB(s *enter.Session, bin *sro.CurRaidBattleInfo) *proto.RaidBattleDB {
	if bin == nil || bin.IsClose {
		return nil
	}
	info := &proto.RaidBattleDB{
		ContentType:        proto.ContentType(bin.ContentType),
		CurrentBossHP:      0,
		RaidMembers:        make([]*proto.RaidMemberDescription, 0),
		RaidUniqueId:       bin.RaidUniqueId,
		CurrentBossAIPhase: bin.AiPhase,
		CurrentBossGroggy:  0,
		IsClear:            bin.IsClose,
		RaidBossIndex:      0,

		BIEchelon:   "",
		SubPartsHPs: make([]int64, 0),
	}

	rmd := &proto.RaidMemberDescription{
		AccountId:        s.AccountServerId,
		AccountName:      GetNickname(s),
		CharacterId:      0,
		DamageCollection: make([]*proto.RaidDamageCollection, 0),
	}

	info.RaidMembers = append(info.RaidMembers, rmd)

	isIndex := false
	for _, raidBoosInfo := range bin.RaidBoosList {
		rdc := &proto.RaidDamageCollection{
			Index:            raidBoosInfo.Index,
			GivenGroggyPoint: raidBoosInfo.BossGroggyPoint,
			GivenDamage:      raidBoosInfo.GivenDamage,
		}
		rmd.DamageCollection = append(rmd.DamageCollection, rdc)

		curHp := raidBoosInfo.MaxHp - raidBoosInfo.GivenDamage

		if curHp != 0 && !isIndex {
			isIndex = true
			info.CurrentBossGroggy = raidBoosInfo.BossGroggyPoint
			info.RaidBossIndex = raidBoosInfo.Index
		}

		info.CurrentBossHP += curHp
	}

	return info
}

func GetRaidDB(s *enter.Session, bin *sro.CurRaidBattleInfo) *proto.RaidDB {
	if bin == nil || bin.IsClose {
		return nil
	}
	info := &proto.RaidDB{
		AccountLevelWhenCreateDB: int64(GetAccountLevel(s)),
		Begin:                    mx.Unix(bin.Begin, 0),
		ContentType:              proto.ContentType(bin.ContentType),
		End:                      mx.Unix(bin.Begin, 0), // .Add(1 * time.Hour),
		Owner: &proto.RaidMemberDescription{
			AccountId:   s.AccountServerId,
			AccountName: GetNickname(s),
			CharacterId: 0,
		},
		PlayerCount:                   1,
		RaidBossDBs:                   make([]*proto.RaidBossDB, 0),
		RaidState:                     proto.RaidStatus_Playing,
		SeasonId:                      bin.SeasonId,
		UniqueId:                      bin.RaidUniqueId,
		IsPractice:                    bin.IsPractice,
		ClanAssistUsed:                bin.IsAssist,
		ParticipateCharacterServerIds: make(map[int64][]int64),
		ServerId:                      bin.ServerId,
		LastBossIndex:                 0,

		SecretCode:           "0",
		OwnerAccountServerId: 0,
		OwnerNickname:        "",
		BossGroup:            "",
		BossDifficulty:       0,
		Tags:                 make([]int32, 0),
		IsEnterRoom:          false,
		SessionHitPoint:      0,
	}
	for index, teamInfo := range bin.RaidTeamList {
		list := make([]int64, 0)
		for _, rcInfo := range teamInfo.MainCharacterList {
			list = append(list, rcInfo.ServerId)
		}
		for _, rcInfo := range teamInfo.SupportCharacterList {
			list = append(list, rcInfo.ServerId)
		}
		info.ParticipateCharacterServerIds[int64(index)] = list
	}

	isIndex := false

	for _, raidBoosInfo := range bin.RaidBoosList {
		curHp := raidBoosInfo.MaxHp - raidBoosInfo.GivenDamage
		rb := &proto.RaidBossDB{
			ContentType:     proto.ContentType(bin.ContentType),
			BossIndex:       raidBoosInfo.Index,
			BossCurrentHP:   curHp,
			BossGroggyPoint: raidBoosInfo.BossGroggyPoint,
		}
		info.RaidBossDBs = append(info.RaidBossDBs, rb)

		if curHp != 0 && !isIndex {
			isIndex = true
			info.LastBossIndex = raidBoosInfo.Index
		}
	}

	return info
}

func GetRaidGiveUpDB(s *enter.Session) *proto.RaidGiveUpDB {
	curBattle := GetCurRaidBattleInfo(s)
	if curBattle == nil {
		return nil
	}
	ranking, bast := rank.GetRaidRankAndScore(curBattle.SeasonId, s.AccountServerId)
	info := &proto.RaidGiveUpDB{
		Ranking:          ranking,
		RankingPoint:     curBattle.ClearTimePoint + curBattle.HpScorePoint + curBattle.DefaultPoint,
		BestRankingPoint: int64(bast),
	}
	return info
}

// CheckRaidCharacter 参战角色验证
func CheckRaidCharacter(s *enter.Session, echelonId int64,
	summary *proto.BattleSummary, curBattle *sro.CurRaidBattleInfo) bool {
	if curBattle == nil {
		return false
	}
	echelonType := proto.EchelonType_Raid
	if curBattle.ContentType == proto.ContentType_EliminateRaid {
		conf := gdconf.GetEliminateRaidStageExcelTable(curBattle.RaidUniqueId)
		echelonType = gdconf.GetEliminateRaidEchelonType(curBattle.SeasonId, conf.GetRaidBossGroup())
	}
	if echelonType == proto.EchelonType_None {
		logger.Warn("未知的总力战boss类型")
		return false
	}
	echelonInfo := GetEchelonInfo(s, echelonType.Value(), echelonId)
	if curBattle.RaidTeamList == nil {
		curBattle.RaidTeamList = make(map[int32]*sro.RaidTeamInfo)
	}
	raidTeamInfo := &sro.RaidTeamInfo{
		LeaderCharacter:      echelonInfo.LeaderCharacter,
		MainCharacterList:    make(map[int32]*sro.RaidCharacterInfo),
		SupportCharacterList: make(map[int32]*sro.RaidCharacterInfo),
		SkillCharacterList:   echelonInfo.SkillCharacterList,
		TryNumber:            curBattle.ServerId,
		EchelonType:          echelonInfo.EchelonType,
	}

	if len(summary.Group01Summary.Heroes) > 6 {
		return false
	}
	if len(summary.Group01Summary.Supporters) > 4 {
		return false
	}
	getRaidCharacterInfo := func(hero *proto.HeroSummary) *sro.RaidCharacterInfo {
		raidCharacterInfo := &sro.RaidCharacterInfo{
			CharacterId:     hero.CharacterId,
			HasWeapon:       false,
			Level:           hero.Level,
			ServerId:        hero.ServerId,
			StarGrade:       hero.Grade,
			WeaponStarGrade: 0,
			IsAssist:        hero.OwnerAccountId != s.AccountServerId,
			AccountId:       hero.OwnerAccountId,
		}
		if hero.CharacterWeapon != nil {
			raidCharacterInfo.HasWeapon = true
			raidCharacterInfo.WeaponStarGrade = hero.CharacterWeapon.StarGrade
		}
		return raidCharacterInfo
	}
	// 添加主角色
	for index, heroe := range summary.Group01Summary.Heroes {
		if heroe.OwnerAccountId != s.AccountServerId {
			if curBattle.IsAssist {
				return false // 援助角色超限
			}

			curBattle.IsAssist = true
		}
		raidTeamInfo.MainCharacterList[int32(index)] = getRaidCharacterInfo(heroe)
	}
	// 添加副角色
	for index, heroe := range summary.Group01Summary.Supporters {
		if heroe.OwnerAccountId != s.AccountServerId {
			if curBattle.IsAssist {
				return false // 援助角色超限
			}
			curBattle.IsAssist = true
		}
		raidTeamInfo.SupportCharacterList[int32(index)] = getRaidCharacterInfo(heroe)
	}
	// 参战角色验证通过
	curBattle.RaidTeamList[int32(curBattle.ServerId)] = raidTeamInfo
	return true
}

func RaidClose(s *enter.Session) []*ParcelResult {
	curBattle := GetCurRaidBattleInfo(s)
	cur := GetCurRaidInfo(s)
	if curBattle == nil {
		return nil
	}
	conf := gdconf.GetRaidStageExcelTable(curBattle.RaidUniqueId)
	if conf == nil {
		return nil
	}
	list := make([]*ParcelResult, 0)

	givenDamage := int64(0)
	mxHp := int64(0)
	for _, raidBoosInfo := range curBattle.RaidBoosList {
		givenDamage += raidBoosInfo.GivenDamage
		mxHp += raidBoosInfo.MaxHp
	}

	// 计算分数
	curBattle.DefaultPoint = conf.DefaultClearScore
	curBattle.HpScorePoint = conf.HPPercentScore * givenDamage / mxHp
	curBattle.ClearTimePoint = alg.MaxInt64(conf.MaximumScore-conf.PerSecondMinusScore/300*int64(curBattle.Frame), 0)

	// 如果不是模拟,且战斗结束
	if !curBattle.IsPractice && curBattle.IsClose && len(curBattle.RaidTeamList) > 0 {
		if mxHp-givenDamage == 0 {
			// 更新通关难度
			cur.Difficulty = alg.MaxInt32(cur.Difficulty, int32(proto.GetDifficultyByStr(conf.Difficulty)))
		}
		rankingPoint := curBattle.ClearTimePoint + curBattle.HpScorePoint + curBattle.DefaultPoint
		cur.TotalScore += rankingPoint // 累积分数
		if cur.BestScore < rankingPoint {
			cur.BestScore = rankingPoint              // 标记最高分
			cur.RaidTeamList = curBattle.RaidTeamList // 标记最高分队伍
			rank.SetRaidScore(curBattle.SeasonId, s.AccountServerId, float64(rankingPoint))
		}
		// 计算奖励
		for _, rewardConf := range gdconf.GetRaidStageRewardExcelTable(conf.RaidRewardGroupId) {
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
