package protocol

type ProtoMessage interface {
	ProtoReflect() Message    // 顶层反射
	GetProtocolKey() uint16   // 获取protocol id
	GetProtocolValue() string // 获取protocol name
}

type Message = ProtoMessage
