package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/proto"
)

func CharacterList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.CharacterListResponse)

	rsp.TSSCharacterDBs = make([]*proto.CharacterDB, 0)
	rsp.WeaponDBs = make([]*proto.WeaponDB, 0)
	rsp.CostumeDBs = make([]*proto.CostumeDB, 0)
	rsp.CharacterDBs = game.GetCharacterDBs(s)
}

func CharacterGearList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.CharacterGearListResponse)

	rsp.GearDBs = make([]*proto.GearDB, 0)
}
