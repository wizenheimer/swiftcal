// internal/models/event.go
package models

import (
	"time"

	"github.com/google/uuid"
)

// Event represents a calendar event
type Event struct {
	ID             uuid.UUID     `json:"id" db:"id"`
	UserID         uuid.UUID     `json:"user_id" db:"user_id"`
	Title          string        `json:"title" db:"title"`
	Description    *string       `json:"description,omitempty" db:"description"`
	StartTime      time.Time     `json:"start_time" db:"start_time"`
	EndTime        time.Time     `json:"end_time" db:"end_time"`
	Location       *string       `json:"location,omitempty" db:"location"`
	AllDay         bool          `json:"all_day" db:"all_day"`
	Recurring      bool          `json:"recurring" db:"recurring"`
	RecurrenceRule *string       `json:"recurrence_rule,omitempty" db:"recurrence_rule"`
	Status         EventStatus   `json:"status" db:"status"`
	Priority       EventPriority `json:"priority" db:"priority"`
	Color          *string       `json:"color,omitempty" db:"color"`
	Tags           []string      `json:"tags,omitempty" db:"tags"`
	Attendees      []Attendee    `json:"attendees,omitempty" db:"attendees"`
	Reminders      []Reminder    `json:"reminders,omitempty" db:"reminders"`
	Source         EventSource   `json:"source" db:"source"`
	SourceID       *string       `json:"source_id,omitempty" db:"source_id"`
	CreatedAt      time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" db:"updated_at"`
}

// EventStatus represents the status of an event
type EventStatus string

const (
	EventStatusConfirmed EventStatus = "confirmed"
	EventStatusTentative EventStatus = "tentative"
	EventStatusCancelled EventStatus = "cancelled"
	EventStatusPending   EventStatus = "pending"
)

// EventPriority represents the priority of an event
type EventPriority string

const (
	EventPriorityLow    EventPriority = "low"
	EventPriorityNormal EventPriority = "normal"
	EventPriorityHigh   EventPriority = "high"
	EventPriorityUrgent EventPriority = "urgent"
)

// EventSource represents the source of an event
type EventSource string

const (
	EventSourceManual   EventSource = "manual"
	EventSourceEmail    EventSource = "email"
	EventSourceCalendar EventSource = "calendar"
	EventSourceAI       EventSource = "ai"
	EventSourceImport   EventSource = "import"
)

// Attendee represents an event attendee
type Attendee struct {
	ID        uuid.UUID        `json:"id" db:"id"`
	EventID   uuid.UUID        `json:"event_id" db:"event_id"`
	Email     string           `json:"email" db:"email"`
	Name      *string          `json:"name,omitempty" db:"name"`
	Response  AttendeeResponse `json:"response" db:"response"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
}

// AttendeeResponse represents the response of an attendee
type AttendeeResponse string

const (
	AttendeeResponseAccepted  AttendeeResponse = "accepted"
	AttendeeResponseDeclined  AttendeeResponse = "declined"
	AttendeeResponseTentative AttendeeResponse = "tentative"
	AttendeeResponsePending   AttendeeResponse = "pending"
)

// Reminder represents an event reminder
type Reminder struct {
	ID        uuid.UUID    `json:"id" db:"id"`
	EventID   uuid.UUID    `json:"event_id" db:"event_id"`
	Type      ReminderType `json:"type" db:"type"`
	Minutes   int          `json:"minutes" db:"minutes"`
	Sent      bool         `json:"sent" db:"sent"`
	SentAt    *time.Time   `json:"sent_at,omitempty" db:"sent_at"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
}

// ReminderType represents the type of reminder
type ReminderType string

const (
	ReminderTypeEmail ReminderType = "email"
	ReminderTypePush  ReminderType = "push"
	ReminderTypeSMS   ReminderType = "sms"
	ReminderTypeInApp ReminderType = "in_app"
)

// EventFilter represents filters for querying events
type EventFilter struct {
	UserID    *uuid.UUID     `json:"user_id,omitempty"`
	StartDate *time.Time     `json:"start_date,omitempty"`
	EndDate   *time.Time     `json:"end_date,omitempty"`
	Status    *EventStatus   `json:"status,omitempty"`
	Priority  *EventPriority `json:"priority,omitempty"`
	Source    *EventSource   `json:"source,omitempty"`
	Tags      []string       `json:"tags,omitempty"`
	Limit     int            `json:"limit,omitempty"`
	Offset    int            `json:"offset,omitempty"`
}

// EventCreateRequest represents a request to create an event
type EventCreateRequest struct {
	Title          string           `json:"title" binding:"required"`
	Description    *string          `json:"description,omitempty"`
	StartTime      time.Time        `json:"start_time" binding:"required"`
	EndTime        time.Time        `json:"end_time" binding:"required"`
	Location       *string          `json:"location,omitempty"`
	AllDay         bool             `json:"all_day"`
	Recurring      bool             `json:"recurring"`
	RecurrenceRule *string          `json:"recurrence_rule,omitempty"`
	Priority       EventPriority    `json:"priority"`
	Color          *string          `json:"color,omitempty"`
	Tags           []string         `json:"tags,omitempty"`
	Attendees      []string         `json:"attendees,omitempty"`
	Reminders      []ReminderCreate `json:"reminders,omitempty"`
}

// EventUpdateRequest represents a request to update an event
type EventUpdateRequest struct {
	Title          *string        `json:"title,omitempty"`
	Description    *string        `json:"description,omitempty"`
	StartTime      *time.Time     `json:"start_time,omitempty"`
	EndTime        *time.Time     `json:"end_time,omitempty"`
	Location       *string        `json:"location,omitempty"`
	AllDay         *bool          `json:"all_day,omitempty"`
	Recurring      *bool          `json:"recurring,omitempty"`
	RecurrenceRule *string        `json:"recurrence_rule,omitempty"`
	Status         *EventStatus   `json:"status,omitempty"`
	Priority       *EventPriority `json:"priority,omitempty"`
	Color          *string        `json:"color,omitempty"`
	Tags           []string       `json:"tags,omitempty"`
}

// ReminderCreate represents a request to create a reminder
type ReminderCreate struct {
	Type    ReminderType `json:"type" binding:"required"`
	Minutes int          `json:"minutes" binding:"required"`
}

// EventConflict represents a time conflict between events
type EventConflict struct {
	EventID      uuid.UUID    `json:"event_id"`
	Title        string       `json:"title"`
	StartTime    time.Time    `json:"start_time"`
	EndTime      time.Time    `json:"end_time"`
	ConflictType ConflictType `json:"conflict_type"`
}

// ConflictType represents the type of conflict
type ConflictType string

const (
	ConflictTypeOverlap  ConflictType = "overlap"
	ConflictTypeAdjacent ConflictType = "adjacent"
	ConflictTypeTravel   ConflictType = "travel"
)
