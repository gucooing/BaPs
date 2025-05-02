package cmd

import (
	"github.com/gucooing/BaPs/protocol/mx"
	"github.com/gucooing/BaPs/protocol/proto"
	"reflect"
	"sync"

	"github.com/gucooing/BaPs/pkg/logger"
)

var sharedCmdProtoMap *CmdProtoMap
var cmdProtoMapOnce sync.Once

type CmdProtoMap struct {
	protoObjCmdIdMap map[reflect.Type]proto.Protocol

	requestPacketMap  map[proto.Protocol]reflect.Type
	responsePacketMap map[proto.Protocol]reflect.Type
}

func Get() *CmdProtoMap {
	cmdProtoMapOnce.Do(func() {
		sharedCmdProtoMap = NewCmdProtoMap()
	})
	return sharedCmdProtoMap
}

func NewCmdProtoMap() (r *CmdProtoMap) {
	r = new(CmdProtoMap)
	r.protoObjCmdIdMap = make(map[reflect.Type]proto.Protocol)

	r.requestPacketMap = make(map[proto.Protocol]reflect.Type)
	r.responsePacketMap = make(map[proto.Protocol]reflect.Type)
	r.registerAllMessage()

	return r
}

func (c *CmdProtoMap) regMsg(cmdId proto.Protocol, protoObjNewFunc func() any, isRequest bool) {
	protoObj := protoObjNewFunc().(mx.Message)
	refType := reflect.TypeOf(protoObj)
	// protoObj -> cmdId
	c.protoObjCmdIdMap[refType] = cmdId

	if isRequest {
		c.requestPacketMap[cmdId] = refType
	} else {
		c.responsePacketMap[cmdId] = refType
	}
}

// 反射方法

func (c *CmdProtoMap) GetRequestPacketByCmdId(cmdId proto.Protocol) mx.Message {
	refType, exist := c.requestPacketMap[cmdId]
	if !exist {
		return nil
	}
	protoObjInst := reflect.New(refType.Elem())
	protoObj := protoObjInst.Interface().(mx.Message)
	return protoObj
}

func (c *CmdProtoMap) GetResponsePacketByCmdId(cmdId proto.Protocol) mx.Message {
	refType, exist := c.responsePacketMap[cmdId]
	if !exist {
		return nil
	}
	protoObjInst := reflect.New(refType.Elem())
	protoObj := protoObjInst.Interface().(mx.Message)
	return protoObj
}

func (c *CmdProtoMap) GetCmdIdByProtoObj(protoObj mx.Message) proto.Protocol {
	cmdId, exist := c.protoObjCmdIdMap[reflect.TypeOf(protoObj)]
	if !exist {
		logger.Debug("unknown proto object: %v\n", protoObj)
		return 0
	}
	return cmdId
}

func (c *CmdProtoMap) GetCmdIdByCmdName(cmdName string) proto.Protocol {
	cmdId, exist := mx.Protocol_value[cmdName]
	if !exist {
		logger.Debug("unknown cmd name: %v\n", cmdName)
		return 0
	}
	return proto.Protocol(cmdId)
}
