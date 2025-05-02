package game

import (
	"github.com/gucooing/BaPs/protocol/mx"
	"math/rand"
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/protocol/proto"
)

func DefaultCafeBin(s *enter.Session) *sro.CafeBin {
	bin := &sro.CafeBin{
		CafeInfoList:      make(map[int64]*sro.CafeInfo),
		FurnitureInfoList: make(map[int64]*sro.FurnitureInfo),
	}
	for _, conf := range gdconf.GetCafeInfoExcelTables() {
		if !conf.IsDefault {
			continue
		}
		NewCafeBin(s, bin, conf.CafeId)
	}

	return bin
}

func NewCafeBin(s *enter.Session, bin *sro.CafeBin, cafeId int64) *sro.CafeInfo {
	if bin == nil {
		return nil
	}
	if bin.CafeInfoList == nil {
		bin.CafeInfoList = make(map[int64]*sro.CafeInfo)
	}
	if bin.FurnitureInfoList == nil {
		bin.FurnitureInfoList = make(map[int64]*sro.FurnitureInfo)
	}
	conf := gdconf.GetCafeInfoExcelTableInfo(cafeId)
	if conf == nil {
		return nil
	}
	sid := GetServerId(s)
	info := &sro.CafeInfo{
		CafeId:         conf.CafeId,
		CafeRank:       1,
		ServerId:       sid,
		FurnitureList:  make(map[int64]bool),
		ProductionList: NewRroductionList(cafeId),
	}
	for _, furnitureConf := range gdconf.GetDefaultFurnitureExcelList() {
		furnitureSid := GetServerId(s)
		furnitureInfo := &sro.FurnitureInfo{
			FurnitureId:  furnitureConf.Id,
			StackCount:   1,
			ServerId:     furnitureSid,
			Location:     proto.FurnitureLocation_value[furnitureConf.Location], // 摆放位置
			CafeServerId: sid,
			PositionX:    furnitureConf.PositionX,
			PositionY:    furnitureConf.PositionY,
			Rotation:     furnitureConf.Rotation,
		}
		bin.FurnitureInfoList[furnitureSid] = furnitureInfo
		info.FurnitureList[furnitureSid] = true
	}
	bin.CafeInfoList[sid] = info

	return info
}

func NewRroductionList(cafeId int64) map[int64]*sro.ProductionInfo {
	list := make(map[int64]*sro.ProductionInfo)
	for _, conf := range gdconf.GetCafeProductionExcelTableList(cafeId, 1) {
		list[conf.CafeProductionParcelId] = &sro.ProductionInfo{
			ParcelId:   conf.CafeProductionParcelId,
			ParcelType: int64(proto.ParcelType_value[conf.CafeProductionParcelType]),
			Amount:     0,
		}
	}
	return list
}

func GetCafeBin(s *enter.Session) *sro.CafeBin {
	bin := GetPlayerBin(s)
	if bin == nil {
		return nil
	}
	if bin.CafeBin == nil {
		bin.CafeBin = DefaultCafeBin(s)
	}
	return bin.CafeBin
}

func GetCafeInfoList(s *enter.Session) map[int64]*sro.CafeInfo {
	bin := GetCafeBin(s)
	if bin == nil {
		return nil
	}
	if bin.CafeInfoList == nil {
		bin.CafeInfoList = make(map[int64]*sro.CafeInfo)
	}
	return bin.CafeInfoList
}

func GetCafeInfo(s *enter.Session, serverId int64) *sro.CafeInfo {
	bin := GetCafeInfoList(s)
	if bin == nil {
		return nil
	}
	UpCafeVisitCharacterDB(bin[serverId])
	return bin[serverId]
}

