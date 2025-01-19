package pack

import (
	"strings"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/gdconf"
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

func CharacterUpdateSkillLevel(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.CharacterSkillLevelUpdateRequest)
	rsp := response.(*proto.CharacterSkillLevelUpdateResponse)

	characterInfo := game.GetCharacterInfoByServerId(s, req.TargetCharacterDBId)
	if characterInfo == nil {
		return
	}
	defer func() {
		rsp.CharacterDB = game.GetCharacterDB(s, characterInfo.CharacterId)
	}()
	parcelResultList := UpCharacterSkill(characterInfo, req.Level, req.SkillSlot)
	rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResultList)
}

func CharacterBatchSkillLevelUpdate(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.CharacterBatchSkillLevelUpdateRequest)
	rsp := response.(*proto.CharacterBatchSkillLevelUpdateResponse)

	characterInfo := game.GetCharacterInfoByServerId(s, req.TargetCharacterDBId)
	if characterInfo == nil {
		return
	}
	defer func() {
		rsp.CharacterDB = game.GetCharacterDB(s, characterInfo.CharacterId)
	}()
	var parcelResultList []*game.ParcelResult
	for _, skillUp := range req.SkillLevelUpdateRequestDBs {
		parcelResultList = append(parcelResultList, UpCharacterSkill(characterInfo, skillUp.Level, skillUp.SkillSlot)...)
	}
	rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResultList)
}

func UpCharacterSkill(characterInfo *sro.CharacterInfo, reqLevel int32, SkillSlot string) []*game.ParcelResult {
	conf := gdconf.GetCharacterSkillListExcelTable(characterInfo.CharacterId, 0)
	if conf == nil {
		return nil
	}

	forFunc := func(groupIdList []string) string {
		var forNum = 0
		for _, id := range groupIdList {
			forNum++
			if !strings.Contains(id, "EmptySkill") {
				return id
			}
			if forNum > 10 {
				break
			}
		}
		return ""
	}

	var level int32
	var parcelResultList []*game.ParcelResult
	for level < reqLevel {
		var groupId string
		var skillType string

		if strings.Contains(SkillSlot, "ExSkill") { // ex 技能
			level = characterInfo.ExSkillLevel
			skillType = "ExSkill"
			groupId = forFunc(conf.ExSkillGroupId)
		} else if strings.Contains(SkillSlot, "PublicSkill") { // 第二个 技能
			level = characterInfo.CommonSkillLevel
			skillType = "PublicSkill"
			groupId = forFunc(conf.PublicSkillGroupId)
		} else if strings.Contains(SkillSlot, "ExtraPassive") { // 第四个 技能
			level = characterInfo.ExtraPassiveSkillLevel
			skillType = "ExtraPassive"
			groupId = forFunc(conf.ExtraPassiveSkillGroupId)
		} else if strings.Contains(SkillSlot, "Passive") { // 第三个 技能
			level = characterInfo.PassiveSkillLevel
			skillType = "Passive"
			if len(conf.PassiveSkillGroupId) != 0 {
				groupId = forFunc(conf.PassiveSkillGroupId)
			} else {
				groupId = forFunc(conf.NormalSkillGroupId)
			}
		}

		skillConf := gdconf.GetSkillExcelTable(groupId, level)
		if skillConf == nil {
			return parcelResultList
		}
		recConf := gdconf.GetRecipeIngredientExcelTable(skillConf.RequireLevelUpMaterial)
		if recConf == nil {
			return parcelResultList
		}
		// TODO 验证材料是否有那么多！！！！！！！！！！！！
		parcelResultList = append(parcelResultList, game.GetParcelResultList(recConf.CostParcelType,
			recConf.CostId, recConf.CostAmount, true)...)
		parcelResultList = append(parcelResultList, game.GetParcelResultList(recConf.IngredientParcelType,
			recConf.IngredientId, recConf.IngredientAmount, true)...)
		switch skillType {
		case "ExSkill":
			characterInfo.ExSkillLevel++
		case "PublicSkill":
			characterInfo.CommonSkillLevel++
		case "Passive":
			characterInfo.PassiveSkillLevel++
		case "ExtraPassive":
			characterInfo.ExtraPassiveSkillLevel++
		}
		level++
	}
	return parcelResultList
}
