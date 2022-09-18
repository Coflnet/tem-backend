package mongo

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"os"
	"time"
)

var (
	client           *mongo.Client
	itemCollection   *mongo.Collection
	playerCollection *mongo.Collection
	petsCollection   *mongo.Collection
)

func Start() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opt := options.Client()
	opt.Monitor = otelmongo.NewMonitor()
	opt.ApplyURI(mongoUrl())

	var err error
	client, err = mongo.Connect(ctx, opt)
	if err != nil {
		log.Panic().Err(err).Msg("Error connecting to mongo")
	}

	itemCollection = client.Database("inventories").Collection("items")
	playerCollection = client.Database("inventories").Collection("players")
	petsCollection = client.Database("inventories").Collection("pets")
}

func Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return client.Disconnect(ctx)
}

func mongoUrl() string {
	u := os.Getenv("MONGO_URL")

	if u == "" {
		log.Panic().Msg("MONGO_URL not set")
	}

	return u
}
