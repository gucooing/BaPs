package proto

type ParcelCost struct {
	ParcelInfos      []*ParcelInfo
	Currency         CurrencyTransaction
	EquipmentDBs     []*EquipmentDB
	ItemDBs          []*ItemDB
	FurnitureDBs     []*FurnitureDB
	ConsumeCondition ConsumeCondition
}
