// Code auto-generated by gucooing & Hiro420 DO NOT EDIT.

package proto

import (
	"github.com/bytedance/sonic"
	"github.com/gucooing/BaPs/protocol/mx"
)

type ProofTokenSubmitRequest struct {
	*RequestPacket

	Answer int64
}

func (x *ProofTokenSubmitRequest) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}

func (x *ProofTokenSubmitRequest) SetPacket(packet mx.Message) {
	if x == nil {
		return
	}
	x.RequestPacket = packet.(*RequestPacket)
}
