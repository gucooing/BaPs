package game

import (
	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
)

func NewItemList(s *enter.Session) map[int64]*sro.ItemInfo {
	list := make(map[int64]*sro.ItemInfo)
	list[2] = &sro.ItemInfo{
		ServerId:   GetServerId(s),
		UniqueId:   2,
		StackCount: 5,
	}
	return list
}

func GetItemBin(s *enter.Session) *sro.ItemBin {
	bin := GetPlayerBin(s)
	if bin == nil {
		return nil
	}
	if bin.ItemBin == nil {
		bin.ItemBin = &sro.ItemBin{}
	}
	return bin.ItemBin
}

func GetItemList(s *enter.Session) map[int64]*sro.ItemInfo {
	bin := GetItemBin(s)
	if bin == nil {
		return nil
	}
	if bin.ItemInfoList == nil {
		bin.ItemInfoList = NewItemList(s)
	}
	return bin.ItemInfoList
}

func GetItemInfo(s *enter.Session, itemId int64) *sro.ItemInfo {
	bin := GetItemList(s)
	if bin == nil {
		return nil
	}
	return bin[itemId]
}

func AddItem(s *enter.Session, id int64, num int32) int64 {
	bin := GetItemBin(s)
	if bin == nil {
		return 0
	}
	if bin.ItemInfoList == nil {
		bin.ItemInfoList = NewItemList(s)
	}
	if info, ok := bin.ItemInfoList[id]; ok {
		info.StackCount += num
		return info.ServerId
	}
	info := &sro.ItemInfo{
		ServerId:   GetServerId(s),
		UniqueId:   id,
		StackCount: num,
	}
	bin.ItemInfoList[id] = info
	return info.ServerId
}
