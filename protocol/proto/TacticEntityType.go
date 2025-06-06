// Code auto-generated by gucooing & Hiro420 DO NOT EDIT.
package proto

type TacticEntityType int32

const (
	TacticEntityType_None                   TacticEntityType = 0
	TacticEntityType_Student                TacticEntityType = 1
	TacticEntityType_Minion                 TacticEntityType = 2
	TacticEntityType_Elite                  TacticEntityType = 4
	TacticEntityType_Champion               TacticEntityType = 8
	TacticEntityType_Boss                   TacticEntityType = 16
	TacticEntityType_Obstacle               TacticEntityType = 32
	TacticEntityType_Servant                TacticEntityType = 64
	TacticEntityType_Vehicle                TacticEntityType = 128
	TacticEntityType_Summoned               TacticEntityType = 256
	TacticEntityType_Hallucination          TacticEntityType = 512
	TacticEntityType_DestructibleProjectile TacticEntityType = 1024
)

var (
	TacticEntityType_name = map[int32]string{
		0:    "None",
		1:    "Student",
		2:    "Minion",
		4:    "Elite",
		8:    "Champion",
		16:   "Boss",
		32:   "Obstacle",
		64:   "Servant",
		128:  "Vehicle",
		256:  "Summoned",
		512:  "Hallucination",
		1024: "DestructibleProjectile",
	}
	TacticEntityType_value = map[string]int32{
		"None":                   0,
		"Student":                1,
		"Minion":                 2,
		"Elite":                  4,
		"Champion":               8,
		"Boss":                   16,
		"Obstacle":               32,
		"Servant":                64,
		"Vehicle":                128,
		"Summoned":               256,
		"Hallucination":          512,
		"DestructibleProjectile": 1024,
	}
)

func (x TacticEntityType) String() string {
	return TacticEntityType_name[int32(x)]
}

func (x TacticEntityType) Value(sr string) TacticEntityType {
	return TacticEntityType(TacticEntityType_value[sr])
}
