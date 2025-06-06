// Code auto-generated by gucooing & Hiro420 DO NOT EDIT.

package proto

import (
	"github.com/bytedance/sonic"
	"github.com/gucooing/BaPs/protocol/mx"
)

type WeekDungeonListResponse struct {
	*ResponsePacket

	AdditionalStageIdList         []int64
	WeekDungeonStageHistoryDBList []*WeekDungeonStageHistoryDB
}

func (x *WeekDungeonListResponse) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}

func (x *WeekDungeonListResponse) SetPacket(packet mx.Message) {
	if x == nil {
		return
	}
	x.ResponsePacket = packet.(*ResponsePacket)
}
