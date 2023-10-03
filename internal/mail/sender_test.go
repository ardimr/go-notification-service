package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/gomail.v2"
)

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_SENDER_NAME = "Rizky Ardi Maulana <rizkyardimaulana@gmail.com>"
const CONFIG_AUTH_EMAIL = ""
const CONFIG_AUTH_PASSWORD = ""

func TestSendEmailWithGmail(t *testing.T) {
	// Dial
	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)
	// Create a new sender
	senderGmail := NewGmailSender(
		dialer,
		CONFIG_SENDER_NAME,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	err := senderGmail.SendEmail(
		"test mail",
		"Hello, <b>This is a test email</b>",
		[]string{"fixenog400@estudys.com"},
		[]string{},
		[]string{},
		[]string{},
	)

	require.NoError(t, err)
}
