package proto

import (
	"github.com/gucooing/BaPs/pkg/mx"
)

type CafeCharacterDB struct {
	IsSummon         bool
	LastInteractTime mx.MxTime
	UniqueId         int64
	ServerId         int64
}
