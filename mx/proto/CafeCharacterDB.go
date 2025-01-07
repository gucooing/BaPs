package proto

import (
	"time"
)

type CafeCharacterDB struct {
	IsSummon         bool
	LastInteractTime time.Time
}
