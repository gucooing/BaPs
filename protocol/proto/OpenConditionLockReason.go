package proto

type OpenConditionLockReason int32

const (
	OpenConditionLockReason_None              = 0
	OpenConditionLockReason_Level             = 1
	OpenConditionLockReason_StageClear        = 2
	OpenConditionLockReason_Time              = 4
	OpenConditionLockReason_Day               = 8
	OpenConditionLockReason_CafeRank          = 16
	OpenConditionLockReason_ScenarioModeClear = 32
	OpenConditionLockReason_CafeOpen          = 64
)
