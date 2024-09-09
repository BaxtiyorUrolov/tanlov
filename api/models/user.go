//models/user.go

package models

type User struct {
	Phone     	string `json:"phone"`
	VideoID     string `json:"video_id"`
}

type VerifyCodeRequest struct {
    Phone    string `json:"phone"`
    Code     string `json:"code"`
    VideoID  string `json:"video_id"`
}