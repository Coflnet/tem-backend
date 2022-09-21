package api

import (
	"github.com/Coflnet/tem-backend/internal/mongo"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"sync"
)

// Player By UUID
// @Summary PlayerUUID
// @Description get a player by his player uuid
// @Tags player
// @Accept json
// @Produce json
// @Param uuid path string true "uuid"
// @Success 200 {object} PlayerResponse
// @Router /player/{uuid} [get]
func playerByUuid(c *gin.Context) {
	uuid := c.Param("uuid")
	c.Writer.Header().Set("Cache-Control", "public, max-age=30, immutable")

	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
		return
	}

	player := &PlayerResponse{}
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func(u string, p *PlayerResponse, waitGroup *sync.WaitGroup) {
		defer waitGroup.Done()

		items, err := mongo.ItemsForPlayerUuid(c.Request.Context(), u)

		if err != nil {
			log.Error().Err(err).Msgf("error searching items for player with uuid %s", u)
			return
		}

		p.Items = items
	}(uuid, player, &wg)

	wg.Add(1)
	go func(u string, p *PlayerResponse, waitGroup *sync.WaitGroup) {
		defer waitGroup.Done()

		player, err := mongo.PlayerByUuid(c.Request.Context(), u)
		if err != nil {
			log.Error().Err(err).Msgf("error searching player with uuid %s", u)
			return
		}

		log.Info().Msgf("found id: %s, %d items, %d pets; for uuid: %s", player.Id, len(player.GenericItems), len(player.GenericPets), u)

		p.Id = player.Id
		p.GenericItems = player.GenericItems
		p.GenericPets = player.GenericPets
	}(uuid, player, &wg)

	wg.Add(1)
	go func(u string, p *PlayerResponse, waitGroup *sync.WaitGroup) {
		defer wg.Done()

		pets, err := mongo.PetsOfPlayerUuid(c.Request.Context(), uuid)
		if err != nil {
			log.Error().Err(err).Msgf("error searching pets for player with uuid %s", u)
			return
		}

		p.Pets = pets
	}(uuid, player, &wg)

	wg.Wait()

	c.JSON(http.StatusOK, player)
}

// Player by Profile UUID
// @Summary ProfileUUID
// @Description get a player by his profile uuid
// @Tags player
// @Accept json
// @Produce json
// @Param uuid path string true "uuid"
// @Success 200 {object} PlayerResponse
// @Router /playerProfile/{uuid} [get]
func playerByProfileUuid(c *gin.Context) {
	uuid := c.Param("uuid")
	c.Writer.Header().Set("Cache-Control", "public, max-age=30, immutable")

	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
		return
	}

	player := &PlayerResponse{}
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func(u string, p *PlayerResponse, waitGroup *sync.WaitGroup) {
		defer waitGroup.Done()

		items, err := mongo.ItemsForProfileUuid(c.Request.Context(), u)

		if err != nil {
			log.Error().Err(err).Msgf("error searching items for player with uuid %s", u)
			return
		}

		p.Items = items
	}(uuid, player, &wg)

	wg.Add(1)
	go func(u string, p *PlayerResponse, waitGroup *sync.WaitGroup) {
		defer wg.Done()

		player, err := mongo.PlayerByProfileUuid(c.Request.Context(), uuid)
		if err != nil {
			log.Error().Err(err).Msgf("error searching player with uuid %s", uuid)
			return
		}

		log.Info().Msgf("found id: %s, %d items, %d pets; for uuid: %s", player.Id, len(player.GenericItems), len(player.GenericPets), uuid)

		p.Id = player.Id
		p.GenericItems = player.GenericItems
		p.GenericPets = player.GenericPets
	}(uuid, player, &wg)

	wg.Add(1)
	go func(u string, p *PlayerResponse, waitGroup *sync.WaitGroup) {
		defer waitGroup.Done()

		pets, err := mongo.PetsOfProfileUuid(c.Request.Context(), uuid)
		if err != nil {
			log.Error().Err(err).Msgf("error searching pets for player with uuid %s", uuid)
			return
		}

		p.Pets = pets
	}(uuid, player, &wg)

	wg.Wait()
	c.JSON(http.StatusOK, player)
}

type PlayerResponse struct {
	Id           mongo.PlayerId      `json:"id" bson:"_id"`
	GenericItems []mongo.GenericItem `json:"generic_items" bson:"generic_items"`
	GenericPets  []string            `json:"generic_pets" bson:"generic_pets"`
	Items        []mongo.Item        `json:"items"`
	Pets         []*mongo.Pet        `json:"pets"`
}
