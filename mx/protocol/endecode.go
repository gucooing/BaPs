package protocol

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io/ioutil"

	"github.com/gucooing/BaPs/pkg/alg"
)

// Decode 解密客户端传入的数据
func Decode(bin []byte) (string, error) {
	if len(bin) <= 12 {
		return "", errors.New("binary too short")
	}
	alg.Xor(bin, []byte{0xD9})
	z, err := gzip.NewReader(bytes.NewReader(bin[12:]))
	if err != nil {
		return "", err
	}
	defer z.Close()
	p, err := ioutil.ReadAll(z)
	return string(p), err
}
