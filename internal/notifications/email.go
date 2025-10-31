package email

import (
	"fmt"
	"net/smtp"
	"postman-task/pkg/config"
	"strings"
)

// Sends a simple text email.
func Send(to, subject, body string) error {
	cfg := config.Load()

	// Get smtp details from config
	host := cfg.Email.SMTPHost
	port := cfg.Email.SMTPPort
	username := cfg.Email.SMTPUsername
	password := cfg.Email.SMTPPassword
	from := cfg.Email.FromEmail

	if host == "" || port == "" || username == "" || password == "" || from == "" {
		return fmt.Errorf("smtp credentials not configured")
	}

	addr := host + ":" + port

	auth := smtp.PlainAuth("", username, password, host)

	msg := strings.Join([]string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-version: 1.0;",
		"Content-Type: text/plain; charset=\"UTF-8\";",
		"",
		body,
	}, "\r\n")

	return smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
}
