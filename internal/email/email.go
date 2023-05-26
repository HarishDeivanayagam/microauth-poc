package email

import (
	"errors"
	"net/smtp"
)

type Service struct {
	username string
	password string
	url      string
	port     string
}

var (
	UnableToSend = errors.New("unable to send email")
	EmailSent    = "email sent successfully"
)

func New(username string, password string, url string, port string) *Service {
	return &Service{
		username: username,
		password: password,
		url:      url,
		port:     port,
	}
}

func (s *Service) SendEmail(to string, subject string, body string) (string, error) {
	auth := smtp.PlainAuth("", s.username, s.password, s.url)
	addr := s.url + ":" + s.port

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(addr, auth, s.username, []string{to}, msg)
	if err != nil {
		return "", UnableToSend
	}

	return EmailSent, nil
}
