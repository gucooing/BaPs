package game

import (
	"math/rand"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/protocol/proto"
)

func GetGachaBin(s *enter.Session) *sro.GachaBin {
	bin := GetPlayerBin(s)
	if bin == nil {
		return nil
	}
	if bin.GachaBin == nil {
		bin.GachaBin = &sro.GachaBin{}
	}
	return bin.GachaBin
}

func GetBeforehandInfo(s *enter.Session) *sro.BeforehandInfo {
	bin := GetGachaBin(s)
	if bin == nil {
		return nil
	}
	if bin.BeforehandInfo == nil {
		bin.BeforehandInfo = &sro.BeforehandInfo{}
	}
	return bin.BeforehandInfo
}

func GetBeforehandGachaSnapshotDB(s *enter.Session) *proto.BeforehandGachaSnapshotDB {
	bin := GetBeforehandInfo(s)
	info := &proto.BeforehandGachaSnapshotDB{
		ShopUniqueId: bin.ShopUniqueId,
		GoodsId:      bin.GoodsId,
		LastIndex:    bin.LastIndex,
		LastResults:  bin.LastResults,
		SavedIndex:   bin.SavedIndex,
		SavedResults: bin.SavedResults,
		PickedIndex:  0,
	}

	return info
}

/*
3 0.03 -> 3000 -> 700 + 2300
2 0.185 -> 18500
1 0.785 -> 78500
*/

func GenGachaResults(goodsId int64) []*ParcelResult {
	results := make([]*ParcelResult, 0)
	goods := gdconf.GetGoodsExcelTable(goodsId)
	addDef := func() {
		results = append(results, &ParcelResult{
			ParcelType: proto.ParcelType_Character,
			ParcelId:   16000,
			Amount:     1,
		})
	}
	if goods == nil {
		addDef()
		return results
	}
	genGach := func(gachaGoodsId, gachaNum int64) {
		// 概率生成
		ges := make(map[int]*gdconf.GachaElementGroupId)
		var probabilityList []int
		probability := 0
		for _, ge := range gdconf.GetGachaElementExcelTableByGachaGroupId(gachaGoodsId) {
			if ge == nil {
				continue
			}
			switch ge.Rarity {
			case "SSR": // 0.03
				if len(ge.GachaElementExcelList) == 1 {
					probability += 700
				} else {
					probability += 2300
				}
			case "SR": // 0.185
				probability += 18500
			case "R": // 0.785
				probability += 78500
			}
			ges[probability] = ge
			probabilityList = append(probabilityList, probability)
		}

		// 抽卡
		for i := int64(0); i < gachaNum; i++ {
			index := rand.Intn(probability) + 1
			var ge *gdconf.GachaElementGroupId
			for _, prob := range probabilityList {
				if index < prob {
					ge = ges[prob]
					break
				}
			}
			if ge == nil {
				addDef()
			} else {
				gee := ge.GachaElementExcelList[rand.Intn(len(ge.GachaElementExcelList))]
				if gee == nil {
					addDef()
					continue
				}
				results = append(results, &ParcelResult{
					ParcelType: proto.GetParcelTypeValue(gee.ParcelType_),
					ParcelId:   gee.ParcelId,
					Amount:     1,
				})
			}
		}
	}

	for index, pt := range goods.ParcelType_ {
		switch pt {
		case "Item", "Character":
			results = append(results, &ParcelResult{
				ParcelType: proto.GetParcelTypeValue(pt),
				ParcelId:   goods.ParcelId[index],
				Amount:     goods.ParcelAmount[index],
			})
		case "GachaGroup":
			genGach(goods.ParcelId[index], goods.ParcelAmount[index])
		default:
			logger.Error("未处理的卡池属性:%s", pt)
		}
	}

	return results
}

// 弃用,已有正规方法
// func GachaRun(num int64, ssr bool, sr bool) []int64 {
// 	results := make([]int64, 0)
// 	fn := func(conf []*sro.CharacterExcelTable) int64 {
// 		index := rand.Intn(len(conf))
// 		return conf[index].Id
// 	}
// 	conf := gdconf.GetCharacterExcelStruct()
//
// 	isSr := false
// 	for i := int64(0); i < num; i++ {
// 		index := rand.Intn(100000) + 1
// 		var result int64
// 		if ssr && (num == 1 || i == num-2) { // 服务端控制是否必出ssr
// 			ssr = false
// 			result = fn(conf.CharacterSSRMap)
// 		} else if sr && (num == 1 || i == num-1) { // 服务端控制是否必出sr
// 			sr = false
// 			result = fn(conf.CharacterSRMap)
// 		} else if num == 10 && i == 9 && !isSr { // 保底四星
// 			isSr = true
// 			result = fn(conf.CharacterSRMap)
//
// 			// 下面是正常概率计算
// 		} else if index < 700 {
// 			result = fn(conf.CharacterSSRMap) // up
// 		} else if index < 700+2300 {
// 			result = fn(conf.CharacterSSRMap) // ssr
// 		} else if index < 700+2300+18500 {
// 			isSr = true
// 			result = fn(conf.CharacterSRMap) // sr
// 		} else {
// 			result = fn(conf.CharacterRMap) // r
// 		}
//
// 		results = append(results, result)
// 	}
//
// 	return results
// }

// 保存抽卡结果
func SaveGachaResults(s *enter.Session, results []*ParcelResult) ([]*proto.GachaResult, map[int64]bool) {
	addItemList := make(map[int64]bool, 0)
	list := make([]*proto.GachaResult, 0)
	for _, result := range results {
		addCharacter := func(id int64) {
			gachaResult := &proto.GachaResult{
				CharacterId: id,
				Character:   nil,
				Stone:       nil,
			}
			if AddCharacter(s, id) {
				gachaResult.Character = GetCharacterDB(s, id)
			} else { // 重复添加
				conf := gdconf.GetCharacterExcel(id)
				if conf == nil {
					return
				}
				// 添加秘石
				secretStoneServerId := AddItem(s, conf.SecretStoneItemId, conf.SecretStoneItemAmount)
				gachaResult.Stone = &proto.ItemDB{
					Type:       proto.ParcelType_Item,
					ServerId:   secretStoneServerId,
					UniqueId:   conf.SecretStoneItemId,
					StackCount: conf.SecretStoneItemAmount,
				}
				addItemList[conf.SecretStoneItemId] = true
				// 添加碎片
				AddItem(s, conf.CharacterPieceItemId, conf.CharacterPieceItemAmount)
				addItemList[conf.CharacterPieceItemId] = true
			}
			list = append(list, gachaResult)
		}
		switch result.ParcelType {
		case proto.ParcelType_Character:
			addCharacter(result.ParcelId)
		case proto.ParcelType_Item:
			AddItem(s, result.ParcelId, int32(result.Amount))
			addItemList[result.ParcelId] = true
		}

	}
	return list, addItemList
}
