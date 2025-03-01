package proto

import (
	"time"
)

type MissionProgressDB struct {
	AccountServerId    int64
	ServerId           int64
	MissionUniqueId    int64
	Complete           bool
	StartTime          time.Time
	ProgressParameters map[int64]int64
}
