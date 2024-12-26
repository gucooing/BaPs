package protocol

import (
	"encoding/json"
	"errors"
)

type NetworkProtocol struct {
	Packet   string // 包数据
	Protocol string // 协议名称
}

func Unmarshal(b []byte, m Message) (string, error) {
	if m == nil {
		return "", errors.New("message nil")
	}
	network := new(NetworkProtocol)
	err := json.Unmarshal(b, network)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal([]byte(network.Packet), m)
	if err != nil {
		return "", err
	}

	return network.Protocol, nil
}

func Marshal(m Message) ([]byte, error) {
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

	return json.Marshal(v)
}
