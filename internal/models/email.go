// internal/models/email.go
package models

import (
	"time"

	"github.com/google/uuid"
)

// Email represents an email message
type Email struct {
	ID              uuid.UUID         `json:"id" db:"id"`
	UserID          uuid.UUID         `json:"user_id" db:"user_id"`
	MessageID       string            `json:"message_id" db:"message_id"`
	Subject         string            `json:"subject" db:"subject"`
	From            string            `json:"from" db:"from"`
	To              []string          `json:"to" db:"to"`
	Cc              []string          `json:"cc,omitempty" db:"cc"`
	Bcc             []string          `json:"bcc,omitempty" db:"bcc"`
	TextContent     *string           `json:"text_content,omitempty" db:"text_content"`
	HTMLContent     *string           `json:"html_content,omitempty" db:"html_content"`
	Headers         map[string]string `json:"headers,omitempty" db:"headers"`
	Attachments     []Attachment      `json:"attachments,omitempty" db:"attachments"`
	Status          EmailStatus       `json:"status" db:"status"`
	Priority        EmailPriority     `json:"priority" db:"priority"`
	Category        EmailCategory     `json:"category" db:"category"`
	Processed       bool              `json:"processed" db:"processed"`
	ProcessedAt     *time.Time        `json:"processed_at,omitempty" db:"processed_at"`
	ExtractedEvents []uuid.UUID       `json:"extracted_events,omitempty" db:"extracted_events"`
	CreatedAt       time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at" db:"updated_at"`
}

// EmailStatus represents the status of an email
type EmailStatus string

const (
	EmailStatusReceived  EmailStatus = "received"
	EmailStatusProcessed EmailStatus = "processed"
	EmailStatusFailed    EmailStatus = "failed"
	EmailStatusSpam      EmailStatus = "spam"
	EmailStatusArchived  EmailStatus = "archived"
)

// EmailPriority represents the priority of an email
type EmailPriority string

const (
	EmailPriorityLow    EmailPriority = "low"
	EmailPriorityNormal EmailPriority = "normal"
	EmailPriorityHigh   EmailPriority = "high"
	EmailPriorityUrgent EmailPriority = "urgent"
)

// EmailCategory represents the category of an email
type EmailCategory string

const (
	EmailCategoryMeeting      EmailCategory = "meeting"
	EmailCategoryAppointment  EmailCategory = "appointment"
	EmailCategoryReminder     EmailCategory = "reminder"
	EmailCategoryInvitation   EmailCategory = "invitation"
	EmailCategoryUpdate       EmailCategory = "update"
	EmailCategoryCancellation EmailCategory = "cancellation"
	EmailCategoryGeneral      EmailCategory = "general"
	EmailCategorySpam         EmailCategory = "spam"
)

// Attachment represents an email attachment
type Attachment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	EmailID   uuid.UUID `json:"email_id" db:"email_id"`
	Filename  string    `json:"filename" db:"filename"`
	MimeType  string    `json:"mime_type" db:"mime_type"`
	Size      int64     `json:"size" db:"size"`
	Content   []byte    `json:"content,omitempty" db:"content"`
	URL       *string   `json:"url,omitempty" db:"url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// EmailFilter represents filters for querying emails
type EmailFilter struct {
	UserID    *uuid.UUID     `json:"user_id,omitempty"`
	Status    *EmailStatus   `json:"status,omitempty"`
	Priority  *EmailPriority `json:"priority,omitempty"`
	Category  *EmailCategory `json:"category,omitempty"`
	Processed *bool          `json:"processed,omitempty"`
	StartDate *time.Time     `json:"start_date,omitempty"`
	EndDate   *time.Time     `json:"end_date,omitempty"`
	From      *string        `json:"from,omitempty"`
	Subject   *string        `json:"subject,omitempty"`
	Limit     int            `json:"limit,omitempty"`
	Offset    int            `json:"offset,omitempty"`
}

// EmailWebhook represents an incoming email webhook
type EmailWebhook struct {
	Subject     string              `json:"subject"`
	Text        string              `json:"text"`
	HTML        string              `json:"html"`
	From        string              `json:"from"`
	To          string              `json:"to"`
	Headers     map[string]string   `json:"headers"`
	Envelope    string              `json:"envelope"`
	SPF         string              `json:"SPF"`
	DKIM        string              `json:"dkim"`
	Timestamp   string              `json:"timestamp,omitempty"`
	Attachments []WebhookAttachment `json:"attachments,omitempty"`
}

// WebhookAttachment represents an attachment from a webhook
type WebhookAttachment struct {
	Filename string `json:"filename"`
	Content  []byte `json:"content"`
	MimeType string `json:"mime_type"`
}

// EmailProcessingResult represents the result of email processing
type EmailProcessingResult struct {
	EmailID         uuid.UUID `json:"email_id"`
	EventsExtracted int       `json:"events_extracted"`
	Events          []Event   `json:"events,omitempty"`
	Confidence      float64   `json:"confidence"`
	ProcessingTime  float64   `json:"processing_time"`
	Errors          []string  `json:"errors,omitempty"`
}

// EmailTemplate represents an email template
type EmailTemplate struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Subject     string    `json:"subject" db:"subject"`
	HTMLContent string    `json:"html_content" db:"html_content"`
	TextContent string    `json:"text_content" db:"text_content"`
	Category    string    `json:"category" db:"category"`
	Variables   []string  `json:"variables" db:"variables"`
	Active      bool      `json:"active" db:"active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// EmailNotification represents an email notification
type EmailNotification struct {
	ID        uuid.UUID          `json:"id" db:"id"`
	UserID    uuid.UUID          `json:"user_id" db:"user_id"`
	Type      NotificationType   `json:"type" db:"type"`
	Subject   string             `json:"subject" db:"subject"`
	Content   string             `json:"content" db:"content"`
	Recipient string             `json:"recipient" db:"recipient"`
	Status    NotificationStatus `json:"status" db:"status"`
	SentAt    *time.Time         `json:"sent_at,omitempty" db:"sent_at"`
	Error     *string            `json:"error,omitempty" db:"error"`
	CreatedAt time.Time          `json:"created_at" db:"created_at"`
}

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeEventReminder NotificationType = "event_reminder"
	NotificationTypeEventUpdate   NotificationType = "event_update"
	NotificationTypeEventCancel   NotificationType = "event_cancel"
	NotificationTypeSystemAlert   NotificationType = "system_alert"
	NotificationTypeWelcome       NotificationType = "welcome"
)

// NotificationStatus represents the status of a notification
type NotificationStatus string

const (
	NotificationStatusPending NotificationStatus = "pending"
	NotificationStatusSent    NotificationStatus = "sent"
	NotificationStatusFailed  NotificationStatus = "failed"
	NotificationStatusBounced NotificationStatus = "bounced"
)
