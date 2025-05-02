package pack

import (
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/protocol/mx"
	"github.com/gucooing/BaPs/protocol/proto"
)

func AccountCurrencySync(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.AccountCurrencySyncResponse)

	rsp.AccountCurrencyDB = game.GetAccountCurrencyDB(s)
	rsp.ExpiredCurrency = make(map[proto.CurrencyTypes]int64)
}

func ItemList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.ItemListResponse)

	rsp.ExpiryItemDBs = make([]*proto.ItemDB, 0)
	rsp.ItemDBs = make([]*proto.ItemDB, 0)

	for id, conf := range game.GetItemList(s) {
		if !gdconf.IsItem(conf.UniqueId) {
			delete(game.GetItemList(s), id)
			continue
		}
		rsp.ItemDBs = append(rsp.ItemDBs, &proto.ItemDB{
			Type: proto.ParcelType_Item,
			ConsumableItemBaseDB: &proto.ConsumableItemBaseDB{
				Key:        nil,
				CanConsume: false,
				ServerId:   conf.ServerId,
				UniqueId:   conf.UniqueId,
				StackCount: int64(conf.StackCount),
			},
		})
	}
}

func EquipmentList(s *enter.Session, request, response mx.Message) {
	rsp := response.(*proto.EquipmentItemListResponse)

	rsp.EquipmentDBs = game.GetEquipmentDBs(s)
}

func EquipmentLevelUp(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.EquipmentItemLevelUpRequest)
	rsp := response.(*proto.EquipmentItemLevelUpResponse)

	defer func() {
		rsp.EquipmentDB = game.GetEquipmentDB(s, req.TargetServerId)
		rsp.AccountCurrencyDB = game.GetAccountCurrencyDB(s)
	}()
	bin := game.GetEquipmentInfo(s, req.TargetServerId)
	if bin == nil {
		return
	}
	conf := gdconf.GetEquipmentExcelTable(bin.UniqueId)
	if conf == nil {
		return
	}
	// 材料处理
	consumeResultDB := &proto.ConsumeResultDB{
		RemovedItemServerIds:                    make([]int64, 0),
		RemovedEquipmentServerIds:               make([]int64, 0),
		RemovedFurnitureServerIds:               make([]int64, 0),
		UsedItemServerIdAndRemainingCounts:      make(map[int64]int64),
		UsedEquipmentServerIdAndRemainingCounts: make(map[int64]int64),
		UsedFurnitureServerIdAndRemainingCounts: make(map[int64]int64),
	}
	rsp.ConsumeResultDB = consumeResultDB
	for serverId, num := range req.ConsumeRequestDB.ConsumeEquipmentServerIdAndCounts {
		ok, feedExp := game.DelEquipment(s, serverId, num)
		if !ok {
			continue
		}
		bin.Exp += feedExp
		equipInfo := game.GetEquipmentInfo(s, serverId)
		if equipInfo == nil {
			continue
		}
		equipConf := gdconf.GetEquipmentExcelTable(equipInfo.UniqueId)
		if equipConf.MaxLevel < 10 {
			consumeResultDB.UsedEquipmentServerIdAndRemainingCounts[serverId] = equipInfo.StackCount
		} else {
			consumeResultDB.RemovedEquipmentServerIds = append(consumeResultDB.RemovedEquipmentServerIds,
				serverId)
		}
	}
	// 该计算可以升到几级了
	for {
		levelConf := gdconf.GetEquipmentLevelExcelTable(bin.Level)
		if levelConf == nil {
			break
		}
		if bin.Exp < levelConf.TierLevelExp[len(levelConf.TierLevelExp)-1] ||
			bin.Level >= conf.MaxLevel {
			break
		}
		bin.Exp -= levelConf.TierLevelExp[len(levelConf.TierLevelExp)-1]
		bin.Level++
	}
}

