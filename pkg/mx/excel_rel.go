//go:build rel
// +build rel

package mx

import (
	"fmt"

	"github.com/gucooing/BaPs/pkg/logger"
)

const (
	ExcelVersion = 1
)

func GetMxToken(uid int64, len int) string {
	return fmt.Sprintf("%d%v", uid, len)
}

func DeExcelBytes(bin []byte) ([]byte, error) {
	return bin, nil
}

func LoadExcelJson[T any](path string, table *[]*T) {
	logger.Error("文件:%s 读取失败,请自行补充相关逻辑s", path)
}
