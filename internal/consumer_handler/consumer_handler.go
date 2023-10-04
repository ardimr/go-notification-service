package consumerhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"go_project_template/internal/mail"
	"go_project_template/internal/user/model"
	"log"
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

	content := fmt.Sprintf(`Your OTP Code is %s or click <a href="%s">here</a>`, userOTPVerificationEmailContent.OTPCode, userOTPVerificationEmailContent.Url)
	if err := ch.sender.SendEmail(
		"OTP Request",
		content,
		[]string{"fixenog400@estudys.com"},
		[]string{},
		[]string{},
		[]string{},
	); err != nil {
		return err
	}

	return nil
}
