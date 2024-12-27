package protocol

import (
	"encoding/json"
	"errors"

	"github.com/gucooing/BaPs/mx"
	"github.com/gucooing/BaPs/mx/cmd"
)

type NetworkProtocol struct {
	Packet   string `json:"packet"`   // 包数据
	Protocol string `json:"protocol"` // 协议名称
}

// Unmarshal 解码
func Unmarshal(b []byte) (mx.Message, error) {
	network := new(NetworkProtocol)
	err := json.Unmarshal(b, network)
	if err != nil {
		return nil, err
	}
	packet := cmd.Get().GetRequestPacketByCmdId(network.Protocol)
	err = json.Unmarshal([]byte(network.Packet), packet)
	if err != nil {
		return nil, err
	}

	return packet, nil
}

// Marshal 编码
func Marshal(m mx.Message) (*NetworkProtocol, error) {
	if m == nil {
		return nil, errors.New("message nil")
	}
	jsonData, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	v := &NetworkProtocol{
		Packet:   string(jsonData),
		Protocol: m.GetProtocolValue(),
	}

	return v, nil
}