func UpCafeVisitCharacterDB(bin *sro.CafeInfo) {
	if bin == nil {
		return
	}
	cafeRankConf := gdconf.GetCafeRankExcelTable(bin.CafeId, bin.CafeRank)
	produConfList := gdconf.GetCafeProductionExcelTableList(bin.CafeId, bin.CafeRank)
	if cafeRankConf == nil || produConfList == nil {
		return
	}
	// 学生刷新
	if alg.GetLastTimeHourH(1).After(time.Unix(bin.LastUpdate, 0)) {
		bin.IsNew = true
		bin.VisitCharacterList = make(map[int64]*sro.VisitCharacterInfo)
		characterNum := int32(0)
		if cafeRankConf.CharacterVisitMax > cafeRankConf.CharacterVisitMin {
			characterNum = rand.Int31n(cafeRankConf.CharacterVisitMax-cafeRankConf.CharacterVisitMin) +
				cafeRankConf.CharacterVisitMin
		} else {
			characterNum = cafeRankConf.CharacterVisitMin
		}
		// 刷新
		for i := int32(0); i < characterNum; i++ {
			characterId := gdconf.RandCharacter()
			if bin.VisitCharacterList[characterId] != nil {
				i--
				continue
			}
			bin.VisitCharacterList[characterId] = &sro.VisitCharacterInfo{
				CharacterId: characterId,
			}
		}
		bin.LastUpdate = time.Now().Unix()
	}
	// 产量计算
	if bin.ProductionList == nil {
		bin.ProductionList = NewRroductionList(bin.CafeId)
	}
	// 公式 ParcelProductionCorrectionValue + ComfortValue/MaxComfortValue * ParcelProductionCoefficient*(0.45,0.2)
	num := time.Now().Sub(time.Unix(bin.ProductionAppliedTime, 0)) / time.Hour
	for _, prodBin := range bin.ProductionList {
		produConf, ok := produConfList[prodBin.ParcelId]
		if !ok {
			continue
		}
		base := produConf.ParcelProductionCorrectionValue / 100
		additional := float32(bin.ComfortValue/cafeRankConf.ComfortMax) * produConf.ParcelProductionCoefficient
		if bin.CafeId == 1 {
			prodBin.Amount = int64((base + additional*45) * float32(num))
		} else {
			prodBin.Amount = int64((base + additional*20) * float32(num))
		}
		prodBin.Amount = alg.MinInt64(prodBin.Amount, produConf.ParcelStorageMax*100) // 最大值限制
	}
	bin.ProductionAppliedNum = int32(num)
	// 不累加，直接覆盖，原因是这个时间是领取的时候才刷新
}

func UpCafeComfortValue(s *enter.Session, cafeServerId int64) {
	cafeBin := GetCafeInfo(s, cafeServerId)
	if cafeBin == nil {
		return
	}
	cafeRankConf := gdconf.GetCafeRankExcelTable(cafeBin.CafeId, cafeBin.CafeRank)
	if cafeRankConf == nil {
		return
	}
	cafeBin.ComfortValue = 0 // 重置掉
	for furnitureServerId, ok := range cafeBin.FurnitureList {
		if !ok {
			continue
		}
		furnitureBin := GetFurnitureInfo(s, furnitureServerId)
		if furnitureBin == nil {
			continue
		}
		conf := gdconf.GetFurnitureExcelTable(furnitureBin.FurnitureId)
		if conf == nil {
			continue
		}
		cafeBin.ComfortValue += conf.ComfortBonus
	}
	cafeBin.ComfortValue = alg.MinInt64(cafeBin.ComfortValue, cafeRankConf.ComfortMax)
}

func GetFurnitureInfoList(s *enter.Session) map[int64]*sro.FurnitureInfo {
	bin := GetCafeBin(s)
	if bin == nil {
		return nil
	}
	if bin.FurnitureInfoList == nil {
		bin.FurnitureInfoList = make(map[int64]*sro.FurnitureInfo)
	}
	return bin.FurnitureInfoList
}

func GetFurnitureInfo(s *enter.Session, serverId int64) *sro.FurnitureInfo {
	bin := GetFurnitureInfoList(s)
	if bin == nil {
		return nil
	}
	return bin[serverId]
}

