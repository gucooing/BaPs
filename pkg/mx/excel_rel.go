//go:build rel
// +build rel

package mx

import (
	"fmt"
)

func GetMxToken(uid int64, len int) string {
	return fmt.Sprintf("%d%v", uid, len)
}

func DeExcelBytes(bin []byte) ([]byte, error) {
	return bin, nil
}
