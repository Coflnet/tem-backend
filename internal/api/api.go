package api

import (
	"github.com/Coflnet/tem-backend/docs"
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

func setupRouter() *gin.Engine {
	r := gin.Default()

	url := ginSwagger.URL("http://localhost:8080/api/tem/swagger/doc.json") // The url pointing to API definition
	r.GET("/api/tem/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/tem"
	docs.SwaggerInfo.Title = "TEM Backend"
	docs.SwaggerInfo.Version = "1.0"

	r.Use(otelgin.Middleware("tem-backend"))

	r.GET("/api/tem/player/:uuid", playerByUuid)
	r.GET("/api/tem/playerProfile/:uuid", playerByProfileUuid)

	r.GET("/api/tem/item/:uuid", itemByUuid)
	r.GET("/api/tem/items/:id", itemsById)

	return r
}
