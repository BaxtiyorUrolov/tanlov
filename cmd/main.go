//cmd/main.go

package main

import (
	"context"
	"it-tanlov/api"
	"it-tanlov/api/handler"
	"it-tanlov/config"
	"it-tanlov/pkg/logger"
	"it-tanlov/service"
	"it-tanlov/storage/postgres"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.ServiceName)

	store, err := postgres.New(context.Background(), cfg, log)
	if err != nil {
		log.Error("Error while connecting to DB: %v", logger.Error(err))
		return
	}
	defer store.Close()

	services := service.New(store, log)

	botToken := cfg.BotToken
	botInstance, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Error("Error creating Telegram bot API: %v", logger.Error(err))
		return
	}

	h := handler.New(services, store, log, botInstance)

	server := api.New(services, store, log, botInstance)

	server.Use(cors.Default())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		offset := 0
		for {
			select {
			case <-ctx.Done():
				log.Info("Shutting down bot...")
				return
			default:
				updates, err := botInstance.GetUpdates(tgbotapi.NewUpdate(offset))
				if err != nil {
					log.Error("Error getting updates: %v", logger.Error(err))
					time.Sleep(5 * time.Second)
					continue
				}
				for _, update := range updates {
					// Umumiy `HandleUpdate` funksiyasini chaqirish
					h.HandleUpdate(update)
					offset = update.UpdateID + 1
				}
			}
		}
	}()

	server.StaticFS("/static", http.Dir("./static"))
	server.StaticFile("/", "./static/index.html")

	if err := server.Run("195.2.84.169:2005"); err != nil {
		log.Error("Error while running server: %v", logger.Error(err))
	}
}
