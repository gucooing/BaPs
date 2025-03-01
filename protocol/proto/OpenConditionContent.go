package proto

type OpenConditionContent string

const (
	OpenConditionContent_Shop                              = 0
	OpenConditionContent_Gacha                             = 1
	OpenConditionContent_LobbyIllust                       = 2
	OpenConditionContent_Raid                              = 3
	OpenConditionContent_Cafe                              = 4
	OpenConditionContent_Unit_Growth_Skill                 = 5
	OpenConditionContent_Unit_Growth_LevelUp               = 6
	OpenConditionContent_Unit_Growth_Transcendence         = 7
	OpenConditionContent_Arena                             = 8
	OpenConditionContent_Academy                           = 9
	OpenConditionContent_Equip                             = 10
	OpenConditionContent_Item                              = 11
	OpenConditionContent_Favor                             = 12
	OpenConditionContent_Prologue                          = 13
	OpenConditionContent_Mission                           = 14
	OpenConditionContent_WeekDungeon_Chase                 = 15
	OpenConditionContent___Deprecated_WeekDungeon_FindGift = 16
	OpenConditionContent___Deprecated_WeekDungeon_Blood    = 17
	OpenConditionContent_Story_Sub                         = 18
	OpenConditionContent_Story_Replay                      = 19
	OpenConditionContent_WeekDungeon                       = 20
	OpenConditionContent_None                              = 21
	OpenConditionContent_Shop_Gem                          = 22
	OpenConditionContent_Craft                             = 23
	OpenConditionContent_Student                           = 24
	OpenConditionContent_GuideMission                      = 25
	OpenConditionContent_Clan                              = 26
	OpenConditionContent_Echelon                           = 27
	OpenConditionContent_Campaign                          = 28
	OpenConditionContent_EventContent                      = 29
	OpenConditionContent_Guild                             = 30
	OpenConditionContent_EventStage_1                      = 31
	OpenConditionContent_EventStage_2                      = 32
	OpenConditionContent_Talk                              = 33
	OpenConditionContent_Billing                           = 34
	OpenConditionContent_Schedule                          = 35
	OpenConditionContent_Story                             = 36
	OpenConditionContent_Tactic_Speed                      = 37
	OpenConditionContent_Cafe_Invite                       = 38
	OpenConditionContent_EventMiniGame_1                   = 39
	OpenConditionContent_SchoolDungeon                     = 40
	OpenConditionContent_TimeAttackDungeon                 = 41
	OpenConditionContent_ShiftingCraft                     = 42
	OpenConditionContent_WorldRaid                         = 43
	OpenConditionContent_Tactic_Skip                       = 44
	OpenConditionContent_Mulligan                          = 45
	OpenConditionContent_EventPermanent                    = 46
	OpenConditionContent_Main_L_1_2                        = 47
	OpenConditionContent_Main_L_1_3                        = 48
	OpenConditionContent_Main_L_1_4                        = 49
	OpenConditionContent_EliminateRaid                     = 50
	OpenConditionContent_Cafe_2                            = 51
	OpenConditionContent_Cafe_Invite_2                     = 52
	OpenConditionContent_MultiFloorRaid                    = 53
	OpenConditionContent_StrategySkip                      = 54
	OpenConditionContent_MinigameDreamMaker                = 55
	OpenConditionContent_Temp                              = 56
)

var (
	OpenConditionContent_name = map[int32]string{
		0:  "Shop",
		1:  "Gacha",
		2:  "LobbyIllust",
		3:  "Raid",
		4:  "Cafe",
		5:  "Unit_Growth_Skill",
		6:  "Unit_Growth_LevelUp",
		7:  "Unit_Growth_Transcendence",
		8:  "Arena",
		9:  "Academy",
		10: "Equip",
		11: "Item",
		12: "Favor",
		13: "Prologue",
		14: "Mission",
		15: "WeekDungeon_Chase",
		16: "__Deprecated_WeekDungeon_FindGift",
		17: "__Deprecated_WeekDungeon_Blood",
		18: "Story_Sub",
		19: "Story_Replay",
		20: "WeekDungeon",
		21: "None",
		22: "Shop_Gem",
		23: "Craft",
		24: "Student",
		25: "GuideMission",
		26: "Clan",
		27: "Echelon",
		28: "Campaign",
		29: "EventContent",
		30: "Guild",
		31: "EventStage_1",
		32: "EventStage_2",
		33: "Talk",
		34: "Billing",
		35: "Schedule",
		36: "Story",
		37: "Tactic_Speed",
		38: "Cafe_Invite",
		39: "EventMiniGame_1",
		40: "SchoolDungeon",
		41: "TimeAttackDungeon",
		42: "ShiftingCraft",
		43: "WorldRaid",
		44: "Tactic_Skip",
		45: "Mulligan",
		46: "EventPermanent",
		47: "Main_L_1_2",
		48: "Main_L_1_3",
		49: "Main_L_1_4",
		50: "EliminateRaid",
		51: "Cafe_2",
		52: "Cafe_Invite_2",
		53: "MultiFloorRaid",
		54: "StrategySkip",
		55: "MinigameDreamMaker",
		56: "Temp",
	}
	OpenConditionContent_value = map[string]int32{
		"Shop":                              0,
		"Gacha":                             1,
		"LobbyIllust":                       2,
		"Raid":                              3,
		"Cafe":                              4,
		"Unit_Growth_Skill":                 5,
		"Unit_Growth_LevelUp":               6,
		"Unit_Growth_Transcendence":         7,
		"Arena":                             8,
		"Academy":                           9,
		"Equip":                             10,
		"Item":                              11,
		"Favor":                             12,
		"Prologue":                          13,
		"Mission":                           14,
		"WeekDungeon_Chase":                 15,
		"__Deprecated_WeekDungeon_FindGift": 16,
		"__Deprecated_WeekDungeon_Blood":    17,
		"Story_Sub":                         18,
		"Story_Replay":                      19,
		"WeekDungeon":                       20,
		"None":                              21,
		"Shop_Gem":                          22,
		"Craft":                             23,
		"Student":                           24,
		"GuideMission":                      25,
		"Clan":                              26,
		"Echelon":                           27,
		"Campaign":                          28,
		"EventContent":                      29,
		"Guild":                             30,
		"EventStage_1":                      31,
		"EventStage_2":                      32,
		"Talk":                              33,
		"Billing":                           34,
		"Schedule":                          35,
		"Story":                             36,
		"Tactic_Speed":                      37,
		"Cafe_Invite":                       38,
		"EventMiniGame_1":                   39,
		"SchoolDungeon":                     40,
		"TimeAttackDungeon":                 41,
		"ShiftingCraft":                     42,
		"WorldRaid":                         43,
		"Tactic_Skip":                       44,
		"Mulligan":                          45,
		"EventPermanent":                    46,
		"Main_L_1_2":                        47,
		"Main_L_1_3":                        48,
		"Main_L_1_4":                        49,
		"EliminateRaid":                     50,
		"Cafe_2":                            51,
		"Cafe_Invite_2":                     52,
		"MultiFloorRaid":                    53,
		"StrategySkip":                      54,
		"MinigameDreamMaker":                55,
		"Temp":                              56,
	}
)

func OpenConditionContentString(x int32) string {
	return OpenConditionContent_name[x]
}

func OpenConditionContentValue(x string) int32 {
	return OpenConditionContent_value[x]
}
