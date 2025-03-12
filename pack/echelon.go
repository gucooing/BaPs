package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/protocol/proto"
)

func EchelonList(s *enter.Session, request, response proto.Message) {
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

func EchelonSave(s *enter.Session, request, response proto.Message) {
	req := request.(*proto.EchelonSaveRequest)
	rsp := response.(*proto.EchelonSaveResponse)

	if req.EchelonDB != nil {
		bin := game.GetEchelonTypeInfo(s, int32(req.EchelonDB.EchelonType))
		if bin == nil {
			return
		}
		conf := &sro.DefaultEchelonExcelTable{
			EchlonId:      int32(req.EchelonDB.EchelonType),
			LeaderId:      s.GetCharacterByKeyId(req.EchelonDB.LeaderServerId).GetCharacterId(),
			MainId:        game.ServerIdsToCharacterIds(s, req.EchelonDB.MainSlotServerIds),
			SupportId:     game.ServerIdsToCharacterIds(s, req.EchelonDB.SupportSlotServerIds),
			TssId:         s.GetCharacterByKeyId(req.EchelonDB.TSSInteractionServerId).GetCharacterId(),
			SkillId:       req.EchelonDB.SkillCardMulliganCharacterIds,
			ExtensionType: int32(req.EchelonDB.ExtensionType),
		}

		info := game.UpEchelonInfo(bin, conf, req.EchelonDB.EchelonNumber)

		rsp.EchelonDB = game.GetEchelonDB(s, info)
	}

}

func EchelonPresetList(s *enter.Session, request, response proto.Message) {
	rsp := response.(*proto.EchelonPresetListResponse)

	rsp.PresetGroupDBs = make([]*proto.EchelonPresetGroupDB, 0)

	for gid, typeInfo := range game.GetEchelonPresetGuidList(s) {
		presetGroupDB := &proto.EchelonPresetGroupDB{
			GroupIndex:    gid,
			ExtensionType: proto.EchelonExtensionType_Base,
			GroupLabel:    "",
			PresetDBs:     make(map[int32]*proto.EchelonPresetDB),
			Item:          nil,
		}
		for index, info := range typeInfo.EchelonInfoList {
			presetGroupDB.PresetDBs[int32(index)] = game.GetEchelonPresetGroupDB(info)
		}
		rsp.PresetGroupDBs = append(rsp.PresetGroupDBs, presetGroupDB)
	}
}
