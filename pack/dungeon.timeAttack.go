package pack

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/protocol/proto"
)

func TimeAttackDungeonLogin(s *enter.Session, request, response proto.Message) {
	rsp := response.(*proto.TimeAttackDungeonLoginResponse)

	bin := game.GetTimeAttackBin(s)
	rsp.PreviousRoomDB = game.GetTimeAttackDungeonRoomDB(s, bin.GetPreviousRoom())
}

func TimeAttackDungeonLobby(s *enter.Session, request, response proto.Message) {
	rsp := response.(*proto.TimeAttackDungeonLobbyResponse)

	bin := game.GetTimeAttackBin(s)
	curBin := game.GetTimeAttackRoom(s, bin.GetCurRoom())
	if curBin != nil {
		// 三次战斗或超时进入战斗结算
		if len(curBin.BattleList) == 3 ||
			time.Now().After(time.Unix(curBin.StartTime, 0).Add(1*time.Hour)) {

			rsp.ParcelResultDB = game.ParcelResultDB(s, nil)
		}
	}

	previous := game.GetTimeAttackRoom(s, bin.GetPreviousRoom())
	if previous != nil {
		rsp.PreviousRoomDB = game.GetTimeAttackDungeonRoomDB(s, bin.GetPreviousRoom())
		rsp.RoomDBs = game.GetTimeAttackDungeonRoomDBs(s)
		rsp.AchieveSeasonBestRecord = false // 是否取得最佳战绩
		rsp.SeasonBestRecord = 0            // 最佳战绩得分
	}
}

func TimeAttackDungeonCreateBattle(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.TimeAttackDungeonCreateBattleRequest)
	rsp := response.(*proto.TimeAttackDungeonCreateBattleResponse)

	game.AddCurTimeAttackRoom(s, req.IsPractice)

	rsp.RoomDB = game.GetCurTimeAttackDungeonRoomDB(s)
	rsp.ParcelResultDB = game.ParcelResultDB(s, nil)
}

func TimeAttackDungeonEnterBattle(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.TimeAttackDungeonEnterBattleRequest)
	rsp := response.(*proto.TimeAttackDungeonEnterBattleResponse)

	assist := req.AssistUseInfo
	if assist != nil {
		ac := enter.GetSessionByUid(assist.CharacterAccountId)
		assistInfo := game.GetAssistInfoByClanAssistUseInfo(ac, assist)
		rsp.AssistCharacterDB = game.GetAssistCharacterDB(ac, assistInfo, assist.AssistRelation)
	}
}

func TimeAttackDungeonEndBattle(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.TimeAttackDungeonEndBattleRequest)
	rsp := response.(*proto.TimeAttackDungeonEndBattleResponse)

	defer func() {
		rsp.RoomDB = game.GetCurTimeAttackDungeonRoomDB(s)
	}()

	summary := req.Summary
	if summary == nil {
		return
	}
}
