package services

import (
	"context"
	"log"
)

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) SendEmail(ctx context.Context, to, subject, body string) error {
	log.Printf("📧 Sending Email\nTo: %s\nSubject: %s\nBody: %s\n", to, subject, body)

	// Later → integrate SMTP / provider here

	return nil
}
