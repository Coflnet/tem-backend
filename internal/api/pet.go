package api

import (
	"github.com/Coflnet/tem-backend/internal/mongo"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Pet by uuid
// @Summary ItemByUUID
// @Description returns the pet by its uuid
// @Tags pets
// @Accept json
// @Produce json
// @Param uuid path string true "id"
// @Success 200 {object} mongo.Pet
// @Router /pet/{uuid} [get]
func petByUuid(c *gin.Context) {
	id := c.Param("uuid")

	c.Writer.Header().Set("Cache-Control", "public, max-age=5, immutable")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
		return
	}

	pet, err := mongo.PetByUuid(c.Request.Context(), id)

	if serr, ok := err.(*mongo.PetNotFoundError); ok {
		c.JSON(http.StatusNotFound, gin.H{"error": serr.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pet)
}
