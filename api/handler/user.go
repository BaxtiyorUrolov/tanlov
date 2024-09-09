package handler

import (
	"context"
	"fmt"
	"it-tanlov/api/models"
	"it-tanlov/pkg/check"
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

// Function to send verification code via SMS (simulated here, should use actual SMS service)
func sendCode(phoneNumber string, code string) {
	fmt.Printf("SMS kod: %s raqamiga yuborildi: %s\n", phoneNumber, code)
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

	// Check if the phone number already exists (has the user already voted?)
	exists, err := check.IUserPhoneExist(user.Phone, h.storage.User())
	if err != nil {
		handleResponse(c, h.log, "Error while checking phone number existence", http.StatusInternalServerError, err)
		return
	}
	if exists {
		handleResponse(c, h.log, "You have already voted", http.StatusBadRequest, "You have already voted")
		return
	}

	// Generate and send code
	code := generateCode()
	cacheStore.Set(user.Phone, code, cache.DefaultExpiration)
	sendCode(user.Phone, code)

	// Inform the user that the SMS has been sent
	handleResponse(c, h.log, "SMS code sent", http.StatusCreated, gin.H{
		"message":    "SMS code sent, please verify",
		"phone":      user.Phone,
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

    // Parse incoming JSON to VerifyCodeRequest struct
    var req models.VerifyCodeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        fmt.Println("VerifyCode: Error while reading body from client:", err)
        handleResponse(c, h.log, "Error while reading body from client", http.StatusBadRequest, err)
        return
    }
    fmt.Println("VerifyCode: Received request body:", req)

    phone := req.Phone
    inputCode := req.Code
	videoID := req.VideoID
	
    // Check the entered code
    cachedCode, found := cacheStore.Get(phone)
    if !found {
        fmt.Println("VerifyCode: Code not found in cache for phone:", phone)
        handleResponse(c, h.log, "Code not found or expired", http.StatusBadRequest, "Code not found or expired")
        return
    }
    fmt.Println("VerifyCode: Cached Code =", cachedCode)

    if cachedCode != inputCode {
        fmt.Println("VerifyCode: Incorrect code entered. Cached Code =", cachedCode, "Input Code =", inputCode)
        handleResponse(c, h.log, "Incorrect code", http.StatusBadRequest, "Incorrect code")
        return
    }


    // Record the vote
    err := h.services.User().Create(context.Background(), phone)
    if err != nil {
        fmt.Println("VerifyCode: Error while creating user:", err)
        handleResponse(c, h.log, "Error while creating user", http.StatusInternalServerError, err)
        return
    }
    fmt.Println("VerifyCode: User created successfully")

    addScore, err := h.services.User().AddScore(context.Background(), videoID)
    if err != nil {
        fmt.Println("VerifyCode: Error while casting vote:", err)
        handleResponse(c, h.log, "Error while casting vote", http.StatusInternalServerError, err.Error())
        return
    }
    fmt.Println("VerifyCode: Vote successfully recorded")

    handleResponse(c, h.log, "Vote successfully recorded", http.StatusOK, addScore)
    fmt.Println("VerifyCode: Response sent successfully")
}
