package service

import (
	"fmt"
	"log"
	"strconv"

	mail "github.com/wneessen/go-mail"
)

// EmailSender định nghĩa interface gửi email.
// Đổi implementation bằng cách thay constructor trong main.go:
//   - NewLogEmailSender      → dev/test (in ra terminal)
//   - NewSMTPEmailSender     → production (gửi qua SMTP)
type EmailSender interface {
	Send(to, subject, htmlBody string) error
}

// LogEmailSender gửi email giả, chỉ in nội dung ra stdout (dùng cho dev/test)
type LogEmailSender struct {
	frontendURL string
}

func NewLogEmailSender(frontendURL string) *LogEmailSender {
	return &LogEmailSender{frontendURL: frontendURL}
}

func (s *LogEmailSender) Send(to, subject, htmlBody string) error {
	_ = htmlBody // body có thể chứa dữ liệu nhạy cảm (reset code/token), không log raw content.
	log.Printf("\n"+
		"╔══════════════════════════════════════════════════╗\n"+
		"║            📧  EMAIL (dev mode)                 ║\n"+
		"╠══════════════════════════════════════════════════╣\n"+
		"║ To:      %s\n"+
		"║ Subject: %s\n"+
		"╠══════════════════════════════════════════════════╣\n"+
		"║ Body:    [REDACTED]\n"+
		"╚══════════════════════════════════════════════════╝\n",
		to, subject,
	)
	return nil
}

func (s *LogEmailSender) FrontendURL() string { return s.frontendURL }

// SMTPEmailSender gửi email thật qua SMTP (dùng cho production)
type SMTPEmailSender struct {
	host        string
	port        int
	user        string
	pass        string
	from        string
	fromName    string
	frontendURL string
}

func NewSMTPEmailSender(host, portStr, user, pass, from, fromName, frontendURL string) *SMTPEmailSender {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 587 // default TLS port
	}
	// Nếu from trống, dùng user làm from
	if from == "" {
		from = user
	}
	// Nếu fromName trống, dùng from làm fromName
	if fromName == "" {
		fromName = from
	}
	return &SMTPEmailSender{
		host:        host,
		port:        port,
		user:        user,
		pass:        pass,
		from:        from,
		fromName:    fromName,
		frontendURL: frontendURL,
	}
}

func (s *SMTPEmailSender) Send(to, subject, htmlBody string) error {
	m := mail.NewMsg()
	if err := m.FromFormat(s.fromName, s.from); err != nil {
		return fmt.Errorf("failed to set From: %w", err)
	}
	if err := m.To(to); err != nil {
		return fmt.Errorf("failed to set To: %w", err)
	}
	m.Subject(subject)
	m.SetBodyString(mail.TypeTextHTML, htmlBody)

	c, err := mail.NewClient(s.host,
		mail.WithPort(s.port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(s.user),
		mail.WithPassword(s.pass),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
	)
	if err != nil {
		return fmt.Errorf("failed to create mail client: %w", err)
	}

	if err = c.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *SMTPEmailSender) FrontendURL() string { return s.frontendURL }
