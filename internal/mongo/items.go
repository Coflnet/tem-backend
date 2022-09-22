package mongo

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Item struct {
	Id              string      `bson:"_id" json:"id"`
	ItemId          string      `bson:"itemId" json:"itemId"`
	Colour          int         `bson:"colour,omitempty" json:"colour"`
	Enchantment     interface{} `bson:"enchantments" json:"enchantments"`
	LastChecked     time.Time   `bson:"lastChecked" json:"lastChecked"`
	Location        string      `bson:"location" json:"location"`
	PreviousOwners  []Owner     `bson:"previousOwners" json:"previousOwners"`
	CurrentOwner    Owner       `bson:"currentOwner" json:"currentOwner"`
	ExtraAttributes interface{} `bson:"extraAttributes" json:"extraAttributes"`
	CreatedAt       time.Time   `bson:"createdAt" json:"createdAt"`
	Start           time.Time   `bson:"start" json:"start"`
	Reforge         string      `bson:"reforge" json:"reforge"`
	Rarity          string      `bson:"rarity" json:"rarity"`
}

type Owner struct {
	PlayerUuid  string `bson:"playerUuid" json:"playerUuid"`
	ProfileUuid string `bson:"profileUuid" json:"profileUuid"`
}

type ItemNotFoundError struct {
	Id      string `json:"id,omitempty"`
	CoflUid string `json:"coflUid,omitempty"`
}

func (i *ItemNotFoundError) Error() string {
	return fmt.Sprintf("item %s not found", i.Id)
}

func ItemsForPlayerUuid(ctx context.Context, uuid string) ([]Item, error) {
	filter := bson.M{"currentOwner.playerUuid": uuid}

	var items []Item
	cur, err := itemCollection.Find(ctx, filter)

	if err != nil {
		log.Error().Err(err).Msgf("error finding items for player uuid %s", uuid)
		return nil, err
	}

	defer func(cur *mongo.Cursor) {
		err := cur.Close(ctx)
		if err != nil {
			log.Error().Err(err).Msgf("error closing cursor for player uuid %s", uuid)
		}
	}(cur)

	if err := cur.All(ctx, &items); err != nil {
		log.Error().Err(err).Msgf("error decoding items for player uuid %s", uuid)
		return nil, err
	}

	return items, nil
}

func ItemsForProfileUuid(ctx context.Context, uuid string) ([]Item, error) {
	filter := bson.M{"currentOwner.profileUuid": uuid}

	var items []Item
	cur, err := itemCollection.Find(ctx, filter)

	if err != nil {
		log.Error().Err(err).Msgf("error finding items for profile uuid %s", uuid)
		return nil, err
	}

	defer func(cur *mongo.Cursor) {
		err := cur.Close(ctx)
		if err != nil {
			log.Error().Err(err).Msgf("error closing cursor for profile uuid %s", uuid)
		}
	}(cur)

	if err := cur.All(ctx, &items); err != nil {
		log.Error().Err(err).Msgf("error decoding items for profile uuid %s", uuid)
		return nil, err
	}

	return items, nil
}

func ItemById(ctx context.Context, id string) (*Item, error) {
	filter := bson.M{"_id": id}

	var item Item
	res := itemCollection.FindOne(ctx, filter)

	if err := res.Decode(&item); err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, &ItemNotFoundError{Id: id}
		}

		log.Error().Err(err).Msgf("error decoding item for id %s", id)
		return nil, err
	}

	return &item, nil
}

func ItemsById(ctx context.Context, id string, offset int64) ([]*Item, error) {

	filter := bson.M{"itemId": id}
	opt := options.Find().
		SetSort(bson.D{{"created_at", -1}}).
		SetLimit(1000).
		SetSkip(offset)

	var items []*Item
	cur, err := itemCollection.Find(ctx, filter, opt)

	if err != nil {
		log.Error().Err(err).Msgf("error finding items for id %s", id)
		return nil, err
	}

	// close cursor
	defer func(cur *mongo.Cursor) {
		err := cur.Close(ctx)
		if err != nil {
			log.Error().Err(err).Msgf("error closing cursor for id %s", id)
		}
	}(cur)

	if err := cur.All(ctx, &items); err != nil {
		log.Error().Err(err).Msgf("error decoding items for id %s", id)
		return nil, err
	}

	return items, nil
}

func ItemsCountById(ctx context.Context, id string) (int64, error) {
	filter := bson.M{"itemId": id}

	count, err := itemCollection.CountDocuments(ctx, filter)

	if err != nil {
		log.Error().Err(err).Msgf("error counting items for id %s", id)
		return 0, err
	}

	return count, nil
}

func ItemByCoflUid(ctx context.Context, uid string) (Item, error) {
	filter := bson.M{"coflUid": uid}

	var item Item
	res := itemCollection.FindOne(ctx, filter)

	if err := res.Decode(&item); err != nil {

		if err == mongo.ErrNoDocuments {
			log.Warn().Msgf("item with cofl uid %s not found", uid)
			return item, &ItemNotFoundError{CoflUid: uid}
		}

		log.Error().Err(err).Msgf("error decoding item for cofl uid %s", uid)
		return item, err
	}

	return item, nil
}
