package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"it-tanlov/api/models"
	"it-tanlov/config"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

func GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999
	return fmt.Sprintf("%06d", rand.Intn(max-min+1)+min)
}

func Send(toNumber, code string) error {
	fmt.Println("Code: ", code)
	fromNumber := "4546"
	apiURL := "https://notify.eskiz.uz/api/message/sms/send"

	cfg := config.Load()

	smsData := models.SMS{
		MobilePhone: toNumber,        
		Message:     "This is test from Eskiz",
		From:        fromNumber,           
	}

	jsonData, err := json.Marshal(smsData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+cfg.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var smsResponse models.SMSResponse
	err = json.Unmarshal(body, &smsResponse)
	if err != nil {
		return err
	}

	fmt.Println("Message ID:", smsResponse.MessageID)
	fmt.Println("Status:", smsResponse.Status)

	return nil
}