func AddFurnitureInfo(s *enter.Session, furnitureId int64, num int64) []int64 {
	list := make([]int64, 0)
	conf := gdconf.GetFurnitureExcelTable(furnitureId)
	if conf == nil {
		return list
	}
	binList := GetFurnitureInfoList(s)
	for i := int64(0); i < num; i++ {
		sid := GetServerId(s)
		info := &sro.FurnitureInfo{
			FurnitureId:  conf.Id,
			StackCount:   1,
			ServerId:     sid,
			Location:     int32(proto.FurnitureLocation_Inventory),
			CafeServerId: 0,
			PositionX:    0,
			PositionY:    0,
			Rotation:     0,
		}
		binList[sid] = info
		list = append(list, sid)
	}

	return list
}

func RemoveFurniture(s *enter.Session, furnitureServerId int64, cafeServerId int64) {
	cafeInfo := GetCafeInfo(s, cafeServerId)
	if cafeInfo == nil {
		return
	}
	furnitureInfo := GetFurnitureInfo(s, furnitureServerId)
	if furnitureInfo == nil ||
		furnitureInfo.PositionX == furnitureInfo.PositionY &&
			furnitureInfo.PositionY == furnitureInfo.Rotation &&
			furnitureInfo.Rotation == 0 {
		return
	}
	furnitureInfo.CafeServerId = 0
	furnitureInfo.Location = int32(proto.FurnitureLocation_Inventory)
	delete(cafeInfo.FurnitureList, furnitureServerId)
}

func DeployRelocateFurniture(s *enter.Session, furnitureDB *proto.FurnitureDB, cafeServerId int64) int64 {
	cafeInfo := GetCafeInfo(s, cafeServerId)
	furnitureInfo := GetFurnitureInfo(s, furnitureDB.ServerId)
	if furnitureInfo == nil || cafeInfo == nil {
		return 0
	}
	if furnitureDB.PositionX == furnitureDB.PositionY &&
		furnitureDB.PositionY == furnitureDB.Rotation &&
		furnitureDB.Rotation == 0 { //唯一家具
		for oldSid, ok := range cafeInfo.FurnitureList {
			if oldInfo := GetFurnitureInfo(s, oldSid); oldInfo != nil && ok &&
				proto.FurnitureLocation(oldInfo.Location) == furnitureDB.Location &&
				oldInfo.PositionX == oldInfo.PositionY &&
				oldInfo.PositionY == oldInfo.Rotation &&
				oldInfo.Rotation == 0 {
				oldInfo.CafeServerId = 0
				oldInfo.Location = int32(proto.FurnitureLocation_Inventory)
				delete(cafeInfo.FurnitureList, oldSid)
				break
			}
		}
	}

	furnitureInfo.CafeServerId = furnitureDB.CafeDBId
	furnitureInfo.Location = int32(furnitureDB.Location)
	furnitureInfo.PositionX = furnitureDB.PositionX
	furnitureInfo.PositionY = furnitureDB.PositionY
	furnitureInfo.Rotation = furnitureDB.Rotation

	if cafeInfo.FurnitureList == nil {
		cafeInfo.FurnitureList = make(map[int64]bool)
	}
	cafeInfo.FurnitureList[furnitureInfo.ServerId] = true
	UpCafeComfortValue(s, cafeServerId)

	return furnitureInfo.ServerId
}

func GetPbCafeDBs(s *enter.Session) []*proto.CafeDB {
	list := make([]*proto.CafeDB, 0)
	for _, bin := range GetCafeInfoList(s) {
		list = append(list, GetCafeDB(s, bin.ServerId))
	}

	return list
}

