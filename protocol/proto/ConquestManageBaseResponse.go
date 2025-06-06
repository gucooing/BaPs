// Code auto-generated by gucooing & Hiro420 DO NOT EDIT.

package proto

import (
	"github.com/bytedance/sonic"
	"github.com/gucooing/BaPs/protocol/mx"
)

type ConquestManageBaseResponse struct {
	*ResponsePacket

	ClearParcels                 [][]*ParcelInfo
	ConquerBonusParcels          [][]*ParcelInfo
	BonusParcels                 []*ParcelInfo
	ParcelResultDB               *ParcelResultDB
	ConquestInfoDB               *ConquestInfoDB
	ConquestEventObjectDBWrapper []*ConquestEventObjectDB
	DisplayInfos                 []*ConquestDisplayInfo
}

func (x *ConquestManageBaseResponse) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}

func (x *ConquestManageBaseResponse) SetPacket(packet mx.Message) {
	if x == nil {
		return
	}
	x.ResponsePacket = packet.(*ResponsePacket)
}
