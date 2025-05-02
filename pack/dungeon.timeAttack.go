package pack

import (
	"github.com/gucooing/BaPs/protocol/mx"
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/protocol/proto"
)

func TimeAttackDungeonLogin(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.TimeAttackDungeonLoginResponse)

	bin := game.GetTimeAttackBin(s)
	rsp.PreviousRoomDB = game.GetTimeAttackDungeonRoomDB(s, bin.GetPreviousRoom())
}

func TimeAttackDungeonLobby(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.TimeAttackDungeonLobbyResponse)

	defer func() {
		rsp.RoomDBs = game.GetTimeAttackDungeonRoomDBs(s)
	}()

	bin := game.GetTimeAttackBin(s)
	curBin := game.GetTimeAttackRoom(s, bin.GetCurRoom())
	if curBin != nil {
		// 三次战斗或超时进入战斗结算
		if len(curBin.BattleList) == 3 ||
			time.Now().After(time.Unix(curBin.StartTime, 0)) {
			parcelResult, _ := game.TimeAttackDungeonClose(s)
			rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResult)
		}
	}

	previous := game.GetTimeAttackRoom(s, bin.GetPreviousRoom())
	if previous != nil {
		rsp.PreviousRoomDB = game.GetTimeAttackDungeonRoomDB(s, bin.GetPreviousRoom())
		rsp.AchieveSeasonBestRecord = true    // 是否取得最佳战绩
		rsp.SeasonBestRecord = previous.Score // 最佳战绩得分
	}
}

func TimeAttackDungeonCreateBattle(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.TimeAttackDungeonCreateBattleRequest)
	rsp := response.(*proto.TimeAttackDungeonCreateBattleResponse)

	game.AddCurTimeAttackRoom(s, req.IsPractice)

	rsp.RoomDB = game.GetCurTimeAttackDungeonRoomDB(s)
	if !req.IsPractice { // 扣票
		game.UpCurrency(s, proto.CurrencyTypes_TimeAttackDungeonTicket, -1)
	}
	rsp.ParcelResultDB = game.ParcelResultDB(s, nil)
}

func TimeAttackDungeonEnterBattle(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.TimeAttackDungeonEnterBattleRequest)
	rsp := response.(*proto.TimeAttackDungeonEnterBattleResponse)

	assist := req.AssistUseInfo
	if assist != nil {
		ac := enter.GetSessionByUid(assist.CharacterAccountId)
		assistInfo := game.GetAssistInfoByClanAssistUseInfo(ac, assist)
		rsp.AssistCharacterDB = game.GetAssistCharacterDB(ac, assistInfo, assist.AssistRelation)
	}
}

func TimeAttackDungeonEndBattle(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.TimeAttackDungeonEndBattleRequest)
	rsp := response.(*proto.TimeAttackDungeonEndBattleResponse)

	defer func() {
		rsp.RoomDB = game.GetCurTimeAttackDungeonRoomDB(s)
	}()

	bin := game.GetTimeAttackBin(s)
	curBin := game.GetTimeAttackRoom(s, bin.GetCurRoom())
	summary := req.Summary
	if summary == nil || curBin == nil || summary.Winner != proto.GroupTag_Group01 {
		return
	}
	conf := gdconf.GetTimeAttackDungeonGeasExcelTable(summary.StageId)
	if conf == nil {
		return
	}
	// 分数计算
	clearTimePoint := conf.ClearTimeWeightPoint
	battleInfo := &sro.TimeAttackRoomBattleHistory{
		GeasId:              conf.Id,
		Frame:               int64(summary.EndFrame),
		ClearTimePoint:      clearTimePoint,
		DefaultPoint:        conf.ClearDefaultPoint,
		DungeonType:         int32(proto.TimeAttackDungeonType_None.Value(conf.TimeAttackDungeonType)),
		MainCharacterDBs:    make([]*sro.TimeAttackDungeonCharacter, 0),
		SupportCharacterDBs: make([]*sro.TimeAttackDungeonCharacter, 0),
	}
	if curBin.BattleList == nil {
		curBin.BattleList = make([]*sro.TimeAttackRoomBattleHistory, 0)
	}
	curBin.BattleList = append(curBin.BattleList, battleInfo)
	// 参战角色保存
	getTimeAttackDungeonCharacter := func(hs *proto.HeroSummary) *sro.TimeAttackDungeonCharacter {
		info := &sro.TimeAttackDungeonCharacter{
			ServerId:  hs.ServerId,
			UniqueId:  hs.CharacterId,
			CostumeId: hs.CostumeId,
			StarGrade: hs.Grade,
			Level:     hs.Level,
			IsAssist:  s.AccountServerId != hs.OwnerAccountId,
		}
		if hs.CharacterWeapon != nil {
			info.HasWeapon = true
			info.WeaponInfo = &sro.WeaponInfo{
				UniqueId:          hs.CharacterWeapon.UniqueId,
				CharacterServerId: hs.CharacterId,
				StarGrade:         hs.CharacterWeapon.StarGrade,
				Level:             hs.CharacterWeapon.Level,
			}
		}
		return info
	}
	for _, heroe := range summary.Group01Summary.Heroes {
		battleInfo.MainCharacterDBs = append(battleInfo.MainCharacterDBs, getTimeAttackDungeonCharacter(heroe))
	}
	for _, heroe := range summary.Group01Summary.Supporters {
		battleInfo.SupportCharacterDBs = append(battleInfo.SupportCharacterDBs, getTimeAttackDungeonCharacter(heroe))
	}

	rsp.TotalPoint = conf.ClearDefaultPoint + clearTimePoint
	rsp.DefaultPoint = conf.ClearDefaultPoint
	rsp.TimePoint = clearTimePoint
}

func TimeAttackDungeonGiveUp(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.TimeAttackDungeonGiveUpRequest)
	rsp := response.(*proto.TimeAttackDungeonGiveUpResponse)

	bin := game.GetTimeAttackBin(s)
	curBin := game.GetTimeAttackRoom(s, bin.GetCurRoom())
	if curBin == nil || curBin.RoomId != req.RoomId {
		return
	}
	defer func() {
		rsp.RoomDB = game.GetTimeAttackDungeonRoomDB(s, req.RoomId)
	}()
	// 结算
	parcelResult, score := game.TimeAttackDungeonClose(s)
	rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResult)
	rsp.SeasonBestRecord = score
}

func TimeAttackDungeonSweep(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.TimeAttackDungeonSweepRequest)
	rsp := response.(*proto.TimeAttackDungeonSweepResponse)

	rsp.Rewards = make([][]*proto.ParcelInfo, 0)

	bin := game.GetTimeAttackBin(s)
	previous := game.GetTimeAttackRoom(s, bin.GetPreviousRoom())
	if previous == nil {
		return
	}
	// 扣票
	game.UpCurrency(s, proto.CurrencyTypes_TimeAttackDungeonTicket, -req.SweepCount)
	if previous.SweepTime == nil {
		previous.SweepTime = make([]int64, 0)
	}
	for i := int64(0); i < req.SweepCount; i++ {
		previous.SweepTime = append(previous.SweepTime, time.Now().Unix())
		// 奖励
		list := game.GetTimeAttackDungeonParcelResultByScore(previous.Score, previous.SeasonId)
		parcelResultDB := game.ParcelResultDB(s, list)
		rsp.Rewards = append(rsp.Rewards, parcelResultDB.ParcelForMission)
	}
	rsp.RoomDB = game.GetTimeAttackDungeonRoomDB(s, previous.RoomId)
	rsp.ParcelResultDB = game.ParcelResultDB(s, nil)
}
