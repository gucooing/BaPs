// Code auto-generated by gucooing & Hiro420 DO NOT EDIT.

package proto

import (
	"github.com/bytedance/sonic"
	"github.com/gucooing/BaPs/protocol/mx"
)

type CharacterGearListResponse struct {
	*ResponsePacket

	GearDBs []*GearDB
}

func (x *CharacterGearListResponse) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}

func (x *CharacterGearListResponse) SetPacket(packet mx.Message) {
	if x == nil {
		return
	}
	x.ResponsePacket = packet.(*ResponsePacket)
}
