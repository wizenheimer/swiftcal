// internal/services/email_service.go
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/wizenheimer/swiftcal/internal/config"
	"github.com/wizenheimer/swiftcal/internal/models"
	"github.com/wizenheimer/swiftcal/internal/utils"
	"github.com/wizenheimer/swiftcal/templates"

	"github.com/google/uuid"
	"github.com/wizenheimer/swiftcal/pkg/logger"
	"go.uber.org/zap"
)

type EmailService struct {
	config          *config.Config
	authService     *AuthService
	calendarService *CalendarService
	openaiService   *OpenAIService
	emailProvider   EmailProvider
}

func NewEmailService(cfg *config.Config, authService *AuthService, calendarService *CalendarService, openaiService *OpenAIService) *EmailService {
	var emailProvider EmailProvider

	if cfg.MailgunAPIKey != "" {
		emailProvider = NewMailgunProvider(cfg)
	}

	return &EmailService{
		config:          cfg,
		authService:     authService,
		calendarService: calendarService,
		openaiService:   openaiService,
		emailProvider:   emailProvider,
	}
}

func (s *EmailService) HandleWebhook(ctx context.Context, webhook *models.EmailWebhook, files []models.EmailFile) error {
	sender := s.getSenderFromEmail(webhook)
	logger.GetLogger().Info("Processing email webhook", zap.String("sender", sender))

	// Verify email authenticity
	if !s.verifyEmail(webhook) {
		logger.GetLogger().Warn("Email failed verification")
		return s.sendUnverifiedEmailResponse(ctx, sender, webhook)
	}

	// Check if this is a support email
	recipients := s.getRecipientsFromEmail(webhook)
	if s.isSupportEmail(recipients, webhook.Subject) {
		return s.forwardToSupport(ctx, sender, webhook)
	}

	// Get user from email
	user, err := s.authService.GetUserByEmail(ctx, sender)
	if err != nil {
		logger.GetLogger().Info("User not found, sending signup invitation", zap.String("sender", sender))
		return s.sendSignupInvitation(ctx, sender, webhook)
	}

	// Determine action from subject
	action := s.parseSubjectAction(webhook.Subject)
	logger.GetLogger().Info("Processing user request",
		zap.String("user_id", user.ID.String()),
		zap.String("action", action))

	switch action {
	case "addUser":
		return s.handleAddEmailAddress(ctx, user, webhook)
	case "removeEmail":
		return s.handleRemoveEmailAddress(ctx, user, webhook)
	case "deleteAccount":
		return s.handleDeleteAccount(ctx, user, webhook)
	case "addEvent":
		return s.handleAddEvent(ctx, user, webhook, files)
	default:
		return s.handleAddEvent(ctx, user, webhook, files)
	}
}

func (s *EmailService) getSenderFromEmail(webhook *models.EmailWebhook) string {
	var envelope struct {
		From string `json:"from"`
	}

	if err := json.Unmarshal([]byte(webhook.Envelope), &envelope); err != nil {
		logger.GetLogger().Warn("Failed to parse envelope, using From header", zap.Error(err))
		return strings.ToLower(webhook.From)
	}

	return strings.ToLower(envelope.From)
}

func (s *EmailService) getRecipientsFromEmail(webhook *models.EmailWebhook) []string {
	var envelope struct {
		To []string `json:"to"`
	}

	if err := json.Unmarshal([]byte(webhook.Envelope), &envelope); err != nil {
		logger.GetLogger().Warn("Failed to parse envelope recipients", zap.Error(err))
		return []string{strings.ToLower(webhook.To)}
	}

	var recipients []string
	for _, to := range envelope.To {
		recipients = append(recipients, strings.ToLower(to))
	}

	return recipients
}

func (s *EmailService) verifyEmail(webhook *models.EmailWebhook) bool {
	// Check SPF
	if webhook.SPF != "pass" {
		return false
	}

	// Check DKIM
	if !strings.Contains(strings.ToLower(webhook.DKIM), "pass") {
		return false
	}

	return true
}

func (s *EmailService) isSupportEmail(recipients []string, subject string) bool {
	for _, recipient := range recipients {
		if strings.Contains(recipient, "support@") ||
			strings.Contains(recipient, "admin@") {
			return true
		}
	}

	// Check for account verification emails
	if strings.Contains(strings.ToLower(subject), "verify your email address") {
		return true
	}

	return false
}

