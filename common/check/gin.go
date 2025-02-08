package check

import (
	"math"
	"sync/atomic"
	"time"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/pkg/logger"
)

var TPS int64
var RT int64

func GinNetInfo() {
	ticker := time.NewTicker(time.Second * 60)
	for {
		<-ticker.C
		tps := atomic.LoadInt64(&TPS)
		rt := float64(atomic.LoadInt64(&RT)) / (float64(tps) * 1000 * 1000)
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
