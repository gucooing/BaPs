package alg

import (
	"strconv"
	"time"

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
	ms, _ := strconv.ParseInt(msg, 10, 32)
	return ms
}

func S2I32(msg string) int32 {
	if msg == "" {
		return 0
	}
	ms, _ := strconv.ParseInt(msg, 10, 32)
	return int32(ms)
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

func MainInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func MinInt64(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}

func MaxInt32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func MaxBool(a, b bool) bool {
	if a {
		return a
	}
	return b
}

func GetDay4() time.Time {
	currentTime := time.Now()
	nextExecution := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 4, 0, 0, 0, currentTime.Location())
	return nextExecution
}

func GetEveryDay4() time.Duration {
	currentTime := time.Now()
	nextExecution := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 4, 0, 0, 0, currentTime.Location())
	if currentTime.Hour() >= 4 {
		nextExecution = nextExecution.AddDate(0, 0, 1)
	}
	return nextExecution.Sub(currentTime)
}

func GetTimeHour4() time.Time {
	currentTime := time.Now()
	hour := currentTime.Hour()
	if hour < 4 {
		previousTwo := currentTime.Add(-24 * time.Hour)
		return time.Date(previousTwo.Year(), previousTwo.Month(), previousTwo.Day(), 16, 0, 0, 0, previousTwo.Location())
	} else if hour < 16 {
		return time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 4, 0, 0, 0, currentTime.Location())
	} else {
		return time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 16, 0, 0, 0, currentTime.Location())
	}
}
