package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type PlayerNotFound struct {
	PlayerUuid  string `json:"player_uuid"`
	ProfileUuid string `json:"profile_uuid"`
}

func (e PlayerNotFound) Error() string {
	if e.PlayerUuid != "" {
		return fmt.Sprintf("Player with player uuid %s not found", e.PlayerUuid)
	}

	if e.ProfileUuid != "" {
		return fmt.Sprintf("Player with profile uuid %s not found", e.ProfileUuid)
	}

	return "player not found, no player uuid or profile uuid given"
}

type Player struct {
	Id           PlayerId       `json:"id" bson:"_id"`
	GenericItems []*GenericItem `json:"generic_items" bson:"generic_items"`
	GenericPets  []*GenericPet  `json:"generic_pets" bson:"generic_pets"`
}

type PlayerId struct {
}

func PlayerByUuid(uuid string) (*Player, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id.playerUuid": uuid}

	var player Player
	res := playerCollection.FindOne(ctx, filter)

	if res.Err() != nil {

		if res.Err() == mongo.ErrNoDocuments {
			return nil, PlayerNotFound{PlayerUuid: uuid}
		}

		return nil, res.Err()
	}

	if err := res.Decode(&player); err != nil {
		return nil, err
	}

	return &player, nil
}
