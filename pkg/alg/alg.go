package alg

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/config"
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

func AutoGucooingApi() gin.HandlerFunc {
	if config.GetGucooingApiKey() == "" {
		return func(c *gin.Context) {}
	} else {
		return func(c *gin.Context) {
			if c.GetHeader("Authorization-Gucooing") != config.GetGucooingApiKey() {
				c.String(401, "Unauthorized")
				c.Abort()
			}
		}
	}
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
