package alg

import (
	crr "crypto/rand"
	"encoding/base64"
	"math/rand"
)

func RandCode() int32 {
	return int32(rand.Intn(900000) + 100000)
}

func RandStr(length int) string {
	key := make([]byte, length)
	crr.Read(key)
	return base64.URLEncoding.EncodeToString(key)
}
