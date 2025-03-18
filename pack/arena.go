package pack

import (
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
