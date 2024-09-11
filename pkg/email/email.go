package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

// Yandex Mail konfiguratsiyasi
var smtpServer = "smtp.yandex.com"
var smtpPort = "587" // TLS uchun port
var senderEmail = "hello@asljons.com" // Yandex pochta manzilingiz
var senderPassword = "bfxhjbyzjdvtckvb" // App Password

// SendEmail funksiyasi tasdiqlash kodini emailga yuboradi
func SendEmail(email, code string) error {
	// SMTP server konfiguratsiyasi
	serverAddress := smtpServer + ":" + smtpPort
	fmt.Println("SMTP serverga ulanish:", serverAddress)

	// SMTP ulanishi
	fmt.Println("TCP orqali ulanishni boshlayapmiz...")
	client, err := smtp.Dial(serverAddress)
	if err != nil {
		fmt.Printf("SMTP serverga ulanishda xatolik: %v\n", err)
		return fmt.Errorf("failed to dial SMTP server: %v", err)
	}
	defer client.Close()

	// STARTTLS orqali ulanish
	fmt.Println("STARTTLS ulanishini boshlayapmiz...")
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer,
	}
	if err = client.StartTLS(tlsConfig); err != nil {
		fmt.Printf("STARTTLS ulanishda xatolik: %v\n", err)
		return fmt.Errorf("failed to start TLS: %v", err)
	}
	fmt.Println("STARTTLS ulanishi muvaffaqiyatli!")

	// SMTP autentifikatsiyasi
	fmt.Println("SMTP serverda autentifikatsiya qilish...")
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)
	if err = client.Auth(auth); err != nil {
		fmt.Printf("SMTP serverda autentifikatsiya xatoligi: %v\n", err)
		return fmt.Errorf("failed to authenticate SMTP client: %v", err)
	}
	fmt.Println("Autentifikatsiya muvaffaqiyatli o'tdi!")

	// Yuboruvchi ma'lumotlarini qo'shish
	fmt.Println("Yuboruvchi manzili:", senderEmail)
	if err = client.Mail(senderEmail); err != nil {
		fmt.Printf("Yuboruvchi ma'lumotlarini qo'shishda xatolik: %v\n", err)
		return fmt.Errorf("failed to set sender: %v", err)
	}
	fmt.Println("Yuboruvchi muvaffaqiyatli o'rnatildi!")

	// Qabul qiluvchi ma'lumotlarini qo'shish
	fmt.Println("Qabul qiluvchi manzili:", email)
	if err = client.Rcpt(email); err != nil {
		fmt.Printf("Qabul qiluvchi ma'lumotlarini qo'shishda xatolik: %v\n", err)
		return fmt.Errorf("failed to set recipient: %v", err)
	}
	fmt.Println("Qabul qiluvchi muvaffaqiyatli o'rnatildi!")

	// Xabarni yuborish
	fmt.Println("Yuboriladigan xabar ma'lumotlarini tayyorlash...")
	wc, err := client.Data()
	if err != nil {
		fmt.Printf("Xabar yuborish jarayonida xatolik: %v\n", err)
		return fmt.Errorf("failed to send data: %v", err)
	}
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: Verification Code\r\n\r\nYour verification code is: %s", email, code))
	fmt.Println("Xabar tayyorlandi, yuborilmoqda...")
	_, err = wc.Write(msg)
	if err != nil {
		fmt.Printf("Xabar yozishda xatolik: %v\n", err)
		return fmt.Errorf("failed to write message: %v", err)
	}
	fmt.Println("Xabar muvaffaqiyatli yozildi!")

	err = wc.Close()
	if err != nil {
		fmt.Printf("Ulanishni yopishda xatolik: %v\n", err)
		return fmt.Errorf("failed to close connection: %v", err)
	}
	fmt.Println("Yuborish jarayoni muvaffaqiyatli yakunlandi!")

	// SMTP ulanishini yopish
	fmt.Println("SMTP ulanishini yopish...")
	if err = client.Quit(); err != nil {
		fmt.Printf("SMTP ulanishini yopishda xatolik: %v\n", err)
		return fmt.Errorf("failed to close SMTP session: %v", err)
	}
	fmt.Println("SMTP ulanishi muvaffaqiyatli yopildi!")

	fmt.Printf("Email muvaffaqiyatli yuborildi: %s\n", email)
	return nil
}
