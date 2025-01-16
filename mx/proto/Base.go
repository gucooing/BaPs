package proto

import (
	"github.com/gucooing/BaPs/mx"
)

// RequestPacket 来自客户端的基础数据包
type RequestPacket struct {
	*mx.BasePacket
	Resendable bool
	Hash       uint64
	IsTest     bool
}

// ResponsePacket 来自服务端的基础数据包
type ResponsePacket struct {
	*mx.BasePacket
	MissionProgressDBs         []*MissionProgressDB
	EventMissionProgressDBDict map[uint64][]*MissionProgressDB
	StaticOpenConditions       map[string]int32 //   map[OpenConditionContent]OpenConditionLockReason
}
