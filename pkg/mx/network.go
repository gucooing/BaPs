package mx

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/pkg/alg"
)

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
	return DeMx(bin)
}

func DeMx(bin []byte) ([]byte, error) {
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
