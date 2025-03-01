package proto

type SingleRaidLobbyInfoDB struct {
	*RaidLobbyInfoDB
	ClearDifficulty []Difficulty
}
