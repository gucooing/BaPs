// Code auto-generated by gucooing & Hiro420 DO NOT EDIT.

package proto

import (
	"github.com/bytedance/sonic"
	"github.com/gucooing/BaPs/protocol/mx"
)

type CafeRelocateFurnitureRequest struct {
	*RequestPacket

	CafeDBId    int64
	FurnitureDB *FurnitureDB
}

func (x *CafeRelocateFurnitureRequest) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}

func (x *CafeRelocateFurnitureRequest) SetPacket(packet mx.Message) {
	if x == nil {
		return
	}
	x.RequestPacket = packet.(*RequestPacket)
}
