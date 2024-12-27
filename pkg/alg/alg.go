package alg

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"

	"github.com/gin-gonic/gin"
)

func Xor(data []byte, key []byte) {
	for i := 0; i < len(data); i++ {
		data[i] ^= key[i%len(key)]
	}
}

func CheckGateWay(c *gin.Context) bool {
	if c.GetHeader("user-agent") != "BestHTTP/2 v2.4.0" ||
		c.GetHeader("accept-encoding") != "gzip" {
		return false
	}
	return true
}

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
	Xor(bin, []byte{0xD9})
	z, err := gzip.NewReader(bytes.NewReader(bin[12:]))
	if err != nil {
		return nil, err
	}
	defer z.Close()
	p, err := io.ReadAll(z)
	return p, err
}
