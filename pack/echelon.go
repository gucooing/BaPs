package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/proto"
)

func EchelonList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.EchelonListResponse)

	rsp.EchelonDBs = make([]*proto.EchelonDB, 0)

	for _, dbType := range game.GetEchelonTypeInfoList(s) {
		if dbType == nil {
			continue
		}
		for _, db := range dbType.EchelonInfoList {
			rsp.EchelonDBs = append(rsp.EchelonDBs, game.GetEchelonDB(s, db))
		}
	}
}
