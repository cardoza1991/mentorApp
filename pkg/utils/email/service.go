// pkg/utils/email/service.go
package email

import (
    
    "log"
)

type EmailService struct {
    // Add email configuration here
    fromEmail string
    // Add any other necessary fields (SMTP settings, etc.)
}

func NewEmailService(fromEmail string) *EmailService {
    return &EmailService{
        fromEmail: fromEmail,
    }
}

func (s *EmailService) SendVerificationEmail(toEmail, token string) error {
    // Implement actual email sending logic here
    log.Printf("Sending verification email to %s with token %s", toEmail, token)
    return nil
}

func (s *EmailService) SendPasswordResetEmail(toEmail, token string) error {
    // Implement actual email sending logic here
    log.Printf("Sending password reset email to %s with token %s", toEmail, token)
    return nil
}
