package email

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendResetCode(toEmail, code string) error
}

type GomailService struct {
	dialer *gomail.Dialer
	from   string
}

func NewGomailService(host string, port int, user, password, from string) *GomailService {
	d := gomail.NewDialer(host, port, user, password)
	return &GomailService{dialer: d, from: from}
}

func (s *GomailService) SendResetCode(toEmail, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Сброс пароля - School With AI")

	htmlBody := fmt.Sprintf(`
		<h1>Сброс пароля</h1>
		<p>Ваш код для сброса пароля: <b>%s</b></p>
		<p>Код действителен 15 минут.</p>
		<p>Если вы не запрашивали сброс, проигнорируйте это письмо.</p>
	`, code)

	m.SetBody("text/html", htmlBody)

	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
