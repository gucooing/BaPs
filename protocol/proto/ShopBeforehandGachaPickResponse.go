// Code auto-generated by gucooing & Hiro420 DO NOT EDIT.

package proto

import (
	"github.com/bytedance/sonic"
	"github.com/gucooing/BaPs/protocol/mx"
)

type ShopBeforehandGachaPickResponse struct {
	*ResponsePacket

	GachaResults  []*GachaResult
	AcquiredItems []*ItemDB
}

func (x *ShopBeforehandGachaPickResponse) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}

func (x *ShopBeforehandGachaPickResponse) SetPacket(packet mx.Message) {
	if x == nil {
		return
	}
	x.ResponsePacket = packet.(*ResponsePacket)
}
