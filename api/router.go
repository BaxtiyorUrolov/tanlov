package api

import (
	"it-tanlov/api/handler"
	"it-tanlov/pkg/logger"
	"it-tanlov/service"
	"it-tanlov/storage"
	"time"

	_ "it-tanlov/api/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func New(services service.IServiceManager, storage storage.IStorage, log logger.ILogger, bot *tgbotapi.BotAPI) *gin.Engine {
	h := handler.New(services, storage, log, bot)

	r := gin.New()

	r.Use(gin.Logger())

	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},  // Frontend manzili
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))	

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Max-Age", "3600") // So'rovni cache qilish muddati

		// Agar OPTIONS so'rovi bo'lsa, javobni 204 status bilan qaytarish
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// Keyingi middlewarelarga o'tkazish
		c.Next()
	})

	r.POST("/partner", h.CreatePartner)
	r.GET("/partners", h.GetPartnerList)

	// Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
