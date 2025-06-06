// Code auto-generated by gucooing & Hiro420 DO NOT EDIT.
package proto

type BattleTypes int32

const (
	BattleTypes_None                BattleTypes = 0
	BattleTypes_Adventure           BattleTypes = 1
	BattleTypes_ScenarioMode        BattleTypes = 2
	BattleTypes_WeekDungeonChaserA  BattleTypes = 4
	BattleTypes_WeekDungeonBlood    BattleTypes = 8
	BattleTypes_WeekDungeonChaserB  BattleTypes = 16
	BattleTypes_WeekDungeonChaserC  BattleTypes = 32
	BattleTypes_WeekDungeonFindGift BattleTypes = 64
	BattleTypes_EventContent        BattleTypes = 128
	BattleTypes_TutorialAdventure   BattleTypes = 256
	BattleTypes_Profiling           BattleTypes = 512
	BattleTypes_SingleRaid          BattleTypes = 2048
	BattleTypes_MultiRaid           BattleTypes = 4096
	BattleTypes_PracticeRaid        BattleTypes = 8192
	BattleTypes_EliminateRaid       BattleTypes = 16384
	BattleTypes_MultiFloorRaid      BattleTypes = 32768
	BattleTypes_MinigameDefense     BattleTypes = 1048576
	BattleTypes_Arena               BattleTypes = 2097152
	BattleTypes_TimeAttack          BattleTypes = 8388608
	BattleTypes_SchoolDungeonA      BattleTypes = 33554432
	BattleTypes_SchoolDungeonB      BattleTypes = 67108864
	BattleTypes_SchoolDungeonC      BattleTypes = 134217728
	BattleTypes_WorldRaid           BattleTypes = 268435456
	BattleTypes_Conquest            BattleTypes = 536870912
	BattleTypes_FieldStory          BattleTypes = 1073741824
	BattleTypes_FieldContent        BattleTypes = -2147483648
	BattleTypes_PvE                 BattleTypes = -301988865
	BattleTypes_WeekDungeon         BattleTypes = 124
	BattleTypes_SchoolDungeon       BattleTypes = 234881024
	BattleTypes_Raid                BattleTypes = 30720
	BattleTypes_PvP                 BattleTypes = 2097152
	BattleTypes_All                 BattleTypes = -1
)

var (
	BattleTypes_name = map[int32]string{
		0:           "None",
		1:           "Adventure",
		2:           "ScenarioMode",
		4:           "WeekDungeonChaserA",
		8:           "WeekDungeonBlood",
		16:          "WeekDungeonChaserB",
		32:          "WeekDungeonChaserC",
		64:          "WeekDungeonFindGift",
		128:         "EventContent",
		256:         "TutorialAdventure",
		512:         "Profiling",
		2048:        "SingleRaid",
		4096:        "MultiRaid",
		8192:        "PracticeRaid",
		16384:       "EliminateRaid",
		32768:       "MultiFloorRaid",
		1048576:     "MinigameDefense",
		2097152:     "Arena",
		8388608:     "TimeAttack",
		33554432:    "SchoolDungeonA",
		67108864:    "SchoolDungeonB",
		134217728:   "SchoolDungeonC",
		268435456:   "WorldRaid",
		536870912:   "Conquest",
		1073741824:  "FieldStory",
		-2147483648: "FieldContent",
		-301988865:  "PvE",
		124:         "WeekDungeon",
		234881024:   "SchoolDungeon",
		30720:       "Raid",
		-1:          "All",
	}
	BattleTypes_value = map[string]int32{
		"None":                0,
		"Adventure":           1,
		"ScenarioMode":        2,
		"WeekDungeonChaserA":  4,
		"WeekDungeonBlood":    8,
		"WeekDungeonChaserB":  16,
		"WeekDungeonChaserC":  32,
		"WeekDungeonFindGift": 64,
		"EventContent":        128,
		"TutorialAdventure":   256,
		"Profiling":           512,
		"SingleRaid":          2048,
		"MultiRaid":           4096,
		"PracticeRaid":        8192,
		"EliminateRaid":       16384,
		"MultiFloorRaid":      32768,
		"MinigameDefense":     1048576,
		"Arena":               2097152,
		"TimeAttack":          8388608,
		"SchoolDungeonA":      33554432,
		"SchoolDungeonB":      67108864,
		"SchoolDungeonC":      134217728,
		"WorldRaid":           268435456,
		"Conquest":            536870912,
		"FieldStory":          1073741824,
		"FieldContent":        -2147483648,
		"PvE":                 -301988865,
		"WeekDungeon":         124,
		"SchoolDungeon":       234881024,
		"Raid":                30720,
		"PvP":                 2097152,
		"All":                 -1,
	}
)

func (x BattleTypes) String() string {
	return BattleTypes_name[int32(x)]
}

func (x BattleTypes) Value(sr string) BattleTypes {
	return BattleTypes(BattleTypes_value[sr])
}
