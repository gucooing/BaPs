package cmd

import (
	"reflect"
	"sync"

	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/proto"
	"github.com/gucooing/BaPs/pkg/logger"
)

var sharedCmdProtoMap *CmdProtoMap
var cmdProtoMapOnce sync.Once

type CmdProtoMap struct {
	protoObjCmdIdMap map[reflect.Type]int32

	requestPacketMap  map[int32]reflect.Type
	responsePacketMap map[int32]reflect.Type
}

func Get() *CmdProtoMap {
	cmdProtoMapOnce.Do(func() {
		sharedCmdProtoMap = NewCmdProtoMap()
	})
	return sharedCmdProtoMap
}

func NewCmdProtoMap() (r *CmdProtoMap) {
	r = new(CmdProtoMap)
	r.protoObjCmdIdMap = make(map[reflect.Type]int32)

	r.requestPacketMap = make(map[int32]reflect.Type)
	r.responsePacketMap = make(map[int32]reflect.Type)
	r.registerAllMessage()

	return r
}

func (c *CmdProtoMap) regMsg(cmdId int32, protoObjNewFunc func() any, isRequest bool) {
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

func (c *CmdProtoMap) GetRequestPacketByCmdId(cmdIdAny any) mx.Message {
	var cmdId int32
	switch cmdIdAny.(type) {
	case int32:
		cmdId = cmdIdAny.(int32)
	case string:
		cmdId = c.GetCmdIdByCmdName(cmdIdAny.(string))
	default:
		return nil
	}
	refType, exist := c.requestPacketMap[cmdId]
	if !exist {
		return nil
	}
	protoObjInst := reflect.New(refType.Elem())
	protoObj := protoObjInst.Interface().(mx.Message)
	return protoObj
}

func (c *CmdProtoMap) GetResponsePacketByCmdId(cmdId int32) mx.Message {
	refType, exist := c.responsePacketMap[cmdId]
	if !exist {
		return nil
	}
	protoObjInst := reflect.New(refType.Elem())
	protoObj := protoObjInst.Interface().(mx.Message)
	return protoObj
}

func (c *CmdProtoMap) GetCmdIdByProtoObj(protoObj mx.Message) int32 {
	cmdId, exist := c.protoObjCmdIdMap[reflect.TypeOf(protoObj)]
	if !exist {
		logger.Debug("unknown proto object: %v\n", protoObj)
		return 0
	}
	return cmdId
}

func (c *CmdProtoMap) GetCmdNameByCmdId(cmdId int32) string {
	cmdName, exist := proto.Protocol_name[cmdId]
	if !exist {
		logger.Debug("unknown cmd id: %v\n", cmdId)
		return ""
	}
	return cmdName
}

func (c *CmdProtoMap) GetCmdIdByCmdName(cmdName string) int32 {
	cmdId, exist := proto.Protocol_value[cmdName]
	if !exist {
		logger.Debug("unknown cmd name: %v\n", cmdName)
		return 0
	}
	return cmdId
}
