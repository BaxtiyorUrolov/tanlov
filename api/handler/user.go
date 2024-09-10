package handler

import (
	"context"
	"fmt"
	"it-tanlov/api/models"
	"it-tanlov/pkg/check"
    "it-tanlov/pkg/email"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

// Global cache object
var cacheStore = cache.New(5*time.Minute, 10*time.Minute)

// Function to generate verification code
func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func sendCode(mail string, code string) {
	err := email.SendEmail(mail, code)
	if err != nil {
		fmt.Printf("Failed to send email to: %s, error: %v\n", mail, err)
	} else {
		fmt.Printf("Verification code sent to email: %s\n", mail)
	}
}

// AddScore godoc
// @Router       /user/{partner_id}/vote [POST]
// @Summary      Send SMS code to user
// @Description  Get the user's phone number and send SMS code for voting
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user body  models.User false "user"
// @Success      201  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) AddScore(c *gin.Context) {
	user := models.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		handleResponse(c, h.log, "Error while reading body from client", http.StatusBadRequest, err)
		return
	}

	// Check if the email already exists (has the user already voted?)
	exists, err := check.IUserEmailExist(user.Email, h.storage.User()) // Change to check for email
	if err != nil {
		handleResponse(c, h.log, "Error while checking email existence", http.StatusInternalServerError, err)
		return
	}
	if exists {
		handleResponse(c, h.log, "You have already voted", http.StatusBadRequest, "You have already voted")
		return
	}

	// Generate and send code via email
	code := generateCode()
	cacheStore.Set(user.Email, code, cache.DefaultExpiration) // Store code with email key
	sendCode(user.Email, code)                                // Send to email instead of phone

	// Inform the user that the email has been sent
	handleResponse(c, h.log, "Email verification code sent", http.StatusCreated, gin.H{
		"message":    "Email verification code sent, please verify",
		"email":      user.Email,
		"partner_id": user.VideoID, // Video ID
	})
}

// VerifyCode godoc
// @Router       /user/{partner_id}/verify [POST]
// @Summary      Verify code and cast vote
// @Description  Verify the code entered by the user and cast their vote
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        code body models.VerifyCodeRequest false "verification"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) VerifyCode(c *gin.Context) {
	var req models.VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleResponse(c, h.log, "Error while reading body from client", http.StatusBadRequest, err)
		return
	}

	email := req.Email
	inputCode := req.Code
	videoID := req.VideoID

	// Check the entered code
	cachedCode, found := cacheStore.Get(email)
	if !found {
		handleResponse(c, h.log, "Code not found or expired", http.StatusBadRequest, "Code not found or expired")
		return
	}

	if cachedCode != inputCode {
		handleResponse(c, h.log, "Incorrect code", http.StatusBadRequest, "Incorrect code")
		return
	}

	// Record the vote
	err := h.services.User().Create(context.Background(), email)
	if err != nil {
		handleResponse(c, h.log, "Error while creating user", http.StatusInternalServerError, err)
		return
	}

	addScore, err := h.services.User().AddScore(context.Background(), videoID)
	if err != nil {
		handleResponse(c, h.log, "Error while casting vote", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "Vote successfully recorded", http.StatusOK, addScore)
}
