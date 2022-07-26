package mongo

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Item struct {
	Id              string        `bson:"_id" json:"id"`
	Colour          string        `bson:"colour" json:"colour"`
	Enchantment     interface{}   `bson:"enchantments" json:"enchantments"`
	LastChecked     time.Time     `bson:"lastChecked" json:"lastChecked"`
	Location        string        `bson:"location" json:"location"`
	PreviousOwners  []interface{} `bson:"previousOwners" json:"previousOwners"`
	CurrentOwner    interface{}   `bson:"currentOwner" json:"currentOwner"`
	ExtraAttributes interface{}   `bson:"extraAttributes" json:"extraAttributes"`
	CreatedAt       time.Time     `bson:"createdAt" json:"createdAt"`
	Start           time.Time     `bson:"start" json:"start"`
	Reforge         string        `bson:"reforge" json:"reforge"`
	Rarity          string        `bson:"rarity" json:"rarity"`
}

func ItemsForPlayerUuid(uuid string) ([]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"currentOwner.playerUuid": uuid}

	var items []interface{}
	cur, err := itemCollection.Find(ctx, filter)

	if err != nil {
		log.Error().Err(err).Msgf("error finding items for player uuid %s", uuid)
		return nil, err
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			log.Error().Err(err).Msgf("error closing cursor for player uuid %s", uuid)
		}
	}(cur, ctx)

	if err := cur.All(ctx, &items); err != nil {
		log.Error().Err(err).Msgf("error decoding items for player uuid %s", uuid)
		return nil, err
	}

	return items, nil
}

func ItemsForProfileUuid(uuid string) ([]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"currentOwner.profileUuid": uuid}

	var items []interface{}
	cur, err := itemCollection.Find(ctx, filter)

	if err != nil {
		log.Error().Err(err).Msgf("error finding items for profile uuid %s", uuid)
		return nil, err
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			log.Error().Err(err).Msgf("error closing cursor for profile uuid %s", uuid)
		}
	}(cur, ctx)

	if err := cur.All(ctx, &items); err != nil {
		log.Error().Err(err).Msgf("error decoding items for profile uuid %s", uuid)
		return nil, err
	}

	return items, nil
}

func ItemById(id string) (*Item, error) {
	filter := bson.M{"_id": id}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var item Item
	if err := itemCollection.FindOne(ctx, filter).Decode(&item); err != nil {
		log.Error().Err(err).Msgf("error decoding item %s", id)
		return nil, err
	}

	return &item, nil
}

func ItemsById(id string, offset int64) ([]*Item, error) {

	filter := bson.M{"_id": id}
	opt := options.Find().SetLimit(1000).SetSkip(offset)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var items []*Item
	cur, err := itemCollection.Find(ctx, filter, opt)

	if err != nil {
		log.Error().Err(err).Msgf("error finding items for id %s", id)
		return nil, err
	}

	// close cursor
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			log.Error().Err(err).Msgf("error closing cursor for id %s", id)
		}
	}(cur, ctx)

	if err := cur.All(ctx, &items); err != nil {
		log.Error().Err(err).Msgf("error decoding items for id %s", id)
		return nil, err
	}

	return items, nil
}

func ItemsCountById(id string) (int64, error) {
	filter := bson.M{"_id": id}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := itemCollection.CountDocuments(ctx, filter)

	if err != nil {
		log.Error().Err(err).Msgf("error counting items for id %s", id)
		return 0, err
	}

	return count, nil
}
