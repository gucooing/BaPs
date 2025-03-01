package proto

import (
	"encoding/json"
)

type ProtoMessage interface {
	// String 序列化成json 用于debug
	String() string
	// ProtoReflect 顶层反射
	ProtoReflect() Message
	// GetProtocol 获取protocol id
	GetProtocol() int32
	// SetSessionKey 设置基础数据
	SetSessionKey(*BasePacket)
}

type Message = ProtoMessage

// BasePacket 基础数据包
type BasePacket struct {
	SessionKey *SessionKey
	Protocol   int32
	AccountId  int64
	// ResponsePacket 来自服务端的基础数据包
	ServerTimeTicks            int64
	ServerNotification         int32
	ErrorCode                  int32
	MissionProgressDBs         []*MissionProgressDB
	EventMissionProgressDBDict map[uint64][]*MissionProgressDB
	StaticOpenConditions       map[string]int32 //   map[OpenConditionContent]OpenConditionLockReason
	// RequestPacket 来自客户端的基础数据包
	ClientUpTime int
	Resendable   bool
	Hash         uint64
	IsTest       bool
}

type SessionKey struct {
	AccountServerId int64
	MxToken         string
}

func (x *BasePacket) GetSessionKey() *SessionKey {
	if x == nil {
		return nil
	}
	return x.SessionKey
}

func (x *BasePacket) GetProtocol() int32 {
	if x == nil {
		return 0
	}
	return x.Protocol
}

func (x *SessionKey) String() string {
	if x == nil {
		return ""
	}
	bin, _ := json.Marshal(x)
	return string(bin)
}

// RequestPacket 来自客户端的基础数据包
type RequestPacket struct {
	*BasePacket
}

// ResponsePacket 来自服务端的基础数据包
type ResponsePacket struct {
	*BasePacket
}
