package proto

type GachaResult struct {
	CharacterId int64
	Character   *CharacterDB
	Stone       *ItemDB
}