func (s *EmailService) parseSubjectAction(subject string) string {
	subject = strings.ToLower(strings.TrimSpace(subject))

	if strings.HasPrefix(subject, "add ") {
		return "addUser"
	} else if strings.HasPrefix(subject, "remove ") {
		return "removeEmail"
	} else if strings.HasPrefix(subject, "delete account") {
		return "deleteAccount"
	} else if strings.HasPrefix(subject, "fwd") {
		return "addEvent"
	}

	return "addEvent"
}

func (s *EmailService) sendUnverifiedEmailResponse(ctx context.Context, sender string, webhook *models.EmailWebhook) error {
	template := templates.GetUnverifiedEmailTemplate(sender, s.config.EmailDomain)
	return s.sendEmailResponse(ctx, sender, webhook, template, true)
}

func (s *EmailService) sendSignupInvitation(ctx context.Context, sender string, webhook *models.EmailWebhook) error {
	template := templates.GetNoUserFoundTemplate(sender, s.config.BaseDomain, s.config.AppDomain, s.config.EmailDomain)
	return s.sendEmailResponse(ctx, sender, webhook, template, true)
}

func (s *EmailService) forwardToSupport(ctx context.Context, sender string, webhook *models.EmailWebhook) error {
	content := fmt.Sprintf("From: %s<br><br>Subject: %s<br><br>%s", sender, webhook.Subject, webhook.HTML)

	return s.emailProvider.SendEmail(ctx,
		"support@"+s.extractDomain(s.config.MainEmailAddress),
		s.config.MainEmailAddress,
		webhook.Subject,
		webhook.Text,
		content,
		nil,
	)
}

func (s *EmailService) handleAddEmailAddress(ctx context.Context, user *models.User, webhook *models.EmailWebhook) error {
	emailRegex := regexp.MustCompile(`^add\s+([a-zA-Z0-9._+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6})$`)
	matches := emailRegex.FindStringSubmatch(strings.ToLower(webhook.Subject))

	if len(matches) != 2 {
		logger.GetLogger().Warn("Invalid add email format, treating as event")
		return s.handleAddEvent(ctx, user, webhook, nil)
	}

	emailToAdd := matches[1]

	// Check if email already exists
	existingUser, err := s.authService.GetUserByEmail(ctx, emailToAdd)
	if err == nil && existingUser != nil {
		template := templates.GetAdditionalEmailInUseTemplate(emailToAdd, s.config.EmailDomain)
		return s.sendEmailResponse(ctx, user.Email, webhook, template, true)
	}

	// Add to pending emails
	verificationCode, err := s.authService.AddPendingEmailAddress(ctx, user.ID, user.Email, emailToAdd)
	if err != nil {
		return fmt.Errorf("failed to add pending email: %w", err)
	}

	// Send verification email
	template := templates.GetAddAdditionalEmailTemplate(verificationCode.String(), user.Email, s.config.AppDomain, s.config.EmailDomain)
	return s.sendEmailResponse(ctx, emailToAdd, webhook, template, false)
}

func (s *EmailService) handleRemoveEmailAddress(ctx context.Context, user *models.User, webhook *models.EmailWebhook) error {
	emailRegex := regexp.MustCompile(`^remove\s+([a-zA-Z0-9._+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6})$`)
	matches := emailRegex.FindStringSubmatch(strings.ToLower(webhook.Subject))

	if len(matches) != 2 {
		logger.GetLogger().Warn("Invalid remove email format, treating as event")
		return s.handleAddEvent(ctx, user, webhook, nil)
	}

	emailToRemove := matches[1]

	// Check if email belongs to this user
	existingUser, err := s.authService.GetUserByEmail(ctx, emailToRemove)
	if err != nil || existingUser.ID != user.ID {
		template := templates.GetRemovalEmailInUseTemplate(emailToRemove, s.config.EmailDomain)
		return s.sendEmailResponse(ctx, user.Email, webhook, template, true)
	}

	// Remove email address
	if err := s.authService.RemoveEmailAddress(ctx, emailToRemove); err != nil {
		return fmt.Errorf("failed to remove email: %w", err)
	}

	template := templates.GetEmailAddressRemovedTemplate(emailToRemove, s.config.EmailDomain)
	return s.sendEmailResponse(ctx, user.Email, webhook, template, true)
}

