package proto

type RaidUserDB struct {
	AccountId                   int64
	RepresentCharacterUniqueId  int64
	RepresentCharacterCostumeId int64
	Level                       int64
	Nickname                    string
	Tier                        int32
	Rank                        int64
	BestRankingPoint            int64
	BestRankingPointDetail      float64
	AccountAttachmentDB         *AccountAttachmentDB
}
