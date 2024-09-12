package api

import (
	"it-tanlov/api/handler"
	"it-tanlov/pkg/logger"
	"it-tanlov/service"
	"it-tanlov/storage"
	"net/http"

	_ "it-tanlov/api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func NewRouter(services service.IServiceManager, storage storage.IStorage, log logger.ILogger, bot *tgbotapi.BotAPI) *gin.Engine {
	h := handler.New(services, storage, log, bot)

	r := gin.New()

	r.Use(gin.Logger())

	// Middleware'ni qo'shish
	r.Use(Middleware())

	r.POST("/partner", h.CreatePartner)
	r.GET("/partners", h.GetPartnerList)

	// Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

// Middleware funktsiyasi
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, GET")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With") //nolint:lll
		c.Header("Access-Control-Max-Age", "3600")
		c.Header("Content-Type", "application/json")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
