package api

import (
	_ "github.com/Coflnet/tem-backend/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func StartApi() error {
	r := setupRouter()

	return r.Run()
}

// @title TEM Backend
// @version 1.0
// @description Some endpoints for the tem db

// @contact.name Flou21
// @contact.email muehlhans.f@coflnet.com

// @license.name AGPL v3

// @host sky.coflnet.com
// @BasePath /api/tem/
func setupRouter() *gin.Engine {
	r := gin.Default()

	url := ginSwagger.URL("https://sky.coflnet.com/api/tem/swagger/index.html") // The url pointing to API definition
	r.GET("/api/tem/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	r.GET("/api/tem/player/:uuid", playerByUuid)
	r.GET("/api/tem/playerProfile/:uuid", playerByProfileUuid)

	return r
}
