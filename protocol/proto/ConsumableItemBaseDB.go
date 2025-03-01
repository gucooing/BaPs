package proto

type ConsumableItemBaseDB struct {
	Key        *ParcelKeyPair
	ServerId   int64
	UniqueId   int64
	StackCount int64
}
