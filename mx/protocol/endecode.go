package protocol

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"

	"github.com/gucooing/BaPs/pkg/alg"
)

// Decode 解密客户端传入的数据
func Decode(bin []byte) ([]byte, error) {
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
