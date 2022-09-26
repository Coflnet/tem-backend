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

	r.Use(CORSMiddleware())

	r.Use(otelgin.Middleware("tem-backend"))

	url := ginSwagger.URL("https://sky.coflnet.com/api/tem/swagger/doc.json") // The url pointing to API definition
	r.GET("/api/tem/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	r.GET("/api/tem/player/:uuid", playerByUuid)
	r.GET("/api/tem/playerProfile/:uuid", playerByProfileUuid)

	r.GET("/api/tem/item/:uuid", itemByUuid)
	r.GET("/api/tem/items/:id", itemsById)
	r.GET("/api/tem/coflItem/:uid", itemByCofluid)

	r.GET("/api/tem/pet/:uuid", petByUuid)

	return r
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
