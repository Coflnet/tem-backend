package mongo

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Pet struct {
	Id            string    `json:"id" bson:"_id"`
	Candy         int       `json:"candy" bson:"candy"`
	CurrentOwner  Owner     `json:"current_owner" bson:"currentOwner"`
	HeldItem      string    `json:"held_item" bson:"heldItem"`
	Level         int       `json:"level" bson:"level"`
	Location      string    `json:"location" bson:"location"`
	Name          string    `json:"name" bson:"name"`
	PreviousOwner Owner     `json:"previous_owner" bson:"previousOwner"`
	Rarity        string    `json:"rarity" bson:"rarity"`
	Skin          *string   `json:"skin" bson:"skin"`
	Start         time.Time `json:"start" bson:"start"`
	LastChecked   time.Time `json:"last_checked" bson:"lastChecked"`
}

type PetNotFoundError struct {
	Uuid string `json:"id,omitempty"`
}

func (p *PetNotFoundError) Error() string {
	return fmt.Sprintf("pet with id %s not found", p.Uuid)
}

func PetsOfPlayerUuid(ctx context.Context, uuid string) ([]*Pet, error) {
	filter := bson.M{"currentOwner.playerUuid": uuid}

	var pets []*Pet
	cur, err := petsCollection.Find(ctx, filter)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, PlayerNotFound{PlayerUuid: uuid}
		}

		return nil, err
	}

	if err := cur.All(ctx, &pets); err != nil {
		log.Error().Err(err).Msgf("error decoding player with uuid %s", uuid)
		return nil, err
	}

	return pets, nil
}

func PetsOfProfileUuid(ctx context.Context, uuid string) ([]*Pet, error) {
	filter := bson.M{"currentOwner.profileUuid": uuid}

	var pets []*Pet
	cur, err := petsCollection.Find(ctx, filter)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, PlayerNotFound{ProfileUuid: uuid}
		}

		return nil, err
	}

	if err := cur.All(ctx, &pets); err != nil {
		log.Error().Err(err).Msgf("error decoding player with uuid %s", uuid)
		return nil, err
	}

	return pets, nil
}

func PetByUuid(ctx context.Context, uuid string) (Pet, error) {
	filter := bson.M{"_id": uuid}

	var pet Pet
	err := petsCollection.FindOne(ctx, filter).Decode(&pet)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return pet, &PetNotFoundError{Uuid: uuid}
		}

		return pet, err
	}

	return pet, nil
}