func GetCafeDB(s *enter.Session, serverId int64) *proto.CafeDB {
	bin := GetCafeInfo(s, serverId)
	info := &proto.CafeDB{
		CafeDBId:              bin.ServerId,
		CafeId:                bin.CafeId,                             // 咖啡厅id
		AccountId:             s.AccountServerId,                      // 账号id
		CafeRank:              bin.CafeRank,                           // 咖啡厅等级
		CafeVisitCharacterDBs: make(map[int64]*proto.CafeCharacterDB), // 入场学生
		FurnitureDBs:          make([]*proto.FurnitureDB, 0),          // 摆放的家具
		IsNew:                 bin.IsNew,                              // 是否新的
		ProductionDB:          nil,
		LastUpdate:            mx.Unix(bin.LastUpdate, 0),
		LastSummonDate:        mx.Unix(bin.SummonUpdate, 0),
		ProductionAppliedTime: mx.Unix(bin.ProductionAppliedTime, 0).Add(time.Duration(bin.ProductionAppliedNum) * time.Hour),
	}

	for _, visitCharacterInfo := range bin.VisitCharacterList {
		cafeCharacterDB := &proto.CafeCharacterDB{
			VisitingCharacterDB: &proto.VisitingCharacterDB{
				UniqueId: visitCharacterInfo.CharacterId,
				ServerId: 0,
			},
			LastInteractTime: mx.Unix(visitCharacterInfo.LastInteractTime, 0),
			IsSummon:         visitCharacterInfo.IsSummon,
		}
		if cafeCharacterBin := GetCharacterInfo(s, visitCharacterInfo.CharacterId); cafeCharacterBin != nil {
			cafeCharacterDB.ServerId = cafeCharacterBin.ServerId
		}
		info.CafeVisitCharacterDBs[visitCharacterInfo.CharacterId] = cafeCharacterDB
	}

	for id, ok := range bin.FurnitureList {
		if ok {
			info.FurnitureDBs = append(info.FurnitureDBs, GetFurnitureDB(s, id))
		}
	}

	productionDB := &proto.CafeProductionDB{
		CafeDBId:              bin.ServerId,
		ComfortValue:          bin.ComfortValue,
		AppliedDate:           info.ProductionAppliedTime,
		ProductionParcelInfos: make([]*proto.CafeProductionParcelInfo, 0),
	}
	for _, productionBin := range bin.ProductionList {
		productionDB.ProductionParcelInfos = append(productionDB.ProductionParcelInfos, &proto.CafeProductionParcelInfo{
			Key: &proto.ParcelKeyPair{
				Type: proto.ParcelType(productionBin.ParcelType),
				Id:   productionBin.ParcelId,
			},
			Amount: productionBin.Amount,
		})
	}

	info.ProductionDB = productionDB

	return info
}

func GetFurnitureDBs(s *enter.Session) []*proto.FurnitureDB {
	list := make([]*proto.FurnitureDB, 0)
	for _, bin := range GetFurnitureInfoList(s) {
		info := &proto.FurnitureDB{
			Type:               proto.ParcelType_Furniture,
			Location:           proto.FurnitureLocation(bin.Location),
			CafeDBId:           bin.CafeServerId,
			PositionX:          bin.PositionX,
			PositionY:          bin.PositionY,
			Rotation:           bin.Rotation,
			ItemDeploySequence: 0,
			ConsumableItemBaseDB: &proto.ConsumableItemBaseDB{
				ServerId:   bin.ServerId,
				UniqueId:   bin.FurnitureId,
				StackCount: bin.StackCount,
				Key:        nil,
				CanConsume: false,
			},
		}
		if bin.CafeServerId != 0 {
			info.ItemDeploySequence = bin.ServerId
		}
		list = append(list, info)
	}
	return list
}

func GetFurnitureDB(s *enter.Session, serverId int64) *proto.FurnitureDB {
	bin := GetFurnitureInfo(s, serverId)
	if bin == nil {
		return nil
	}
	info := &proto.FurnitureDB{
		Type:               proto.ParcelType_Furniture,
		Location:           proto.FurnitureLocation(bin.Location),
		CafeDBId:           bin.CafeServerId,
		PositionX:          bin.PositionX,
		PositionY:          bin.PositionY,
		Rotation:           bin.Rotation,
		ItemDeploySequence: 0,
		ConsumableItemBaseDB: &proto.ConsumableItemBaseDB{
			ServerId:   bin.ServerId,
			UniqueId:   bin.FurnitureId,
			StackCount: bin.StackCount,
			Key:        nil,
			CanConsume: false,
		},
	}
	if bin.CafeServerId != 0 {
		info.ItemDeploySequence = bin.ServerId
	}
	return info
}
