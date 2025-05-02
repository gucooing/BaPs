package pack

import (
	"github.com/gucooing/BaPs/protocol/mx"
	"strings"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/logger"
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
	for serverId, info := range game.GetGearInfoList(s) {
		characterInfo := s.GetCharacterByKeyId(info.CharacterServerId)
		if characterInfo != nil {
			characterInfo.GearServerId = serverId
		}
		rsp.GearDBs = append(rsp.GearDBs, game.GetGearDB(s, serverId))
	}
}

func CharacterTranscendence(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.CharacterTranscendenceRequest)
	rsp := response.(*proto.CharacterTranscendenceResponse)

	characterInfo := s.GetCharacterByKeyId(req.TargetCharacterServerId)
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

	characterInfo := s.GetCharacterByKeyId(req.TargetCharacterServerId)
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
		characterInfo := s.GetCharacterByKeyId(sid)
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

	characterInfo := s.GetCharacterByKeyId(req.TargetCharacterDBId)
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

	characterInfo := s.GetCharacterByKeyId(req.TargetCharacterDBId)
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
			s.GetCharacterByKeyId(req.CharacterServerId).GetCharacterId())
		rsp.EquipmentDBs = append(rsp.EquipmentDBs, game.GetEquipmentDB(s, req.EquipmentServerId))
	}()
	game.SetCharacterEquipment(s, req.CharacterServerId, req.EquipmentServerId, req.SlotIndex)
}

func CharacterExpGrowth(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.CharacterExpGrowthRequest)
	rsp := response.(*proto.CharacterExpGrowthResponse)

	characterInfo := s.GetCharacterByKeyId(req.TargetCharacterServerId)
	if characterInfo == nil {
		return
	}
	defer func() {
		rsp.AccountCurrencyDB = game.GetAccountCurrencyDB(s)
		rsp.CharacterDB = game.GetCharacterDB(s, characterInfo.CharacterId)
	}()
	consumeResultDB := &proto.ConsumeResultDB{
		UsedItemServerIdAndRemainingCounts: make(map[int64]int64),
	}
	rsp.ConsumeResultDB = consumeResultDB

	// 计算可以获取多少经验
	if req.ConsumeRequestDB == nil {
		return
	}
	for itemServerId, itemNum := range req.ConsumeRequestDB.ConsumeItemServerIdAndCounts {
		itemInfo := game.GetItemInfo(s, s.GetItemByKeyId(itemServerId).GetUniqueId())
		if itemInfo == nil {
			continue
		}
		itemConf := gdconf.GetItemExcelTable(itemInfo.UniqueId)
		if itemConf == nil {
			continue
		}
		if !game.RemoveItem(s, itemInfo.UniqueId, int32(itemNum)) {
			continue
		}
		consumeResultDB.UsedItemServerIdAndRemainingCounts[itemServerId] = int64(itemInfo.StackCount)

		characterInfo.Exp += itemConf.StackableFunction * itemNum
	}
	// 计算升级后的等级
	newLevel, newExp := gdconf.UpCharacterLevel(characterInfo.Level, characterInfo.Exp)
	if newLevel > game.GetAccountLevel(s) { // 避免高出账号等级
		newLevel = game.GetAccountLevel(s)
		conf := gdconf.GetCharacterLevelExcelTable(newLevel)
		if conf == nil {
			return
		}
		newExp = conf.Exp - 1
	}
	characterInfo.Level = newLevel
	characterInfo.Exp = newExp
}

func CharacterGearUnlock(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.CharacterGearUnlockRequest)
	rsp := response.(*proto.CharacterGearUnlockResponse)

	bin := game.GetGearInfoList(s)
	characterInfo := s.GetCharacterByKeyId(req.CharacterServerId)
	if bin == nil || characterInfo == nil {
		return
	}
	defer func() {
		rsp.CharacterDB = game.GetCharacterDB(s, characterInfo.CharacterId)
	}()

	conf := gdconf.GetUnlockCharacterGear(characterInfo.CharacterId)
	if conf == nil {
		return
	}
	sId := game.GetServerId(s)
	bin[sId] = &sro.GearInfo{
		UniqueId:          conf.Id,
		CharacterServerId: characterInfo.ServerId,
		Level:             1,
		ServerId:          sId,
		SlotIndex:         int64(req.SlotIndex),
		Exp:               0,
		Tier:              conf.Tier,
	}
	characterInfo.GearServerId = sId

	rsp.GearDB = game.GetGearDB(s, sId)
}

func CharacterPotentialGrowth(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.CharacterPotentialGrowthRequest)
	rsp := response.(*proto.CharacterPotentialGrowthResponse)

	characterInfo := s.GetCharacterByKeyId(req.TargetCharacterDBId)
	if characterInfo == nil {
		return
	}
	if characterInfo.PotentialStats == nil {
		characterInfo.PotentialStats = game.NewPotentialStats()
	}
	defer func() {
		rsp.CharacterDB = game.GetCharacterDB(s, characterInfo.CharacterId)
	}()

	parcelResultList := make([]*game.ParcelResult, 0)
	for _, reqInfo := range req.PotentialGrowthRequestDBs {
		oldLevel := characterInfo.PotentialStats[int32(reqInfo.Type)]
		conf := gdconf.GetCharacterPotentialExcelType(characterInfo.CharacterId, reqInfo.Type.String())
		if conf == nil {
			continue
		}

	ty:
		for ; oldLevel < reqInfo.Level; oldLevel++ {
			statConf := gdconf.GetCharacterPotentialStatExcel(conf.PotentialStatGroupId, oldLevel)
			if statConf == nil {
				goto ty
			}
			recConf := gdconf.GetRecipeIngredientExcelTable(statConf.RecipeId)
			if recConf == nil {
				goto ty
			}
			// 根据配方计算需要的东西
			parcelResultList = append(parcelResultList, game.GetParcelResultList(recConf.CostParcelType,
				recConf.CostId, recConf.CostAmount, true)...)
			parcelResultList = append(parcelResultList, game.GetParcelResultList(recConf.IngredientParcelType,
				recConf.IngredientId, recConf.IngredientAmount, true)...)
		}

		characterInfo.PotentialStats[int32(reqInfo.Type)] = oldLevel
	}
	rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResultList)
}
