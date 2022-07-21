package mongo

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

var (
	itemCollection   *mongo.Collection
	playerCollection *mongo.Collection
	petsCollection   *mongo.Collection
)

func Start() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl()))
	if err != nil {
		log.Panic().Err(err).Msg("Error connecting to mongo")
	}

	itemCollection = client.Database("inventories").Collection("items")
	playerCollection = client.Database("inventories").Collection("players")
	petsCollection = client.Database("inventories").Collection("pets")
}

func Stop() {

}

func mongoUrl() string {
	u := os.Getenv("MONGO_URL")

	if u == "" {
		log.Panic().Msg("MONGO_URL not set")
	}

	return u
}
