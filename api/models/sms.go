package models

type SMS struct {
	MobilePhone string `json:"mobile_phone"`
	Message     string `json:"message"`
	From        string `json:"from"`
}

type SMSResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}