func (s *EmailService) handleDeleteAccount(ctx context.Context, user *models.User, webhook *models.EmailWebhook) error {
	if err := s.authService.DeleteUser(ctx, user.ID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	template := templates.GetUserDeletedTemplate(s.config.EmailDomain)
	return s.sendEmailResponse(ctx, user.Email, webhook, template, true)
}

func (s *EmailService) handleAddEvent(ctx context.Context, user *models.User, webhook *models.EmailWebhook, files []models.EmailFile) error {
	// Check for ICS attachments first
	for _, file := range files {
		if strings.HasSuffix(strings.ToLower(file.Filename), ".ics") {
			return s.handleICSEvent(ctx, user, webhook, file)
		}
	}

	// Process email with AI
	return s.handleAIEvent(ctx, user, webhook)
}

func (s *EmailService) handleICSEvent(ctx context.Context, user *models.User, webhook *models.EmailWebhook, icsFile models.EmailFile) error {
	// Parse ICS file using existing utility
	event, err := utils.ParseICSFile(icsFile.Content)
	if err != nil {
		logger.GetLogger().Error("Failed to parse ICS file", zap.Error(err))
		template := templates.GetUnableToParseTemplate(s.config.EmailDomain)
		return s.sendEmailResponse(ctx, user.Email, webhook, template, true)
	}

	// Add to calendar using the parsed event directly
	addedEvent, err := s.calendarService.AddEvent(ctx, user.ID, event)
	if err != nil {
		logger.GetLogger().Error("Failed to add ICS event to calendar", zap.Error(err))
		template := templates.GetOAuthFailedTemplate(s.config.AppDomain, s.config.EmailDomain)
		return s.sendEmailResponse(ctx, user.Email, webhook, template, true)
	}

	// Send success response
	template := templates.GetEventAddedTemplate(
		addedEvent.HTMLLink,
		s.formatEventDate(addedEvent.StartTime, addedEvent.TimeZone),
		s.formatAttendees(addedEvent.Attendees),
		s.config.EmailDomain,
	)
	return s.sendEmailResponse(ctx, user.Email, webhook, template, true)
}

func (s *EmailService) handleAIEvent(ctx context.Context, user *models.User, webhook *models.EmailWebhook) error {
	// Extract headers
	headers := s.parseEmailHeaders(webhook.Headers)

	// Process with OpenAI
	eventsResponse, _, err := s.openaiService.ProcessEmail(
		ctx,
		webhook.Text,
		webhook.Subject,
		webhook.From,
		headers["Date"],
	)

	if err != nil {
		logger.GetLogger().Error("OpenAI processing failed", zap.Error(err))
		template := templates.GetUnableToParseTemplate(s.config.EmailDomain)
		return s.sendEmailResponse(ctx, user.Email, webhook, template, true)
	}

	if eventsResponse.Error != nil {
		template := templates.GetAIParseErrorTemplate(*eventsResponse.Description, s.config.EmailDomain)
		return s.sendEmailResponse(ctx, user.Email, webhook, template, true)
	}

	if len(eventsResponse.Events) == 0 {
		template := templates.GetUnableToParseTemplate(s.config.EmailDomain)
		return s.sendEmailResponse(ctx, user.Email, webhook, template, true)
	}

	// Process events
	var successfulEvents []*models.GoogleCalendarEvent
	var failedEvents []error

	for _, event := range eventsResponse.Events {
		// Validate and filter attendees
		event.Attendees = s.filterValidEmails(event.Attendees)

		calEvent, err := s.calendarService.AddEvent(ctx, user.ID, &event)
		if err != nil {
			logger.GetLogger().Error("Failed to add event",
				zap.Error(err),
				zap.String("summary", event.Summary))
			failedEvents = append(failedEvents, err)
			continue
		}

		successfulEvents = append(successfulEvents, calEvent)
	}

	if len(successfulEvents) == 0 {
		template := templates.GetOAuthFailedTemplate(s.config.AppDomain, s.config.EmailDomain)
		return s.sendEmailResponse(ctx, user.Email, webhook, template, true)
	}

	// Send success response
	if len(successfulEvents) == 1 {
		event := successfulEvents[0]
		if len(event.Attendees) > 1 {
			// Multiple attendees - show invite link
			inviteLink := s.buildInviteLink(user.ID, event.ID, "primary", event.Attendees)
			template := templates.GetEventAddedAttendeesTemplate(
				event.HTMLLink,
				s.formatEventDate(event.StartTime, event.TimeZone),
				inviteLink,
				s.formatAttendees(event.Attendees),
				s.config.EmailDomain,
			)
			return s.sendEmailResponse(ctx, user.Email, webhook, template, true)
		} else {
			// Single attendee
			template := templates.GetEventAddedTemplate(
				event.HTMLLink,
				s.formatEventDate(event.StartTime, event.TimeZone),
				s.formatAttendees(event.Attendees),
				s.config.EmailDomain,
			)
			return s.sendEmailResponse(ctx, user.Email, webhook, template, true)
		}
	} else {
		// Multiple events - custom response
		return s.sendMultipleEventsResponse(ctx, user.Email, webhook, successfulEvents, failedEvents)
	}
}

func (s *EmailService) filterValidEmails(emails []string) []string {
	var valid []string
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	for _, email := range emails {
		if emailRegex.MatchString(email) {
			valid = append(valid, email)
		} else {
			logger.GetLogger().Warn("Filtering out invalid email", zap.String("email", email))
		}
	}

	return valid
}

func (s *EmailService) parseEmailHeaders(headerString string) map[string]string {
	headers := make(map[string]string)
	lines := strings.Split(headerString, "\n")

	for _, line := range lines {
		if colonIndex := strings.Index(line, ":"); colonIndex > 0 {
			key := strings.TrimSpace(line[:colonIndex])
			value := strings.TrimSpace(line[colonIndex+1:])
			headers[key] = value
		}
	}

	return headers
}

func (s *EmailService) buildInviteLink(userID uuid.UUID, eventID, calendarID string, attendees []models.GoogleCalendarAttendee) string {
	params := url.Values{}
	params.Set("eventId", eventID)
	params.Set("calendarId", calendarID)
	params.Set("uid", userID.String())

	var emails []string
	for _, attendee := range attendees {
		emails = append(emails, attendee.Email)
	}
	params.Set("attendees", strings.Join(emails, ","))

	return fmt.Sprintf("%s/auth/inviteAdditionalAttendees?%s", s.config.APIURL, params.Encode())
}

func (s *EmailService) formatEventDate(eventTime time.Time, timezone string) string {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}

	return eventTime.In(loc).Format("Monday, January 2 at 3:04 PM MST")
}

