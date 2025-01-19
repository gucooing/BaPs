package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/pkg/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func CharacterList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.CharacterListResponse)

	rsp.TSSCharacterDBs = make([]*proto.CharacterDB, 0)
	rsp.CostumeDBs = make([]*proto.CostumeDB, 0)
	rsp.WeaponDBs = game.GetWeaponDBs(s)
	rsp.CharacterDBs = game.GetCharacterDBs(s)
}

func CharacterGearList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.CharacterGearListResponse)

	rsp.GearDBs = make([]*proto.GearDB, 0)
}

func CharacterTranscendence(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.CharacterTranscendenceRequest)
	rsp := response.(*proto.CharacterTranscendenceResponse)

	characterInfo := game.GetCharacterInfoByServerId(s, req.TargetCharacterServerId)
	if characterInfo == nil {
		return
	}
	defer func() {
		rsp.CharacterDB = game.GetCharacterDB(s, characterInfo.CharacterId)
	}()
	num := game.GetCharacterUpStarGradeNum(characterInfo.StarGrade + 1)
	itemInfo := game.GetItemInfo(s, characterInfo.CharacterId)
	if itemInfo == nil || num == 0 {
		return
	}
	if itemInfo.StackCount < num {
		return
	}
	rsp.ParcelResultDB = game.ParcelResultDB(s, []*game.ParcelResult{
		{
			ParcelType: proto.ParcelType_Item,
			ParcelId:   characterInfo.CharacterId,
			Amount:     int64(-num),
		},
	})
	characterInfo.StarGrade++
}

func CharacterUnlockWeapon(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.CharacterUnlockWeaponRequest)
	rsp := response.(*proto.CharacterUnlockWeaponResponse)

	characterInfo := game.GetCharacterInfoByServerId(s, req.TargetCharacterServerId)
	if characterInfo == nil {
		return
	}
	game.AddWeapon(s, characterInfo.CharacterId)
	rsp.WeaponDB = game.GetWeaponDB(s, characterInfo.CharacterId)
}

func CharacterWeaponTranscendence(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.CharacterWeaponTranscendenceRequest)
	rsp := response.(*proto.CharacterWeaponTranscendenceResponse)

	characterInfo := game.GetCharacterInfoByServerId(s, req.TargetCharacterServerId)
	if characterInfo == nil {
		return
	}
	waeponInfo := game.GetWeaponInfo(s, characterInfo.CharacterId)
	if waeponInfo == nil || waeponInfo.StarGrade >= 3 {
		return
	}
	num := game.GetWeaponUpStarGradeNum(waeponInfo.StarGrade)
	itemInfo := game.GetItemInfo(s, characterInfo.CharacterId)
	if itemInfo == nil || num == 0 {
		return
	}
	if itemInfo.StackCount < num {
		return
	}
	waeponInfo.StarGrade++
	rsp.ParcelResultDB = game.ParcelResultDB(s, []*game.ParcelResult{
		{
			ParcelType: proto.ParcelType_Item,
			ParcelId:   characterInfo.CharacterId,
			Amount:     int64(-num),
		},
		{
			ParcelType: proto.ParcelType_CharacterWeapon,
			ParcelId:   characterInfo.CharacterId,
			Amount:     0,
		},
	})
}

func CharacterSetFavorites(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.CharacterSetFavoritesRequest)
	rsp := response.(*proto.CharacterSetFavoritesResponse)

	rsp.CharacterDBs = make([]*proto.CharacterDB, 0)
	for sid, ok := range req.ActivateByServerIds {
		characterInfo := game.GetCharacterInfoByServerId(s, sid)
		if characterInfo == nil {
			continue
		}
		characterInfo.IsFavorite = ok
		rsp.CharacterDBs = append(rsp.CharacterDBs, game.GetCharacterDB(s, characterInfo.CharacterId))
	}
}
