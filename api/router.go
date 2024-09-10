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

	// Apply CORS middleware here globally
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Faqat frontend (localhost) uchun ruxsat
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig)) // CORS middleware applied

	r.POST("/partner", h.CreatePartner)
	r.GET("/partner/:id", h.GetPartner)
	r.GET("/partners", h.GetPartnerList)

	// Voting endpoints
	r.POST("/user/:partner_id/vote", h.AddScore)
	r.POST("/user/:partner_id/verify", h.VerifyCode)

	// Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
