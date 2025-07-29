package gomail

import (
	"dms-api/config"
	"fmt"

	"gopkg.in/gomail.v2"
)

func SendOTPViaMail(to string, otp string) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", config.Config("TO_EMAIL"))
	mail.SetHeader("To", to)
	mail.SetBody("text/plain", fmt.Sprintf("Your OTP code is: %s", otp))

	d := gomail.NewDialer("smtp.gmail.com", 587, config.Config("TO_EMAIL"), config.Config("APP_PASSWORD"))

	return d.DialAndSend(mail)
}
