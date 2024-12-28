package alg

import (
	"strconv"

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

func S2I64(msg string) int64 {
	if msg == "" {
		return 0
	}
	ms, _ := strconv.ParseUint(msg, 10, 32)
	return int64(ms)
}
