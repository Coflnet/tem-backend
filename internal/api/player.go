package api

import (
	"github.com/Coflnet/tem-backend/internal/mongo"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"sync"
)

// Player
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

	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
		return
	}

	var player *PlayerResponse
	wg := sync.WaitGroup{}

	go func(u string, p *PlayerResponse, waitGroup *sync.WaitGroup) {

		waitGroup.Add(1)
		defer waitGroup.Done()

		items, err := mongo.ItemsForPlayerUuid(u)

		if err != nil {
			log.Error().Err(err).Msgf("error searching items for player with uuid %s", u)
			return
		}

		p.Items = items
	}(uuid, player, &wg)

	go func(u string, p *PlayerResponse, waitGroup *sync.WaitGroup) {
		waitGroup.Add(1)
		defer waitGroup.Done()

		player, err := mongo.PlayerByUuid(u)
		if err != nil {
			log.Error().Err(err).Msgf("error searching player with uuid %s", u)
			return
		}

		p.Id = player.Id
		p.GenericItems = player.GenericItems
		p.GenericPets = player.GenericPets
	}(uuid, player, &wg)

	wg.Wait()
	c.JSON(http.StatusOK, player)
}

func playerByProfileUuid(c *gin.Context) {
}

type PlayerResponse struct {
	Id           mongo.PlayerId       `json:"id" bson:"_id"`
	GenericItems []*mongo.GenericItem `json:"generic_items" bson:"generic_items"`
	GenericPets  []*mongo.GenericPet  `json:"generic_pets" bson:"generic_pets"`

	Items []*mongo.Item `json:"items"`
}
