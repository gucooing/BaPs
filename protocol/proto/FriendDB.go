package proto

import (
	"github.com/gucooing/BaPs/pkg/mx"
)

type FriendDB struct {
	AccountId                   int64
	Level                       int32
	Nickname                    string
	LastConnectTime             mx.MxTime
	RepresentCharacterUniqueId  int64
	RepresentCharacterCostumeId int64
	ComfortValue                int64
	FriendCount                 int64
	AttachmentDB                *AccountAttachmentDB
}
