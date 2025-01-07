package protocol

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/cmd"
	"github.com/gucooing/BaPs/mx/proto"
)

type NetworkProtocolResponse struct {
	Packet   string `json:"packet"`   // 包数据
	Protocol string `json:"protocol"` // 协议名称
}

// UnmarshalRequest 解码req数据包
func UnmarshalRequest(b []byte) (mx.Message, *mx.BasePacket, error) {
	base := new(mx.BasePacket)
	err := json.Unmarshal(b, base)
	if err != nil {
		return nil, nil, err
	}
	packet := cmd.Get().GetRequestPacketByCmdId(base.GetProtocol())
	if packet == nil {
		return nil, nil, errors.New(fmt.Sprintf("request unknown cmd id: %v", base.GetProtocol()))
	}
	err = json.Unmarshal(b, packet)
	if err != nil {
		return nil, nil, err
	}

	return packet, base, nil
}

// MarshalResponse 编码rsp数据包
func MarshalResponse(m mx.Message) (*NetworkProtocolResponse, error) {
	if m == nil {
		return nil, errors.New("message nil")
	}
	jsonData, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	v := &NetworkProtocolResponse{
		Packet:   string(jsonData),
		Protocol: proto.Protocol(m.GetProtocol()).String(),
	}

	return v, nil
}
