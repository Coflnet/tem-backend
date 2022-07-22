package mongo

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Item struct {
}

func ItemsForPlayerUuid(_ string) ([]*Item, error) {
	return nil, fmt.Errorf("not implemented")
}

func ItemsForProfileUuid(uuid string) ([]*Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"currentOwner.profileUuid": uuid}

	var items []*Item
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
