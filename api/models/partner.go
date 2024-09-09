package models

import "time"

type Partner struct {
	ID        	string `json:"id"`
	FullName  	string `json:"full_name"`
	Phone     	string `json:"phone"`
	Email     	string `json:"email"`
	VideoLink 	string `json:"video_link"`
	VideoStatus bool   `json:"-"`
	Score 		int	   `json:"score"`
	ImagePath   string    `json:"image_path"`
	CreatedAt 	time.Time `json:"created_at"`
	UpdatedAt 	string `json:"updated_at"`
	DeletedAt 	string `json:"deleted_at"`
}

type CreatePartner struct {
	FullName  	string `json:"full_name"`
	Phone     	string `json:"phone"`
	Email     	string `json:"email"`
	VideoLink 	string `json:"video_link"`
}

type PartnerResponse struct {
	Partners []Partner
	Count      int
}