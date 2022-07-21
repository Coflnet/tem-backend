package main

import (
	"github.com/Coflnet/tem-backend/internal/api"
	"github.com/Coflnet/tem-backend/internal/mongo"
	"github.com/rs/zerolog/log"
)

func main() {

	mongo.Start()
	defer mongo.Stop()

	err := api.StartApi()

	if err != nil {
		log.Panic().Err(err).Msgf("Error starting API")
	}
}
