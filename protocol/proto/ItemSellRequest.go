// Code auto-generated by gucooing & Hiro420 DO NOT EDIT.

package proto

import (
	"github.com/bytedance/sonic"
	"github.com/gucooing/BaPs/protocol/mx"
)

type ItemSellRequest struct {
	*RequestPacket

	TargetServerIds []int64
}

func (x *ItemSellRequest) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}

func (x *ItemSellRequest) SetPacket(packet mx.Message) {
	if x == nil {
		return
	}
	x.RequestPacket = packet.(*RequestPacket)
}
