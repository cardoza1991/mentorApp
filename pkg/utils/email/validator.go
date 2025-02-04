package email

import (
	"fmt"
	"net/smtp"
	"regexp"
	"strings"

	"github.com/beego/beego/v2/server/web"
)

// Email configuration constants
const (
	maxEmailLength = 254
	minEmailLength = 3
)

// EmailConfig holds the SMTP configuration
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// GetEmailConfig loads email configuration from app.conf
func GetEmailConfig() (*EmailConfig, error) {
	config := &EmailConfig{}

	config.Host = web.AppConfig.DefaultString("smtp_host", "")
	config.Port = web.AppConfig.DefaultInt("smtp_port", 587)
	config.Username = web.AppConfig.DefaultString("smtp_user", "")
	config.Password = web.AppConfig.DefaultString("smtp_password", "")
	config.From = web.AppConfig.DefaultString("smtp_from", "")

	if config.Host == "" || config.Username == "" || config.Password == "" {
		return nil, fmt.Errorf("incomplete SMTP configuration")
	}

	return config, nil
}

// ValidateEmail performs comprehensive email validation
func ValidateEmail(email string) error {

	// Check length constraints
	if len(email) > maxEmailLength {
		return fmt.Errorf("email exceeds maximum length of %d characters", maxEmailLength)
	}
	if len(email) < minEmailLength {
		return fmt.Errorf("email is shorter than minimum length of %d characters", minEmailLength)
	}

	// Basic format validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	// Additional checks
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return fmt.Errorf("invalid email format: must contain exactly one @ symbol")
	}

	local, domain := parts[0], parts[1]

	// Local part checks
	if len(local) > 64 {
		return fmt.Errorf("local part exceeds maximum length of 64 characters")
	}
	if strings.HasPrefix(local, ".") || strings.HasSuffix(local, ".") {
		return fmt.Errorf("local part cannot start or end with a dot")
	}
	if strings.Contains(local, "..") {
		return fmt.Errorf("local part cannot contain consecutive dots")
	}

	// Domain checks
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return fmt.Errorf("domain cannot start or end with a dot")
	}
	if strings.Contains(domain, "..") {
		return fmt.Errorf("domain cannot contain consecutive dots")
	}

	return nil
}

// SendEmail sends an email using the configured SMTP server
func SendEmail(to []string, subject, body string) error {
	config, err := GetEmailConfig()
	if err != nil {
		return fmt.Errorf("failed to get email config: %v", err)
	}

	// Compose email
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(fmt.Sprintf("Subject: %s\n%s\n%s", subject, mime, body))

	// Configure auth
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	// Send email
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	if err := smtp.SendMail(addr, auth, config.From, to, msg); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// SendVerificationEmail sends an email verification link
func SendVerificationEmail(email, token string) error {
	subject := "Verify Your Email Address"
	verificationURL := fmt.Sprintf("%s/verify-email/%s",
		web.AppConfig.DefaultString("base_url", "http://localhost:8080"),
		token)

	body := fmt.Sprintf(`
        <h2>Welcome to the Mentorship Platform!</h2>
        <p>Please click the link below to verify your email address:</p>
        <p><a href="%s">Verify Email</a></p>
        <p>If you didn't create this account, please ignore this email.</p>
    `, verificationURL)

	return SendEmail([]string{email}, subject, body)
}

// SendPasswordResetEmail sends a password reset link
func SendPasswordResetEmail(email, token string) error {
	subject := "Reset Your Password"
	resetURL := fmt.Sprintf("%s/reset-password/%s",
		web.AppConfig.DefaultString("base_url", "http://localhost:8080"),
		token)

	body := fmt.Sprintf(`
        <h2>Password Reset Request</h2>
        <p>You have requested to reset your password. Click the link below to proceed:</p>
        <p><a href="%s">Reset Password</a></p>
        <p>If you didn't request this, please ignore this email and ensure your account is secure.</p>
        <p>This link will expire in 24 hours.</p>
    `, resetURL)

	return SendEmail([]string{email}, subject, body)
}

// SendMentorshipRequestEmail notifies a mentor about a new mentorship request
func SendMentorshipRequestEmail(mentorEmail, menteeName string, programTitle string) error {
	subject := "New Mentorship Request"
	body := fmt.Sprintf(`
        <h2>New Mentorship Request</h2>
        <p>You have received a new mentorship request from %s for your program "%s".</p>
        <p>Please log in to your dashboard to review and respond to this request.</p>
    `, menteeName, programTitle)

	return SendEmail([]string{mentorEmail}, subject, body)
}
