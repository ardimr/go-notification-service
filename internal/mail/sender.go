package mail

import (
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
	dialer            *gomail.Dialer
}

func NewGmailSender(dialer *gomail.Dialer, name string, fromEmailAddress string, fromEmailPassword string) *GmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailAddress,
		dialer:            dialer,
	}
}

func (sender *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", sender.fromEmailAddress)
	mailer.SetHeader("To", to...)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", content)

	err := sender.dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Mail sent!")

	return err
}

func RenderTemplate(ctx *gin.Context) {
	filepath := path.Join("/home/ardimr/workspace/portfolio/go_notification_service/internal/template", "confirm-email.html")

	tmpl, err := template.ParseFiles(filepath)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	data := map[string]interface{}{
		"Product": "Mata Duitan",
		"OTPCode": "123456",
		"URL":     "http://localhost:8080/api/verify-otp?otp_code=123456",
	}

	err = tmpl.Execute(ctx.Writer, data)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
}
