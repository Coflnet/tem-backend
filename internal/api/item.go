package api

import (
	"github.com/Coflnet/tem-backend/internal/mongo"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"sync"
)

type ItemResponse struct {
	Items  []*mongo.Item `json:"items"`
	Count  int64         `json:"count"`
	Offset int64         `json:"offset"`
}

// Item by uuid
// @Summary ItemByUUID
// @Description returns the item by its uuid
// @Tags items
// @Accept json
// @Produce json
// @Success 200 {object} mongo.Item
// @Router /item/{uuid} [get]
func ItemByUuid(c *gin.Context) {
	id := c.Param("uuid")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
		return
	}

	item, err := mongo.ItemById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// Item by uuid
// @Summary ItemByUUID
// @Description returns the item by its uuid
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Query offset path string true "offset"
// @Success 200 {object} ItemResponse
// @Router /items/{id} [get]
func ItemsById(c *gin.Context) {
	id := c.Param("id")
	offsetStr := c.Query("offset")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var offset int64 = 0
	if offsetStr != "" {
		var err error
		o, err := strconv.Atoi(offsetStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "offset must be an integer"})
			return
		}

		offset = int64(o)
	}

	response := ItemResponse{}
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func(r *ItemResponse, waitGroup *sync.WaitGroup, i string, o int64) {
		defer waitGroup.Done()

		val, e := mongo.ItemsById(i, o)
		if e != nil {
			log.Error().Err(e).Msgf("error searching items for player with id %s", id)
			return
		}

		r.Items = val
	}(&response, &wg, id, offset)

	go func(r *ItemResponse, waitGroup *sync.WaitGroup) {
		defer waitGroup.Done()

		count, e := mongo.ItemsCountById(id)
		if e != nil {
			log.Error().Err(e).Msgf("error searching items for player with id %s", id)
			return
		}

		r.Count = count
	}(&response, &wg)

	response.Offset = offset

	wg.Wait()
	c.JSON(http.StatusOK, response)
}
