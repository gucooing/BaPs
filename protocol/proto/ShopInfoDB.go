// Code auto-generated by gucooing & Hiro420 DO NOT EDIT.

package proto

import (
	"github.com/bytedance/sonic"
	"github.com/gucooing/BaPs/protocol/mx"
)

type ShopInfoDB struct {
	EventContentId      int64
	Category            ShopCategoryType
	ManualRefreshCount  int64
	IsRefresh           bool
	NextAutoRefreshDate mx.MxTime
	LastAutoRefreshDate mx.MxTime
	ShopProductList     []*ShopProductDB
}

func (x *ShopInfoDB) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}
