package pack

import (
	"github.com/gucooing/BaPs/common/enter"
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
	rsp.MapId = 1004
	rsp.AutoRefreshTime = mx.MxTime{}
	rsp.ArenaPlayerInfoDB = game.GetArenaPlayerInfoDB(s)
}

func ArenaOpponentList(s *enter.Session, request, response proto.Message) {
	rsp := response.(*proto.ArenaOpponentListResponse)

	rsp.PlayerRank = 1
	rsp.AutoRefreshTime = mx.MxTime{}
	rsp.OpponentUserDBs = game.GetOpponentUserDBs(s)
}
