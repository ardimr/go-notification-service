package consumerhandler

import (
	"bytes"
	"context"
	"encoding/json"
	"go_project_template/internal/mail"
	"go_project_template/internal/user/model"
	"html/template"
	"log"
	"path"
)

type IConsumerHandler interface {
	SendEmail(ctx context.Context, data []byte) error
}

type ConsumerHandler struct {
	sender mail.EmailSender
}

func NewConsumerHandler(sender mail.EmailSender) *ConsumerHandler {
	return &ConsumerHandler{
		sender: sender,
	}
}

func (ch *ConsumerHandler) SendEmail(ctx context.Context, data []byte) error {
	var userOTPVerificationEmailContent model.OTPVerificationEmailContent

	if err := json.Unmarshal(data, &userOTPVerificationEmailContent); err != nil {
		return err
	}

	log.Println("Sending OTP Code to", userOTPVerificationEmailContent.Email)

	// content := fmt.Sprintf(`Your OTP Code is %s or click <a href="%s">here</a>`, userOTPVerificationEmailContent.OTPCode, userOTPVerificationEmailContent.Url)

	temp := map[string]interface{}{
		"Product": "Mata Duitan",
		"OTPCode": userOTPVerificationEmailContent.OTPCode,
		"URL":     userOTPVerificationEmailContent.Url,
	}
	content, err := RenderTemplate(temp)

	if err != nil {
		return err
	}

	if err := ch.sender.SendEmail(
		"OTP Request",
		content,
		[]string{userOTPVerificationEmailContent.Email},
		[]string{},
		[]string{},
		[]string{},
	); err != nil {
		return err
	}

	return nil
}

func RenderTemplate(data map[string]interface{}) (string, error) {
	filepath := path.Join("/home/ardimr/workspace/portfolio/go_notification_service/internal/template", "confirm-email.html")

	tmpl, err := template.ParseFiles(filepath)

	if err != nil {
		return "", err
	}

	buff := new(bytes.Buffer)
	err = tmpl.Execute(buff, data)

	if err != nil {
		return "", err
	}

	return buff.String(), err
}
