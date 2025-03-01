package proto

type ConsumeCondition int32

const (
	ConsumeCondition_And = 0
	ConsumeCondition_Or  = 1
)

var (
	ConsumeCondition_name = map[int32]string{
		0: "And",
		1: "Or",
	}
	ConsumeCondition_value = map[string]int32{
		"And": 0,
		"Or":  1,
	}
)

func (x ConsumeCondition) String(v int32) string {
	return ConsumeCondition_name[v]
}