func EquipmentTierUp(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.EquipmentItemTierUpRequest)
	rsp := response.(*proto.EquipmentItemTierUpResponse)

	defer func() {
		rsp.EquipmentDB = game.GetEquipmentDB(s, req.TargetEquipmentServerId)
	}()
	bin := game.GetEquipmentInfo(s, req.TargetEquipmentServerId)
	if bin == nil {
		return
	}
	conf := gdconf.GetEquipmentExcelTable(bin.UniqueId)
	if conf == nil || conf.NextTierEquipment == 0 {
		return
	}
	// 升级蓝图扣除
	recConf := gdconf.GetRecipeIngredientExcelTable(conf.RecipeId)
	if recConf == nil {
		return
	}
	// 根据配方计算需要的东西
	parcelResultList := game.GetParcelResultList(recConf.CostParcelType,
		recConf.CostId, recConf.CostAmount, true)
	parcelResultList = append(parcelResultList, game.GetParcelResultList(recConf.IngredientParcelType,
		recConf.IngredientId, recConf.IngredientAmount, true)...)
	bin.UniqueId = conf.NextTierEquipment
	bin.Tier++

	rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResultList)
}

func EquipmentBatchGrowth(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.EquipmentBatchGrowthRequest)
	rsp := response.(*proto.EquipmentBatchGrowthResponse)

	rsp.EquipmentDBs = make([]*proto.EquipmentDB, 0)
	consumeResultDB := &proto.ConsumeResultDB{
		RemovedItemServerIds:                    make([]int64, 0),
		RemovedEquipmentServerIds:               make([]int64, 0),
		RemovedFurnitureServerIds:               make([]int64, 0),
		UsedItemServerIdAndRemainingCounts:      make(map[int64]int64),
		UsedEquipmentServerIdAndRemainingCounts: make(map[int64]int64),
		UsedFurnitureServerIdAndRemainingCounts: make(map[int64]int64),
	}
	rsp.ConsumeResultDB = consumeResultDB
	parcelResultList := make([]*game.ParcelResult, 0)
	// 装备升级
	for _, equipmentBatch := range req.EquipmentBatchGrowthRequestDBs {
		bin := game.GetEquipmentInfo(s, equipmentBatch.TargetServerId)
		if bin == nil {
			continue
		}
		conf := gdconf.GetEquipmentExcelTable(bin.UniqueId)
		if conf == nil {
			continue
		}
		isUp := true
	del:
		for _, consumeRequestDB := range equipmentBatch.ConsumeRequestDBs {
			for serverId, num := range consumeRequestDB.ConsumeEquipmentServerIdAndCounts {
				ok, _ := game.DelEquipment(s, serverId, num)
				if !ok {
					isUp = false
					goto del
				}
				equipInfo := game.GetEquipmentInfo(s, serverId)
				if equipInfo == nil {
					isUp = false
					goto del
				}
				parcelResultList = append(parcelResultList, &game.ParcelResult{
					ParcelType: proto.ParcelType_Equipment,
					ParcelId:   equipInfo.UniqueId,
				})
				equipConf := gdconf.GetEquipmentExcelTable(equipInfo.UniqueId)
				if equipConf.MaxLevel < 10 {
					consumeResultDB.UsedEquipmentServerIdAndRemainingCounts[serverId] = equipInfo.StackCount
				} else {
					consumeResultDB.RemovedEquipmentServerIds = append(consumeResultDB.RemovedEquipmentServerIds,
						serverId)
				}
			}
		}

		if isUp {
			bin.Level = int32(equipmentBatch.AfterLevel)
			bin.Exp = equipmentBatch.AfterExp
		}
		for {
			newConf := gdconf.GetEquipmentExcelTable(bin.UniqueId)
			if newConf == nil || newConf.NextTierEquipment == 0 ||
				bin.Level < newConf.MaxLevel || int32(equipmentBatch.AfterTier) <= newConf.TierInit {
				break
			}

			// 升级蓝图扣除
			recConf := gdconf.GetRecipeIngredientExcelTable(newConf.RecipeId)
			if recConf == nil {
				break
			}
			// 根据配方计算需要的东西
			parcelResultList = append(parcelResultList, game.GetParcelResultList(recConf.CostParcelType,
				recConf.CostId, recConf.CostAmount, true)...)
			parcelResultList = append(parcelResultList, game.GetParcelResultList(recConf.IngredientParcelType,
				recConf.IngredientId, recConf.IngredientAmount, true)...)

			bin.UniqueId = newConf.NextTierEquipment
			bin.Tier++
		}
		rsp.EquipmentDBs = append(rsp.EquipmentDBs, game.GetEquipmentDB(s, equipmentBatch.TargetServerId))
	}

	// 爱用品升级
	if reqGear := req.GearTierUpRequestDB; reqGear != nil {
		defer func() {
			rsp.GearDB = game.GetGearDB(s, reqGear.TargetServerId)
		}()
		bin := game.GetGearInfo(s, reqGear.TargetServerId)
		if bin == nil {
			return
		}
		conf := gdconf.GetCharacterGearExcel(bin.UniqueId)
		if conf == nil {
			return
		}
		recConf := gdconf.GetRecipeIngredientExcelTable(conf.RecipeId)
		if recConf == nil {
			return
		}
		// 根据配方计算需要的东西
		parcelResultList = append(parcelResultList, game.GetParcelResultList(recConf.CostParcelType,
			recConf.CostId, recConf.CostAmount, true)...)
		parcelResultList = append(parcelResultList, game.GetParcelResultList(recConf.IngredientParcelType,
			recConf.IngredientId, recConf.IngredientAmount, true)...)
		// 构造回包 这步是没意义的一步浪费性能
		for _, prInfo := range parcelResultList {
			switch prInfo.ParcelType {
			case proto.ParcelType_Item:
				consumeResultDB.UsedItemServerIdAndRemainingCounts[prInfo.ParcelId] =
					int64(game.GetItemInfo(s, prInfo.ParcelId).GetStackCount())
			}
		}
		// 升级
		bin.UniqueId = conf.NextTierEquipment
		bin.Tier++
	}
	// 删物品
	rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResultList)
}

