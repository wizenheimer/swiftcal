// internal/models/event.go
package models

import (
	"time"
)

type Event struct {
	Summary        string   `json:"summary" validate:"required"`
	Location       *string  `json:"location"`
	Description    *string  `json:"description"`
	ConferenceCall bool     `json:"conference_call"`
	Date           string   `json:"date" validate:"required"`
	StartTime      string   `json:"start_time" validate:"required"`
	EndTime        *string  `json:"end_time"`
	TimeZone       *string  `json:"time_zone"`
	Attendees      []string `json:"attendees" validate:"required"`
}

type EventsResponse struct {
	Events      []Event `json:"events,omitempty"`
	Error       *string `json:"error,omitempty"`
	Description *string `json:"description,omitempty"`
}

type TimezoneResponse struct {
	Reason   string  `json:"reason"`
	Timezone *string `json:"timezone"`
}

type GoogleCalendarEvent struct {
	ID          string                   `json:"id"`
	Summary     string                   `json:"summary"`
	Description string                   `json:"description"`
	Location    string                   `json:"location"`
	StartTime   time.Time                `json:"start_time"`
	EndTime     time.Time                `json:"end_time"`
	TimeZone    string                   `json:"timezone"`
	HTMLLink    string                   `json:"html_link"`
	Attendees   []GoogleCalendarAttendee `json:"attendees"`
}

type GoogleCalendarAttendee struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name,omitempty"`
	Organizer   bool   `json:"organizer,omitempty"`
}

type EmailWebhook struct {
	Subject   string `json:"subject"`
	Text      string `json:"text"`
	HTML      string `json:"html"`
	From      string `json:"from"`
	To        string `json:"to"`
	Headers   string `json:"headers"`
	Envelope  string `json:"envelope"`
	SPF       string `json:"SPF"`
	DKIM      string `json:"dkim"`
	Timestamp string `json:"timestamp,omitempty"`
}

type EmailFile struct {
	Filename string `json:"filename"`
	Content  []byte `json:"content"`
	MimeType string `json:"mime_type"`
}
