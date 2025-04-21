package gdconf

import (
	"sync"
	"time"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	"github.com/gucooing/BaPs/pkg/mx"
)

func (g *GameConfig) loadItemExcelTable() {
	g.GetExcel().ItemExcelTable = make([]*sro.ItemExcelTable, 0)
	name := "ItemExcelTable.json"
	mx.LoadExcelJson(g.excelPath+name, &g.GetExcel().ItemExcelTable)
}

type ItemExcel struct {
	RecruitCoin          *RecruitCoin
	recruitCoinSync      sync.Mutex
	ItemExcelMap         map[int64]*sro.ItemExcelTable
	ItemExcelCategoryMap map[string][]*sro.ItemExcelTable
}

type RecruitCoin struct {
	Item   *sro.ItemExcelTable
	EnTime time.Time
}

func (g *GameConfig) gppItemExcelTable() {
	g.GetGPP().ItemExcel = &ItemExcel{
		ItemExcelMap:         make(map[int64]*sro.ItemExcelTable),
		ItemExcelCategoryMap: make(map[string][]*sro.ItemExcelTable),
		recruitCoinSync:      sync.Mutex{},
	}

	for _, v := range g.GetExcel().GetItemExcelTable() {
		if v.ExpirationDateTime != "" {
			enTime, err := time.Parse("2006-01-02 15:04:05", v.ExpirationDateTime)
			if err != nil {
				continue
			}
			if time.Now().After(enTime) {
				continue
			}
		}
		g.GetGPP().ItemExcel.ItemExcelMap[v.Id] = v
		if g.GetGPP().ItemExcel.ItemExcelCategoryMap[v.ItemCategory] == nil {
			g.GetGPP().ItemExcel.ItemExcelCategoryMap[v.ItemCategory] = make([]*sro.ItemExcelTable, 0)
		}
		g.GetGPP().ItemExcel.ItemExcelCategoryMap[v.ItemCategory] = append(
			g.GetGPP().ItemExcel.ItemExcelCategoryMap[v.ItemCategory], v)
	}

	logger.Info("处理道具配置完成,道具:%v个,类型:%v个", len(g.GetGPP().ItemExcel.ItemExcelMap),
		len(g.GetGPP().ItemExcel.ItemExcelCategoryMap))
}

// 弃用,已经有更好的方法了
// func GetRecruitCoin() *sro.ItemExcelTable {
// 	bin := GC.GetGPP().ItemExcel
// 	bin.recruitCoinSync.Lock()
// 	defer bin.recruitCoinSync.Unlock()
// 	if bin.RecruitCoin != nil &&
// 		!time.Now().After(bin.RecruitCoin.EnTime) {
// 		return bin.RecruitCoin.Item
// 	}
// 	confList := GC.GetGPP().ItemExcel.ItemExcelCategoryMap["RecruitCoin"]
// 	for _, conf := range confList {
// 		if conf.ExpirationDateTime == "2099-12-31 23:59:59" {
// 			continue
// 		}
// 		enTime, err := time.Parse("2006-01-02 15:04:05", conf.ExpirationDateTime)
// 		if err != nil {
// 			continue
// 		}
// 		if !time.Now().After(enTime) {
// 			bin.RecruitCoin = &RecruitCoin{
// 				Item:   conf,
// 				EnTime: enTime,
// 			}
// 			return conf
// 		}
//
// 	}
// 	return nil
// }

func GetItemExcelCategoryMap(itemCategory string) []*sro.ItemExcelTable {
	return GC.GetGPP().ItemExcel.ItemExcelCategoryMap[itemCategory]
}

func GetItemExcelTable(id int64) *sro.ItemExcelTable {
	return GC.GetGPP().ItemExcel.ItemExcelMap[id]
}

func IsItem(id int64) bool {
	return GC.GetGPP().ItemExcel.ItemExcelMap[id] != nil
}
