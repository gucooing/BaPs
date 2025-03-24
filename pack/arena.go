package pack

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/common/rank"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func ArenaLogin(s *enter.Session, request, response proto.Message) {
	rsp := response.(*proto.ArenaLoginResponse)

	rsp.ArenaPlayerInfoDB = game.GetArenaPlayerInfoDB(s)
}

func ArenaEnterLobby(s *enter.Session, request, response proto.Message) {
	rsp := response.(*proto.ArenaEnterLobbyResponse)

	rsp.OpponentUserDBs = game.GetOpponentUserDBs(s)
	rsp.MapId = 1006
	rsp.AutoRefreshTime = mx.MxTime(game.GetArenAutoRefreshTime(s))
	rsp.ArenaPlayerInfoDB = game.GetArenaPlayerInfoDB(s)
}

func ArenaOpponentList(s *enter.Session, request, response proto.Message) {
	rsp := response.(*proto.ArenaOpponentListResponse)

	bin := game.GetArenaBin(s)
	if bin == nil {
		rsp.ErrorCode = 0
		return
	}
	rsp.PlayerRank = rank.GetArenaRank(bin.GetCurSeasonId(), s.AccountServerId)
	rsp.AutoRefreshTime = mx.MxTime(game.GetArenAutoRefreshTime(s))
	rsp.OpponentUserDBs = game.GetOpponentUserDBs(s)
}

func ArenaSyncEchelonSettingTime(s *enter.Session, request, response proto.Message) {
	// req := request.(*proto.ArenaSyncEchelonSettingTimeRequest)
	rsp := response.(*proto.ArenaSyncEchelonSettingTimeResponse)

	rsp.EchelonSettingTime = mx.MxTime(time.Now())
}

func ArenaEnterBattlePart1(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.ArenaEnterBattlePart1Request)
	rsp := response.(*proto.ArenaEnterBattlePart1Response)

	au := s.GetArenaUserByIndex(req.OpponentIndex)
	if au == nil {
		rsp.ErrorCode = proto.WebAPIErrorCode_ArenaInfoNotFound
		return
	}

	// 上锁！！！！！！！！！！！！！！！！！
	if !enter.AddArenaBattleRank(au.Rank) {
		rsp.ErrorCode = proto.WebAPIErrorCode_ArenaInfoNotFound
		return
	}
	if !enter.AddArenaBattleRank(rank.GetArenaRank(game.GetArenaBin(s).GetCurSeasonId(), s.AccountServerId)) {
		enter.DelArenaBattleRank(au.Rank)
		rsp.ErrorCode = proto.WebAPIErrorCode_ArenaInfoNotFound
		return
	}
	go enter.CheckArenaBattle(
		time.NewTicker(3+1*time.Hour),
		rank.GetArenaRank(game.GetArenaBin(s).GetCurSeasonId(), s.AccountServerId),
		au.Rank,
	)

	rsp.ArenaBattleDB = &proto.ArenaBattleDB{
		ArenaBattleServerId: s.AccountServerId,
		Season:              game.GetArenaBin(s).GetCurSeasonId(),
		Group:               game.GetArenaBin(s).GetPlayerGroupId(),
		BattleStartTime:     mx.MxTime(time.Now()),
		BattleEndTime:       mx.MxTime(time.Now().Add(3 * time.Hour)),
		Seed:                114514,
		AttackingUserDB:     game.GetPlayerArenaUserDB(s, proto.EchelonType_ArenaAttack),
		DefendingUserDB:     nil,
		BattleSummary:       nil,
	}

	if ps := enter.GetSessionByUid(au.Uid); ps != nil && !au.IsNpc {
		rsp.ArenaBattleDB.DefendingUserDB = game.GetPlayerArenaUserDB(ps, proto.EchelonType_ArenaDefence)
	} else {
		rsp.ArenaBattleDB.DefendingUserDB = game.GetNPCArenaUserDB(au)
	}
}

func ArenaEnterBattlePart2(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.ArenaEnterBattlePart2Request)
	rsp := response.(*proto.ArenaEnterBattlePart2Response)

	// 判断输赢

	// 交换排名

	// 发送奖励

	// 主动解除锁定
	rsp.ArenaBattleDB = &proto.ArenaBattleDB{
		ArenaBattleServerId: req.ArenaBattleDB.ArenaBattleServerId,
		Season:              req.ArenaBattleDB.Season,
		Group:               req.ArenaBattleDB.Group,
		BattleStartTime:     req.ArenaBattleDB.BattleStartTime,
		BattleEndTime:       req.ArenaBattleDB.BattleEndTime,
		Seed:                req.ArenaBattleDB.Seed,
		AttackingUserDB:     req.ArenaBattleDB.AttackingUserDB,
		DefendingUserDB:     req.ArenaBattleDB.DefendingUserDB,
		BattleSummary:       nil,
	}
	rsp.ArenaPlayerInfoDB = game.GetArenaPlayerInfoDB(s)
	rsp.AccountCurrencyDB = game.GetAccountCurrencyDB(s)
}
