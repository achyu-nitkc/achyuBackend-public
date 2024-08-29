package auth

import (
	"fmt"
	"net/smtp"
	"strings"
)

func SendVerifyCode(email, code string) error {
	from := "hirogoshawko3249@gmail.com"
	subject := "confirm"
	body := "Thank you for using achyu\n" +
		"To confirm, input this number\n" +
		code + "\nhttps://localhost:3000/home"
	hostname, port, username, password := config()
	auth := smtp.PlainAuth("", username, password, hostname)
	msg := []byte(strings.ReplaceAll(fmt.Sprintf("To: %s\nSubject: %s\n\n%s", email, subject, body), "\n", "\r\n"))
	err := smtp.SendMail(fmt.Sprintf("%s:%d", hostname, port), auth, from, []string{email}, msg)
	if err != nil {
		return err
	}
	return nil
}
