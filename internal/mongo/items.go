package mongo

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

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
