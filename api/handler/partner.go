package handler

import (
	"context"
	"fmt"
	"it-tanlov/api/models"
	"it-tanlov/config"
	"it-tanlov/pkg/check"
	"it-tanlov/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

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

	telegramUserID := config.Load().BotAdmin

	createPartner := models.CreatePartner{}
	if err := c.ShouldBindJSON(&createPartner); err != nil {
		handleResponse(c, h.log, "Error while reading body from client", http.StatusBadRequest, err)
		return
	}

	if !check.PhoneNumber(createPartner.Phone) {
		handleResponse(c, h.log, "Incorrect phone number", http.StatusBadRequest, nil)
		return
	}

	exists, err := h.services.Partner().PhoneExist(createPartner.Phone)
	if err != nil {
		handleResponse(c, h.log, "Error while checking phone existence", http.StatusInternalServerError, nil)
		return
	}
	if exists {
		handleResponse(c, h.log, "Phone number already exists", http.StatusBadRequest, "This phone number already exists")
		return
	}

	email, err := h.services.Partner().EmailExist(createPartner.Email)
	if err != nil {
		handleResponse(c, h.log, "Error while checking email existence", http.StatusInternalServerError, nil)
		return
	}
	if email {
		handleResponse(c, h.log, "Email already exists", http.StatusBadRequest, "This Email already exists")
		return
	}

	videoLink, err := h.services.Partner().VideoLinkExist(createPartner.VideoLink)
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

	acceptButton := tgbotapi.NewInlineKeyboardButtonData("Accept", "accept_partner_"+partner.ID)
	rejectButton := tgbotapi.NewInlineKeyboardButtonData("Reject", "reject_partner_"+partner.ID)

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(acceptButton, rejectButton),
	)

	msg := tgbotapi.NewMessage(telegramUserID, message)
	msg.ReplyMarkup = inlineKeyboard

	_, sendErr := h.bot.Send(msg)
	if sendErr != nil {
		h.log.Error("Error while sending message to Telegram user", logger.Error(sendErr))
	}

	handleResponse(c, h.log, "", http.StatusCreated, partner)
}

func (h *Handler) HandleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery) {
	data := callbackQuery.Data
	messageID := callbackQuery.Message.MessageID

	telegramUserID := config.Load().BotAdmin

	if data[:15] == "reject_partner_" {
		partnerID := data[15:]

		err := h.services.Partner().Delete(context.Background(), partnerID)
		if err != nil {
			h.log.Error("Error deleting partner", logger.Error(err))
			return
		}

		editMsg := tgbotapi.NewEditMessageText(telegramUserID, messageID, "Partner has been rejected and deleted.")
		if _, err := h.bot.Send(editMsg); err != nil {
			h.log.Error("Error sending edit message", logger.Error(err))
		}
	} else if data[:15] == "accept_partner_" {
		partnerID := data[15:]

		err := h.services.Partner().Update(context.Background(), partnerID)
		if err != nil {
			h.log.Error("Error verifying partner", logger.Error(err))
			return
		}

		// Tasdiqlangan partner haqida xabar yuborish
		editMsg := tgbotapi.NewEditMessageText(telegramUserID, messageID, "Partner has been accepted and verified.")
		if _, err := h.bot.Send(editMsg); err != nil {
			h.log.Error("Error sending edit message", logger.Error(err))
		}
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
