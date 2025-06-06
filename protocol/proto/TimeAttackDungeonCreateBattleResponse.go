// Code auto-generated by gucooing & Hiro420 DO NOT EDIT.

package proto

import (
	"github.com/bytedance/sonic"
	"github.com/gucooing/BaPs/protocol/mx"
)

type TimeAttackDungeonCreateBattleResponse struct {
	*ResponsePacket

	RoomDB         *TimeAttackDungeonRoomDB
	ParcelResultDB *ParcelResultDB
}

func (x *TimeAttackDungeonCreateBattleResponse) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}

func (x *TimeAttackDungeonCreateBattleResponse) SetPacket(packet mx.Message) {
	if x == nil {
		return
	}
	x.ResponsePacket = packet.(*ResponsePacket)
}
