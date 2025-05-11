// internal/services/mail_provider.go
package services

import (
	"context"
	"fmt"

	"github.com/wizenheimer/swiftcal/internal/config"
	"github.com/wizenheimer/swiftcal/pkg/logger"

	"github.com/mailgun/mailgun-go/v4"
	"go.uber.org/zap"
)

type EmailProvider interface {
	SendEmail(ctx context.Context, to, from, subject, textContent, htmlContent string, headers map[string]string) error
}

type MailgunProvider struct {
	client mailgun.Mailgun
	config *config.Config
}

func NewMailgunProvider(cfg *config.Config) *MailgunProvider {
	client := mailgun.NewMailgun(cfg.MailgunDomain, cfg.MailgunAPIKey)
	return &MailgunProvider{
		client: client,
		config: cfg,
	}
}

func (p *MailgunProvider) SendEmail(ctx context.Context, to, from, subject, textContent, htmlContent string, headers map[string]string) error {
	// if !p.config.IsProduction() {
	// 	logger.GetLogger().Info("Email not sent (development mode)",
	// 		zap.String("to", to),
	// 		zap.String("from", from),
	// 		zap.String("subject", subject),
	// 	)
	// 	return nil
	// }

	message := mailgun.NewMessage(from, subject, textContent, to)

	if htmlContent != "" {
		message.SetHTML(htmlContent)
	}

	// Add custom headers
	for key, value := range headers {
		message.AddHeader(key, value)
	}

	_, id, err := p.client.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send email via Mailgun: %w", err)
	}

	logger.GetLogger().Debug("Email sent via Mailgun",
		zap.String("message_id", id),
		zap.String("to", to),
		zap.String("subject", subject),
	)

	return nil
}
