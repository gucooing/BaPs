// Code generated by gucooing DO NOT EDIT.

package proto

import (
	"encoding/json"

	"github.com/gucooing/BaPs/pkg/mx"
)

type NotificationEventContentReddotRequest struct{
	message ProtoMessage
	RequestPacket

    
}

func (x *NotificationEventContentReddotRequest) String() string {
	j, _ := json.Marshal(x)
	return string(j)
}

func (x *NotificationEventContentReddotRequest) ProtoReflect() Message {
	return x
}

func (x *NotificationEventContentReddotRequest) GetProtocol() int32 {
	return mx.Protocol_Notification_EventContentReddotCheck
}

func (x *NotificationEventContentReddotRequest) SetSessionKey(base *BasePacket) {
	if x == nil {
		return
	}
	x.BasePacket = base
}

