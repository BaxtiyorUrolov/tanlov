package handler

import (
	"context"
	"fmt"
	"it-tanlov/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (h *Handler) AddVote(message *tgbotapi.Message) {
	// Xabardan video ID'ni olish
	videoID := message.Text[12:] // "vote_" dan keyingi qismni olish

	telegramID := message.From.ID

	exists, err := h.services.User().UserTelegramIDExist(context.Background(), telegramID)
	if err != nil {
		fmt.Println("Error while checking user tg exists")
		return
	}
	if exists {
		h.bot.Send(tgbotapi.NewMessage(int64(telegramID), "Kechirasiz siz allaqachon ovoz berib bo'lgansiz"))
		return
	}

	// Foydalanuvchini bazaga qo'shish
	err = h.services.User().Create(context.Background(), telegramID)
	if err != nil {
		h.log.Error("Error while creating user", logger.Error(err))
		h.bot.Send(tgbotapi.NewMessage(int64(telegramID), "Foydalanuvchini yaratishda xatolik yuz berdi"))
		return
	}

	// Video uchun ovoz qo'shish
	_, err = h.services.User().AddScore(context.Background(), videoID)
	if err != nil {
		h.log.Error("Error while adding score", logger.Error(err))
		h.bot.Send(tgbotapi.NewMessage(int64(telegramID), "Ovoz berishda xatolik yuz berdi"))
		return
	}

	// Foydalanuvchini muvaffaqiyatli ovoz bergani haqida xabardor qilish
	h.bot.Send(tgbotapi.NewMessage(int64(telegramID), "Ovoz berish muvaffaqiyatli yakunlandi!"))
}
