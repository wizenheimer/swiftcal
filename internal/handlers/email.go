// internal/handlers/email.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wizenheimer/swiftcal/internal/config"
	"github.com/wizenheimer/swiftcal/internal/models"
	"github.com/wizenheimer/swiftcal/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/swiftcal/pkg/logger"
	"go.uber.org/zap"
)

type EmailHandler struct {
	emailService *services.EmailService
	config       *config.Config
}

func NewEmailHandler(emailService *services.EmailService, cfg *config.Config) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
		config:       cfg,
	}
}

func (h *EmailHandler) HandleMailgunWebhook(c *fiber.Ctx) error {
	if c.Method() != "POST" {
		return c.Status(http.StatusMethodNotAllowed).JSON(fiber.Map{
			"error": "Method not allowed",
		})
	}

	// Check spam filtering headers from Mailgun
	spamFlag := c.Get("X-Mailgun-Sflag")
	spamScore := c.Get("X-Mailgun-Sscore")

	if spamFlag == "Yes" || (spamScore != "" && spamScore > "0.5") {
		logger.GetLogger().Warn("Rejected spam email", zap.String("spam_flag", spamFlag), zap.String("spam_score", spamScore))
		return c.JSON(fiber.Map{
			"message": "rejected spam",
		})
	}

	// Parse form-encoded webhook data
	webhook, files, err := h.parseMailgunWebhook(c)
	if err != nil {
		logger.GetLogger().Error("Failed to parse Mailgun webhook", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse webhook data",
		})
	}

	// Process the email
	if err := h.emailService.HandleWebhook(c.Context(), webhook, files); err != nil {
		logger.GetLogger().Error("Failed to process email webhook", zap.Error(err))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process email",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Email processed successfully",
	})
}

func (h *EmailHandler) parseMailgunWebhook(c *fiber.Ctx) (*models.EmailWebhook, []models.EmailFile, error) {
	// Mailgun sends webhook data as form-encoded data, not multipart form data
	// Parse form data using FormValue instead of MultipartForm
	webhook := &models.EmailWebhook{}
	var files []models.EmailFile

	// Extract form values using FormValue
	getValue := func(key string) string {
		return c.FormValue(key)
	}

	// Transform Mailgun format to match SendGrid expectations
	webhook.Subject = getValue("subject")
	webhook.Text = getValue("body-plain")
	webhook.HTML = getValue("body-html")
	webhook.From = getValue("from")
	webhook.To = getValue("recipient")
	webhook.Timestamp = getValue("timestamp")

	// Construct headers from Mailgun's message-headers
	webhook.Headers = h.constructMailgunHeaders(getValue("message-headers"), getValue("timestamp"), webhook.Subject, webhook.From, webhook.To)

	// Construct envelope
	envelope := map[string]interface{}{
		"from": getValue("sender"),
		"to":   []string{getValue("recipient")},
	}
	envelopeBytes, _ := json.Marshal(envelope)
	webhook.Envelope = string(envelopeBytes)

	// Set spam filtering results (default to pass since we filtered at header level)
	webhook.SPF = "pass"
	webhook.DKIM = "pass"

	// Note: Mailgun typically doesn't send file attachments in webhooks
	// If you need to handle attachments, you would need to fetch them separately
	// using Mailgun's API with the message ID

	return webhook, files, nil
}

func (h *EmailHandler) constructMailgunHeaders(messageHeaders, timestamp, subject, from, to string) string {
	headers := ""
	var messageID string
	var existingReferences string

	// Parse message-headers if available
	if messageHeaders != "" {
		var headersList [][2]string
		if err := json.Unmarshal([]byte(messageHeaders), &headersList); err == nil {
			for _, header := range headersList {
				headers += fmt.Sprintf("%s: %s\n", header[0], header[1])

				if header[0] == "Message-Id" {
					messageID = header[1]
				}
				if header[0] == "References" {
					existingReferences = header[1]
				}
			}
		}
	}

	// Fallback headers if message-headers parsing failed
	if headers == "" {
		if timestamp != "" {
			// Convert timestamp to date
			headers += fmt.Sprintf("Date: %s\n", timestamp)
		}
		if subject != "" {
			headers += fmt.Sprintf("Subject: %s\n", subject)
		}
		if from != "" {
			headers += fmt.Sprintf("From: %s\n", from)
		}
		if to != "" {
			headers += fmt.Sprintf("To: %s\n", to)
		}
	}

	// Add threading headers for responses
	if messageID != "" {
		headers += fmt.Sprintf("In-Reply-To: %s\n", messageID)

		if existingReferences != "" {
			headers += fmt.Sprintf("References: %s %s\n", existingReferences, messageID)
		} else {
			headers += fmt.Sprintf("References: %s\n", messageID)
		}
	}

	return headers
}
