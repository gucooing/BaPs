// Code auto-generated by gucooing & Hiro420 DO NOT EDIT.

package proto

import (
	"github.com/bytedance/sonic"
	"github.com/gucooing/BaPs/protocol/mx"
)

type CafeRankUpRequest struct {
	*RequestPacket

	AccountServerId  int64
	CafeDBId         int64
	ConsumeRequestDB *ConsumeRequestDB
}

func (x *CafeRankUpRequest) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}

func (x *CafeRankUpRequest) SetPacket(packet mx.Message) {
	if x == nil {
		return
	}
	x.RequestPacket = packet.(*RequestPacket)
}
