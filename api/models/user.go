//models/user.go

package models

type User struct {
	Email     	string `json:"email"`
	VideoID     string `json:"video_id"`
}

type VerifyCodeRequest struct {
    Email    string `json:"email"`
    Code     string `json:"code"`
    VideoID  string `json:"video_id"`
}