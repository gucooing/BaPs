// Code auto-generated by gucooing & Hiro420 DO NOT EDIT.

package proto

import (
	"github.com/bytedance/sonic"
	"github.com/gucooing/BaPs/protocol/mx"
)

type EventContentDeployEchelonRequest struct {
	*RequestPacket

	EventContentId   int64
	StageUniqueId    int64
	DeployedEchelons []*HexaUnit
}

func (x *EventContentDeployEchelonRequest) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}

func (x *EventContentDeployEchelonRequest) SetPacket(packet mx.Message) {
	if x == nil {
		return
	}
	x.RequestPacket = packet.(*RequestPacket)
}
