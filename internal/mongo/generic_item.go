package mongo

type GenericItem struct {
	Id              string      `json:"id" bson:"_id"`
	ItemId          string      `json:"item_id" bson:"itemId"`
	Rarity          string      `json:"rarity" bson:"rarity"`
	Reforge         *string     `json:"reforge" bson:"reforge"`
	Enchantments    interface{} `json:"enchantments" bson:"enchantments"`
	ExtraAttributes interface{} `json:"extra_attributes" bson:"extraAttributes"`
}
