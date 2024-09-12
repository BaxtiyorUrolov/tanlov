//api/handler/partner.go

package handler

import (
	"context"
	"fmt"
	"it-tanlov/api/models"
	"it-tanlov/pkg/check"
	"it-tanlov/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const telegramUserID int64 = 6358749851

// CreatePartner godoc
// @Router       /partner [POST]
// @Summary      Creates a new partner
// @Description  create a new partner
// @Tags         partner
// @Accept       json
// @Produce      json
// @Param        partner body models.CreatePartner false "partner"
// @Success      201  {object}  models.Partner
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) CreatePartner(c *gin.Context) {
	createPartner := models.CreatePartner{}
	if err := c.ShouldBindJSON(&createPartner); err != nil {
		handleResponse(c, h.log, "Error while reading body from client", http.StatusBadRequest, err)
		return
	}

	if !check.PhoneNumber(createPartner.Phone) {
		handleResponse(c, h.log, "Incorrect phone number", http.StatusBadRequest, nil)
		return
	}

	exists, err := check.IPhoneExist(createPartner.Phone, h.storage.Partner())
	if err != nil {
		handleResponse(c, h.log, "Error while checking phone existence", http.StatusInternalServerError, nil)
		return
	}
	if exists {
		handleResponse(c, h.log, "Phone number already exists", http.StatusBadRequest, "This phone number already exists")
		return
	}

	email, err := check.IEmailExist(createPartner.Email, h.storage.Partner())
	if err != nil {
		handleResponse(c, h.log, "Error while checking email existence", http.StatusInternalServerError, nil)
		return
	}
	if email {
		handleResponse(c, h.log, "Email already exists", http.StatusBadRequest, "This Email already exists")
		return
	}

	videoLink, err := check.IVideoLinkExist(createPartner.VideoLink, h.storage.Partner())
	if err != nil {
		handleResponse(c, h.log, "Error while checking video link existence", http.StatusInternalServerError, nil)
		return
	}
	if videoLink {
		handleResponse(c, h.log, "Video link already exists", http.StatusBadRequest, "This Vide  link already exists")
		return
	}

	partner, err := h.services.Partner().Create(context.Background(), createPartner)
	if err != nil {
		handleResponse(c, h.log, "error is while creating partners", http.StatusInternalServerError, err.Error())
		return
	}

	message := fmt.Sprintf("New partner created!\nName: %s\nPhone: %s\nEmail: %s\nVideo link: %s", partner.FullName, partner.Phone, partner.Email, partner.VideoLink)
	
	// Create inline buttons with callback data
	acceptButton := tgbotapi.NewInlineKeyboardButtonData("Accept", "accept_partner_"+partner.ID)
	rejectButton := tgbotapi.NewInlineKeyboardButtonData("Reject", "reject_partner_"+partner.ID)
	
	// Arrange buttons in a single row
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(acceptButton, rejectButton),
	)
	
	msg := tgbotapi.NewMessage(telegramUserID, message)
	msg.ReplyMarkup = inlineKeyboard // Attach inline keyboard to the message

	_, sendErr := h.bot.Send(msg) // h.bot is the Telegram bot instance
	if sendErr != nil {
		h.log.Error("Error while sending message to Telegram user", logger.Error(sendErr))
	}

	handleResponse(c, h.log, "", http.StatusCreated, partner)
}

func (h *Handler) HandleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery) {
	data := callbackQuery.Data
	messageID := callbackQuery.Message.MessageID

	if data[:15] == "reject_partner_" {
		partnerID := data[15:]

		// Delete the partner from the storage
		err := h.storage.Partner().Delete(context.Background(), partnerID)
		if err != nil {
			h.log.Error("Error deleting partner", logger.Error(err))
			return
		}

		// Notify the user about the rejection
		editMsg := tgbotapi.NewEditMessageText(telegramUserID, messageID, "Partner has been rejected and deleted.")
		if _, err := h.bot.Send(editMsg); err != nil {
			h.log.Error("Error sending edit message", logger.Error(err))
		}
	} else if data[:15] == "accept_partner_" {
		partnerID := data[15:]

		// So'rang: admin rasmini yuborishi kerak
		msg := tgbotapi.NewMessage(telegramUserID, "Please upload the image for the partner.")
		if _, err := h.bot.Send(msg); err != nil {
			h.log.Error("Error sending image request message", logger.Error(err))
		}

		// Admin rasm yuborishni kutilmoqda, partnerID ni saqlaymiz
		h.pendingPartnerImages[telegramUserID] = partnerID
	}
}

func (h *Handler) HandleImageMessage(update tgbotapi.Update) {
	if update.Message.Photo != nil {
		// Eng katta o'lchamdagi rasmni olamiz
		photo := (*update.Message.Photo)[len(*update.Message.Photo)-1]
		fileID := photo.FileID

		// PartnerID ni oldindan kiritilgan partner uchun olamiz
		partnerID, exists := h.pendingPartnerImages[int64(update.Message.From.ID)]
		if !exists {
			h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "No pending partner for this image."))
			return
		}

		// Telegram API orqali rasmni yuklab olish uchun URL'ni olamiz
		fileURL, err := h.bot.GetFileDirectURL(fileID)
		if err != nil {
			h.log.Error("Error retrieving file from Telegram", logger.Error(err))
			return
		}

		// Partner yozuvini yangilash, rasm yo'li bilan
		err = h.services.Partner().Update(context.Background(), partnerID, fileURL)
		if err != nil {
			h.log.Error("Error updating partner with image", logger.Error(err))
			return
		}

		// Adminni muvaffaqiyatli yangilash haqida xabardor qilish
		h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Image uploaded and partner updated successfully."))

		// Pending image'ni olib tashlash
		delete(h.pendingPartnerImages, int64(update.Message.From.ID))
	}
}


// GetPartnerList godoc
// @Router       /partners [GET]
// @Summary      Get partner list
// @Description  get partner list
// @Tags         partner
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param        limit query string false "limit"
// @Param        search query string false "search"
// @Success      201  {object}  models.PartnerResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetPartnerList(c *gin.Context) {
	var (
		page, limit int
		search      string
		err         error
	)

	pageStr := c.DefaultQuery("page", "1")
	page, err = strconv.Atoi(pageStr)
	if err != nil {
		handleResponse(c, h.log, "error is while converting page", http.StatusBadRequest, err.Error())
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		handleResponse(c, h.log, "error is while converting limit", http.StatusBadRequest, err.Error())
		return
	}

	search = c.Query("search")

	partners, err := h.services.Partner().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	if err != nil {
		handleResponse(c, h.log, "error is while get list", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, partners)
}
