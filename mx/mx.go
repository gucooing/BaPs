package mx

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/pkg/alg"
)

type ProtoMessage interface {
	String() string
	ProtoReflect() Message     // 顶层反射
	GetProtocol() int32        // 获取protocol
	SetSessionKey(*BasePacket) // 设置 base
}

type Message = ProtoMessage

func GetFormMx(c *gin.Context) ([]byte, error) {
	file, err := c.FormFile("mx")
	if err != nil {
		return nil, err
	}
	fileContent, err := file.Open()
	if err != nil {
		return nil, err
	}
	bin := make([]byte, file.Size)
	_, err = fileContent.Read(bin)
	if err != nil {
		return nil, err
	}
	// 下面是解密
	if len(bin) <= 12 {
		return nil, errors.New("binary too short")
	}
	alg.Xor(bin, []byte{0xD9})
	z, err := gzip.NewReader(bytes.NewReader(bin[12:]))
	if err != nil {
		return nil, err
	}
	defer z.Close()
	p, err := io.ReadAll(z)
	return p, err
}

// BasePacket 基础数据包
type BasePacket struct {
	SessionKey *SessionKey
	Protocol   int32
	AccountId  int64
	// s
	ServerTimeTicks    int64
	ServerNotification int32 // proto.ServerNotificationFlag
	// c
	ClientUpTime int
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
