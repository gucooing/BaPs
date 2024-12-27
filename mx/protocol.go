package mx

type ProtoMessage interface {
	String() string
	ProtoReflect() Message    // 顶层反射
	GetProtocolKey() int32    // 获取protocol id
	GetProtocolValue() string // 获取protocol name
}

type Message = ProtoMessage