func CharacterWeaponTranscendence(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.CharacterWeaponTranscendenceRequest)
	rsp := response.(*proto.CharacterWeaponTranscendenceResponse)

	characterInfo := s.GetCharacterByKeyId(req.TargetCharacterServerId)
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

func CharacterWeaponExpGrowth(s *enter.Session, request, response mx.Message) {
	req := request.(*proto.CharacterWeaponExpGrowthRequest)
	rsp := response.(*proto.CharacterWeaponExpGrowthResponse)

	characterInfo := s.GetCharacterByKeyId(req.TargetCharacterServerId)
	if characterInfo == nil {
		return
	}
	waeponInfo := game.GetWeaponInfo(s, characterInfo.CharacterId)
	if waeponInfo == nil {
		return
	}
	parcelResult := make([]*game.ParcelResult, 0)
	parcelResult = append(parcelResult, &game.ParcelResult{
		ParcelType: proto.ParcelType_CharacterWeapon,
		ParcelId:   characterInfo.CharacterId,
	})
	for id, num := range req.ConsumeUniqueIdAndCounts {
		serverId := s.GetEquipmentByKeyId(id).GetServerId()
		ok, feedExp := game.DelEquipment(s, serverId, num)
		if !ok {
			continue
		}
		parcelResult = append(parcelResult, &game.ParcelResult{
			ParcelType: proto.ParcelType_Equipment,
			ParcelId:   id,
		})
		waeponInfo.Exp += feedExp
	}
	// 该升到几级了 excel表没有导出拿不到配置，所以直接拉满得了
	waeponInfo.Exp = 0
	switch waeponInfo.StarGrade {
	case 1:
		waeponInfo.Level = 30
	case 2:
		waeponInfo.Level = 40
	case 3:
		waeponInfo.Level = 50

	}

	rsp.ParcelResultDB = game.ParcelResultDB(s, parcelResult)
}
