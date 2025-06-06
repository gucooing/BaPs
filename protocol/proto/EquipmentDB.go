// Code auto-generated by gucooing & Hiro420 DO NOT EDIT.

package proto

import (
	"github.com/bytedance/sonic"
	"github.com/gucooing/BaPs/protocol/mx"
)

type EquipmentDB struct {
	*ConsumableItemBaseDB

	Type                   ParcelType
	ParcelInfos            []*ParcelInfo
	Level                  int32
	Exp                    int64
	Tier                   int32
	BoundCharacterServerId int64
	CanConsume             bool
}

func (x *EquipmentDB) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}

func (x *EquipmentDB) SetPacket(packet mx.Message) {
	if x == nil {
		return
	}
	x.ConsumableItemBaseDB = packet.(*ConsumableItemBaseDB)
}
