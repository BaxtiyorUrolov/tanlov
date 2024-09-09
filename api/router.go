package api

import (
	"it-tanlov/api/handler"
	"it-tanlov/pkg/logger"
	"it-tanlov/service"
	"it-tanlov/storage"

	_ "it-tanlov/api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func New(services service.IServiceManager, storage storage.IStorage, log logger.ILogger, bot *tgbotapi.BotAPI) *gin.Engine {
	h := handler.New(services, storage, log, bot)

	r := gin.New()

	r.Use(gin.Logger())

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
