package utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/abstractmelon/is-site-live/internal/config"
	"gopkg.in/gomail.v2"
)

// EmailSender handles sending emails
type EmailSender struct {
	config config.SMTPConfig
}

// NewEmailSender creates a new email sender
func NewEmailSender(config config.SMTPConfig) *EmailSender {
	return &EmailSender{
		config: config,
	}
}

// SendDowntimeAlert sends a downtime alert email
func (e *EmailSender) SendDowntimeAlert(to, username, siteName, siteURL string, statusCode int, errorMessage string) error {
	// Check if SMTP is configured
	if e.config.Host == "" || e.config.User == "" || e.config.Password == "" || e.config.From == "" {
		return fmt.Errorf("SMTP not configured")
	}

	// Create message
	m := gomail.NewMessage()
	m.SetHeader("From", e.config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", fmt.Sprintf("Downtime Alert: %s is down", siteName))

	// Create email body
	body := fmt.Sprintf(`
		<h2>Downtime Alert</h2>
		<p>Hello %s,</p>
		<p>Your site <strong>%s</strong> is currently down.</p>
		<p><strong>URL:</strong> %s</p>
		<p><strong>Status Code:</strong> %s</p>
		<p><strong>Error:</strong> %s</p>
		<p>We'll notify you when the site is back up.</p>
		<p>Regards,<br>Is It Live Monitoring</p>
	`, username, siteName, siteURL, getStatusCodeText(statusCode), errorMessage)

	m.SetBody("text/html", body)

	// Create dialer
	d := gomail.NewDialer(e.config.Host, e.config.Port, e.config.User, e.config.Password)

	// Send email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// SendRecoveryAlert sends a recovery alert email
func (e *EmailSender) SendRecoveryAlert(to, username, siteName, siteURL string, statusCode int, downtime string) error {
	// Check if SMTP is configured
	if e.config.Host == "" || e.config.User == "" || e.config.Password == "" || e.config.From == "" {
		return fmt.Errorf("SMTP not configured")
	}

	// Create message
	m := gomail.NewMessage()
	m.SetHeader("From", e.config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", fmt.Sprintf("Recovery Alert: %s is back up", siteName))

	// Create email body
	body := fmt.Sprintf(`
		<h2>Recovery Alert</h2>
		<p>Hello %s,</p>
		<p>Good news! Your site <strong>%s</strong> is back up.</p>
		<p><strong>URL:</strong> %s</p>
		<p><strong>Status Code:</strong> %s</p>
		<p><strong>Downtime Duration:</strong> %s</p>
		<p>Regards,<br>Is It Live Monitoring</p>
	`, username, siteName, siteURL, getStatusCodeText(statusCode), downtime)

	m.SetBody("text/html", body)

	// Create dialer
	d := gomail.NewDialer(e.config.Host, e.config.Port, e.config.User, e.config.Password)

	// Send email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// getStatusCodeText returns a human-readable status code text
func getStatusCodeText(statusCode int) string {
	if statusCode == 0 {
		return "Connection Failed"
	}
	
	return strconv.Itoa(statusCode) + " " + http.StatusText(statusCode)
}
