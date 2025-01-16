package game

import (
	"math/rand"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/mx/proto"
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

func GachaRun(num int, ssr bool, sr bool) []int64 {
	results := make([]int64, 0)
	fn := func(conf []*sro.CharacterExcelTable) int64 {
		index := rand.Intn(len(conf))
		return conf[index].Id
	}
	conf := gdconf.GetCharacterExcelStruct()

	isSr := false
	for i := 0; i < num; i++ {
		index := rand.Intn(100000) + 1
		var result int64
		if ssr && (num == 1 || i == num-2) { // 服务端控制是否必出ssr
			result = fn(conf.CharacterSSRMap)
		} else if sr && (num == 1 || i == num-1) { // 服务端控制是否必出sr
			result = fn(conf.CharacterSRMap)
		} else if num == 10 && i == 9 && !isSr { // 保底四星
			isSr = true
			result = fn(conf.CharacterSRMap)

			// 下面是正常概率计算
		} else if index < 700 {
			result = fn(conf.CharacterSSRMap) // up
		} else if index < 700+2300 {
			result = fn(conf.CharacterSSRMap) // ssr
		} else if index < 700+2300+18500 {
			isSr = true
			result = fn(conf.CharacterSRMap) // sr
		} else {
			result = fn(conf.CharacterRMap) // r
		}

		results = append(results, result)
	}

	return results
}

// 保存抽卡结果
func SaveGachaResults(s *enter.Session, results []int64) ([]*proto.GachaResult, map[int64]bool) {
	addItemList := make(map[int64]bool, 0)
	list := make([]*proto.GachaResult, 0)
	for _, id := range results {
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
				continue
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
	return list, addItemList
}
