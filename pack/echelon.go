package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
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

func EchelonSave(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.EchelonSaveRequest)
	rsp := response.(*proto.EchelonSaveResponse)

	if req.EchelonDB != nil {
		bin := game.GetEchelonTypeInfo(s, int32(req.EchelonDB.EchelonType))
		if bin == nil {
			return
		}
		conf := &sro.DefaultEchelonExcelTable{
			EchlonId:      int32(req.EchelonDB.EchelonType),
			LeaderId:      game.ServerIdToCharacterId(s, req.EchelonDB.LeaderServerId),
			MainId:        game.ServerIdsToCharacterIds(s, req.EchelonDB.MainSlotServerIds),
			SupportId:     game.ServerIdsToCharacterIds(s, req.EchelonDB.SupportSlotServerIds),
			TssId:         game.ServerIdToCharacterId(s, req.EchelonDB.TSSInteractionServerId),
			SkillId:       game.ServerIdsToCharacterIds(s, req.EchelonDB.SkillCardMulliganCharacterIds),
			ExtensionType: int32(req.EchelonDB.ExtensionType),
		}

		info := game.UpEchelonInfo(bin, conf, req.EchelonDB.EchelonNumber)

		rsp.EchelonDB = game.GetEchelonDB(s, info)
	}

}
