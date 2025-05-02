package mx

type ProtoMessage interface {
	String() string
	SetPacket(packet Message)
}

type Message = ProtoMessage
