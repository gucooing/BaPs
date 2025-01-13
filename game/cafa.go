package game

import (
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/mx/proto"
)

func NewCafe() *sro.CafeBin {
	bin := &sro.CafeBin{
		CafeInfoList: make(map[int64]*sro.CafeInfo),
	}
	for _, conf := range gdconf.GetCafeInfoExcelTables() {
		if !conf.IsDefault {
			continue
		}
		bin.CafeInfoList[conf.CafeId] = &sro.CafeInfo{
			CafeId:                conf.CafeId,
			LaseUpdateTime:        time.Now().Unix(),
			CafeRank:              1,
			ProductionAppliedTime: time.Now().Unix(),
		}
	}

	return bin
}

func GetCafeBin(s *enter.Session) *sro.CafeBin {
	bin := GetPlayerBin(s)
	if bin == nil {
		return nil
	}
	if bin.CafeBin == nil {
		bin.CafeBin = NewCafe()
	}
	return bin.CafeBin
}

func GetPbCafeDBs(s *enter.Session) []*proto.CafeDB {
	list := make([]*proto.CafeDB, 0)
	for _, bin := range GetCafeBin(s).GetCafeInfoList() {
		db := &proto.CafeDB{
			CafeDBId:              GetDBId(),
			CafeId:                bin.CafeId,        // 咖啡厅id
			AccountId:             s.AccountServerId, // 账号id
			CafeRank:              bin.CafeRank,      // 咖啡厅等级
			LastUpdate:            time.Unix(bin.LaseUpdateTime, 0),
			LastSummonDate:        time.Time{},
			IsNew:                 true,                                   // 是否新的
			CafeVisitCharacterDBs: make(map[int64]*proto.CafeCharacterDB), // 入场学生
			FurnitureDBs:          make([]*proto.FurnitureDB, 0),          // 摆放的家具
			ProductionAppliedTime: time.Unix(bin.ProductionAppliedTime, 0),
			ProductionDB:          nil,
		}

		for _, id := range []int64{16000, 16007, 10031} {
			db.CafeVisitCharacterDBs[id] = &proto.CafeCharacterDB{
				UniqueId: id,
			}
		}

		productionDB := &proto.CafeProductionDB{
			CafeDBId:     GetDBId(),
			ComfortValue: 0,
			AppliedDate:  time.Unix(bin.ProductionAppliedTime, 0),
			ProductionParcelInfos: []*proto.CafeProductionParcelInfo{
				{
					Key: proto.ParcelKeyPair{
						Type: proto.ParcelType_Currency,
						Id:   1,
					},
					Amount: 0,
				},
				{
					Key: proto.ParcelKeyPair{
						Type: proto.ParcelType_Currency,
						Id:   5,
					},
					Amount: 0,
				},
			},
		}

		db.ProductionDB = productionDB

		list = append(list, db)
	}

	return list
}