func (s *EmailService) formatAttendees(attendees []models.GoogleCalendarAttendee) string {
	var emails []string
	for _, attendee := range attendees {
		emails = append(emails, attendee.Email)
	}

	return strings.Join(emails, ", ")
}

func (s *EmailService) sendMultipleEventsResponse(ctx context.Context, to string, webhook *models.EmailWebhook, events []*models.GoogleCalendarEvent, failures []error) error {
	html := fmt.Sprintf("%d events added to your calendar.<br><br>", len(events))

	for _, event := range events {
		html += fmt.Sprintf("<strong>%s</strong><br>", event.Summary)
		html += fmt.Sprintf("Date: %s<br>", s.formatEventDate(event.StartTime, event.TimeZone))
		if event.Location != "" {
			html += fmt.Sprintf("Location: %s<br>", event.Location)
		}
		html += fmt.Sprintf(`<a href="%s" style="display:inline-block; padding:10px 20px; margin:5px 0; background-color:#3498db; color:white; text-align:center; text-decoration:none; font-weight:bold; border-radius:5px;">View Event</a><br><br>`, event.HTMLLink)
	}

	if len(failures) > 0 {
		html += fmt.Sprintf("<p>Failed to add %d event(s). Please try again or contact support.</p>", len(failures))
	}

	html += fmt.Sprintf(`<br><br>You can always ask for help: <a href="mailto:hey@%s">hey@%s</a><br>`, s.config.EmailDomain, s.config.EmailDomain)

	return s.emailProvider.SendEmail(ctx, to, s.config.MainEmailAddress, fmt.Sprintf("Re: %s", webhook.Subject), "", html, s.getThreadHeaders(webhook.Headers))
}

func (s *EmailService) sendEmailResponse(ctx context.Context, to string, webhook *models.EmailWebhook, template templates.EmailTemplate, includeThread bool) error {
	html := template.HTML
	subject := webhook.Subject

	if template.Subject != "" {
		subject = template.Subject
	}

	if includeThread {
		html = s.threadEmailHTML(webhook, html)
	}

	headers := s.getThreadHeaders(webhook.Headers)

	return s.emailProvider.SendEmail(ctx, to, s.config.MainEmailAddress, subject, "", html, headers)
}

func (s *EmailService) threadEmailHTML(original *models.EmailWebhook, responseHTML string) string {
	return fmt.Sprintf("%s%s", responseHTML, original.HTML)
}

func (s *EmailService) getThreadHeaders(headerString string) map[string]string {
	headers := make(map[string]string)
	headerMap := s.parseEmailHeaders(headerString)

	if inReplyTo, exists := headerMap["In-Reply-To"]; exists {
		headers["In-Reply-To"] = inReplyTo
	}

	if references, exists := headerMap["References"]; exists {
		headers["References"] = references
	}

	return headers
}

func (s *EmailService) extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return s.config.EmailDomain
}
