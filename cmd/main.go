package main

import (
	"context"
	"it-tanlov/api"
	"it-tanlov/api/handler"
	"it-tanlov/config"
	"it-tanlov/pkg/logger"
	"it-tanlov/service"
	"it-tanlov/storage/postgres"
	"os"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
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

	h := handler.New(services, log, botInstance)

	server := api.New(cfg, log, services, store, botInstance, h)

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
					h.HandleUpdate(update)
					offset = update.UpdateID + 1
				}
			}
		}
	}()

	go func() {
		if err := server.Run(); err != nil {
			log.Error("Error run http server", zap.Error(err))
		}
	}()

	<-ctx.Done()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("failed http graceful shutdown", zap.Error(err))
	}
}
