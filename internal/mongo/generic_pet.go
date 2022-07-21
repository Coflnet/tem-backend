package mongo

type GenericPet struct {
	Id       string  `json:"id" bson:"_id"`
	Candy    int     `json:"candy" bson:"candy"`
	HeldItem string  `json:"held_item" bson:"heldItem"`
	Level    int     `json:"level" bson:"level"`
	Name     string  `json:"name" bson:"name"`
	Rarity   string  `json:"rarity" bson:"rarity"`
	Skin     *string `json:"skin" bson:"skin"`
}
