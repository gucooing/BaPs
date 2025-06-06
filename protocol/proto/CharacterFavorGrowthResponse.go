// Code auto-generated by gucooing & Hiro420 DO NOT EDIT.

package proto

import (
	"github.com/bytedance/sonic"
	"github.com/gucooing/BaPs/protocol/mx"
)

type CharacterFavorGrowthResponse struct {
	*ResponsePacket

	CharacterDB                  *CharacterDB
	ConsumeStackableItemDBResult []*ItemDB
	ParcelResultDB               *ParcelResultDB
}

func (x *CharacterFavorGrowthResponse) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}

func (x *CharacterFavorGrowthResponse) SetPacket(packet mx.Message) {
	if x == nil {
		return
	}
	x.ResponsePacket = packet.(*ResponsePacket)
}
