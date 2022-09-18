package mongo

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	Id           PlayerId      `json:"id" bson:"_id"`
	GenericItems []interface{} `json:"generic_items" bson:"generic_items"`
	GenericPets  []string      `json:"generic_pets" bson:"generic_pets"`
}

type PlayerId struct {
}

func PlayerByUuid(ctx context.Context, uuid string) (*Player, error) {
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
		log.Error().Err(err).Msgf("error decoding player with uuid %s", uuid)
		return nil, err
	}

	return &player, nil
}

func PlayerByProfileUuid(ctx context.Context, uuid string) (*Player, error) {
	filter := bson.M{"_id.profileUuid": uuid}

	var player Player
	res := playerCollection.FindOne(ctx, filter)

	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return nil, PlayerNotFound{ProfileUuid: uuid}
		}

		return nil, res.Err()
	}

	if err := res.Decode(&player); err != nil {
		log.Error().Err(err).Msgf("error decoding player with profile uuid %s", uuid)
		return nil, err
	}

	return &player, nil
}
