package pack

import (
	"strings"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/logger"
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

func UpCharacterSkill(characterInfo *sro.CharacterInfo, reqLevel int32, skillSlot string) []*game.ParcelResult {
	conf := gdconf.GetCharacterSkillListExcelTable(characterInfo.CharacterId, 0)
	if conf == nil || len(skillSlot) <= 2 {
		return nil
	}

	getGidFunc := func(groupIdList []string) string {
		for _, id := range groupIdList {
			if !strings.Contains(id, "EmptySkill") {
				return id
			}
		}
		return ""
	}

	var level *int32
	var parcelResultList []*game.ParcelResult
	var forNum int32
	for {
		var groupId string
		if forNum >= reqLevel { // 最大递归次数限制
			return parcelResultList
		}
		forNum++

		skillType := skillSlot[:len(skillSlot)-2]
		switch skillType {
		case "ExSkill": // ex 技能
			level = &characterInfo.ExSkillLevel
			groupId = getGidFunc(conf.ExSkillGroupId)
		case "PublicSkill": // 第二个 技能
			level = &characterInfo.CommonSkillLevel
			groupId = getGidFunc(conf.PublicSkillGroupId)
		case "Passive": // 第三个 技能
			level = &characterInfo.PassiveSkillLevel
			if len(conf.PassiveSkillGroupId) != 0 {
				groupId = getGidFunc(conf.PassiveSkillGroupId)
			} else {
				groupId = getGidFunc(conf.NormalSkillGroupId)
			}
		case "ExtraPassive": // 第四个 技能
			level = &characterInfo.ExtraPassiveSkillLevel
			groupId = getGidFunc(conf.ExtraPassiveSkillGroupId)
		default:
			logger.Warn("未知的角色技能升级请求skillSlot:%s", skillSlot)
			return parcelResultList
		}

		skillConf := gdconf.GetSkillExcelTable(groupId, *level)
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
		*level++
		if *level >= reqLevel {
			return parcelResultList
		}
	}
}

func EquipmentEquip(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.EquipmentItemEquipRequest)
	rsp := response.(*proto.EquipmentItemEquipResponse)

	rsp.EquipmentDBs = make([]*proto.EquipmentDB, 0)
	defer func() {
		rsp.CharacterDB = game.GetCharacterDB(s,
			game.ServerIdToCharacterId(s, req.CharacterServerId))
		rsp.EquipmentDBs = append(rsp.EquipmentDBs, game.GetEquipmentDB(s, req.EquipmentServerId))
	}()
	game.SetCharacterEquipment(s, req.CharacterServerId, req.EquipmentServerId, req.SlotIndex)
}
