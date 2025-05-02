package mx

import (
	"encoding/json"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/gucooing/BaPs/pkg/logger"
)

type MxTime time.Time

func (t MxTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format("2006-01-02T15:04:05"))
}

func (t *MxTime) UnmarshalJSON(data []byte) error {
	parsedTime, err := time.Parse("2006-01-02T15:04:05.9999999", strings.Trim(string(data), "\""))
	if err != nil {
		return err
	}
	*t = MxTime(parsedTime)
	return nil
}

// Add returns the time t+d.
func (t MxTime) Add(d time.Duration) MxTime {
	return MxTime(time.Time(t).Add(d))
}

// Unix returns the local Time corresponding to the given Unix time,
// sec seconds and nsec nanoseconds since January 1, 1970 UTC.
// It is valid to pass nsec outside the range [0, 999999999].
// Not all sec values have a corresponding time value. One such
// value is 1<<63-1 (the largest int64 value).
func Unix(sec int64, nsec int64) MxTime {
	return MxTime(time.Unix(sec, nsec))
}

func TimeData(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) MxTime {
	return MxTime(time.Date(year, month, day, hour, min, sec, nsec, loc))
}

func Now() MxTime {
	return MxTime(time.Now())
}

func (t MxTime) After(u MxTime) bool {
	return time.Time(t).After(time.Time(u))
}

func (t MxTime) Before(u MxTime) bool {
	return time.Time(t).Before(time.Time(u))
}

func SetTZ() {
	logger.OriginalLocal = time.Local
	//loc, err := time.LoadLocation("Asia/Tokyo")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//time.Local = loc
}
