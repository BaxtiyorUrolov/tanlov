//api/handler/handler.go

package handler

import (
	"it-tanlov/api/models"
	"it-tanlov/pkg/logger"
	"it-tanlov/service"
	"it-tanlov/storage"

	"github.com/gin-gonic/gin"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Handler struct {
	services service.IServiceManager
	storage  storage.IStorage
	log      logger.ILogger
	bot      *tgbotapi.BotAPI
	pendingPartnerImages map[int64]string 
}

func New(services service.IServiceManager, storage storage.IStorage, log logger.ILogger, bot *tgbotapi.BotAPI) Handler {
	return Handler{
		services: services,
		storage:  storage,
		log:      log,
		bot: bot,
		pendingPartnerImages: make(map[int64]string),
	}
}

func (h *Handler) HandleUpdate(update tgbotapi.Update) {
    if update.Message != nil {
        // Agar xabar rasm bo'lsa
        if update.Message.Photo != nil {
            h.HandleImageMessage(update)
        } else if update.Message.Text != "" {
            // Xabar matn bo'lsa, bu yerda qayta ishlashingiz mumkin
            h.log.Info("Text message received", logger.String("text", update.Message.Text))
        }
    } else if update.CallbackQuery != nil {
        h.HandleCallbackQuery(update.CallbackQuery)
    }
}


func handleResponse(c *gin.Context, log logger.ILogger, msg string, statusCode int, data interface{}) {
	resp := models.Response{}

	switch code := statusCode; {
	case code < 400:
		resp.Description = "OK"
		log.Info("~~~~> OK", logger.String("msg", msg), logger.Any("status", code))
	case code < 500:
		resp.Description = "Bad Request"
		log.Error("!!!!! BAD REQUEST", logger.String("msg", msg), logger.Any("status", code))
	default:
		resp.Description = "Internal Server Error"
		log.Error("!!!!! INTERNAL SERVER ERROR", logger.String("msg", msg), logger.Any("status", code))
	}

	resp.StatusCode = statusCode
	resp.Data = data

	c.JSON(resp.StatusCode, resp)
}