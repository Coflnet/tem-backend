package api

import (
	_ "github.com/Coflnet/tem-backend/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func StartApi() error {
	r := setupRouter()

	return r.Run()
}

// @title           TEM Backend
// @version         1.0
// @description     A little backend for the tem db

// @contact.name   Flou21
// @contact.url    flou.dev
// @contact.email  muehlhans.f@coflnet.com

// @license.name  AGPL-3.0

// @host      sky.coflnet.com
// @BasePath  /api/tem/
func setupRouter() *gin.Engine {
	r := gin.Default()

	url := ginSwagger.URL("http://localhost:8080/api/tem/swagger/doc.json") // The url pointing to API definition
	r.GET("/api/tem/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	r.Use(otelgin.Middleware("tem-backend"))

	r.GET("/api/tem/player/:uuid", playerByUuid)
	r.GET("/api/tem/playerProfile/:uuid", playerByProfileUuid)

	r.GET("/api/tem/item/:uuid", itemByUuid)
	r.GET("/api/tem/items/:id", itemsById)
	r.GET("/api/tem/coflItem/:uid", itemByCofluid)

	return r
}
