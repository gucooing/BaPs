package proto

import (
	"time"
)

type CafeCharacterDB struct {
	IsSummon         bool
	LastInteractTime time.Time
	UniqueId         int64
	ServerId         int64
}
