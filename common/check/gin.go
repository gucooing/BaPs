package check

import (
	"math"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/pkg/logger"
)

var TPS int64
var RT int64
var OLDTPS int64
var OLDRT float64

func GinNetInfo() {
	ticker := time.NewTicker(time.Second * 60)
	for {
		<-ticker.C
		tps := atomic.LoadInt64(&TPS)
		rt := float64(atomic.LoadInt64(&RT)) / (float64(tps) * 1000 * 1000)
		OLDTPS = tps
		OLDRT = rt
		if tps == 0 && math.IsNaN(rt) {
			continue
		}
		logger.Info("SessionNum: %v", enter.GetSessionNum())
		logger.Info("TPS: %v", tps)
		logger.Info("RT: %.6f ms", rt)
		atomic.StoreInt64(&TPS, 0)
		atomic.StoreInt64(&RT, 0)
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}