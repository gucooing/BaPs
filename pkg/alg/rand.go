package alg

import (
	crr "crypto/rand"
	"encoding/base64"
	"math/rand"
)

func RandCode() int32 {
	return rand.Int31n(900000) + 100000
}

func RandCodeInt64() int64 {
	return rand.Int63n(90000000) + 100000
}

func RandStr(length int) string {
	key := make([]byte, length)
	crr.Read(key)
	return base64.URLEncoding.EncodeToString(key)
}
