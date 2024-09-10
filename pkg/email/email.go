package email

import (
	"fmt"
	"net/smtp"
)

// Example Email configuration
var smtpServer = "smtp.example.com"
var smtpPort = "587"
var senderEmail = "no-reply@example.com"
var senderPassword = "yourpassword"

// SendEmail sends a verification code via email to the user's email address
func SendEmail(email, code string) error {
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)

	to := []string{email}
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: Verification Code\r\n\r\nYour verification code is: %s", email, code))

	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, senderEmail, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	fmt.Printf("Email sent to %s with code %s\n", email, code)
	return nil
}
