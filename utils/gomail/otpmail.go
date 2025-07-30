package gomail

import (
	"dms-api/config"
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

func SendOTPViaMail(to string, otp string) error {
	mail := gomail.NewMessage()
	senderEmail := config.Config("TO_EMAIL")
	appPassword := config.Config("APP_PASSWORD")

	mail.SetHeader("From", senderEmail)
	mail.SetHeader("To", to)
	mail.SetBody("text/plain", fmt.Sprintf("Your OTP code is: %s", otp))

	d := gomail.NewDialer("smtp.gmail.com", 587, senderEmail, appPassword)

	//return d.DialAndSend(mail)

	err := d.DialAndSend(mail)
	if err != nil {
		log.Printf("Error sending email: %v", err) // Log the error
		return fmt.Errorf("failed to send email: %w", err)
	}
	log.Printf("OTP sent successfully to %s", to)
	return nil
}
