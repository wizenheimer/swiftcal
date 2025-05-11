// internal/utils/email_parser.go
package utils

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/wizenheimer/swiftcal/internal/models"

	"github.com/emersion/go-ical"
	"github.com/emersion/go-message"
	"github.com/emersion/go-message/mail"
)

func ParseICSFile(content []byte) (*models.Event, error) {
	reader := strings.NewReader(string(content))
	cal, err := ical.NewDecoder(reader).Decode()
	if err != nil {
		return nil, fmt.Errorf("failed to decode ICS: %w", err)
	}

	// Find the first event using the Events() method
	events := cal.Events()
	if len(events) == 0 {
		return nil, fmt.Errorf("no event found in ICS file")
	}

	event := events[0]

	// Extract event details
	summary := ""
	if prop := event.Props.Get("SUMMARY"); prop != nil {
		summary = prop.Value
	}

	description := ""
	if prop := event.Props.Get("DESCRIPTION"); prop != nil {
		description = prop.Value
	}

	location := ""
	if prop := event.Props.Get("LOCATION"); prop != nil {
		location = prop.Value
	}

	// Parse start time
	startProp := event.Props.Get("DTSTART")
	if startProp == nil {
		return nil, fmt.Errorf("no start time found")
	}

	startTime, err := time.Parse("20060102T150405Z", startProp.Value)
	if err != nil {
		// Try without timezone
		startTime, err = time.Parse("20060102T150405", startProp.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse start time: %w", err)
		}
	}

	// Parse end time
	var endTime *time.Time
	if endProp := event.Props.Get("DTEND"); endProp != nil {
		if et, err := time.Parse("20060102T150405Z", endProp.Value); err == nil {
			endTime = &et
		} else if et, err := time.Parse("20060102T150405", endProp.Value); err == nil {
			endTime = &et
		}
	}

	// Get timezone
	timezone := "UTC"
	if tzParam := startProp.Params.Get("TZID"); tzParam != "" {
		timezone = tzParam
	}

	// Extract attendees
	var attendees []string
	for _, prop := range event.Props.Values("ATTENDEE") {
		email := strings.TrimPrefix(prop.Value, "mailto:")
		attendees = append(attendees, email)
	}

	// Add organizer
	if orgProp := event.Props.Get("ORGANIZER"); orgProp != nil {
		organizer := strings.TrimPrefix(orgProp.Value, "mailto:")
		attendees = append(attendees, organizer)
	}

	// Format for our event model
	modelEvent := &models.Event{
		Summary:     summary,
		Description: &description,
		Location:    &location,
		Date:        startTime.Format("2 January 2006"),
		StartTime:   startTime.Format("15:04"),
		TimeZone:    &timezone,
		Attendees:   attendees,
	}

	if endTime != nil {
		endTimeStr := endTime.Format("15:04")
		modelEvent.EndTime = &endTimeStr
	}

	return modelEvent, nil
}

func ParseEmailMessage(raw []byte) (*message.Header, string, string, error) {
	reader := strings.NewReader(string(raw))
	entity, err := message.Read(reader)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to parse email: %w", err)
	}

	// Create mail reader
	mailReader := mail.NewReader(entity)

	// Read the first part (should be the main message)
	part, err := mailReader.NextPart()
	if err != nil && err != io.EOF {
		return nil, "", "", fmt.Errorf("failed to read mail part: %w", err)
	}

	var textBody, htmlBody string

	// Read all parts to find text and HTML content
	for part != nil {
		switch h := part.Header.(type) {
		case *mail.InlineHeader:
			contentType := h.Get("Content-Type")
			if strings.Contains(contentType, "text/plain") {
				scanner := bufio.NewScanner(part.Body)
				var lines []string
				for scanner.Scan() {
					lines = append(lines, scanner.Text())
				}
				textBody = strings.Join(lines, "\n")
			} else if strings.Contains(contentType, "text/html") {
				scanner := bufio.NewScanner(part.Body)
				var lines []string
				for scanner.Scan() {
					lines = append(lines, scanner.Text())
				}
				htmlBody = strings.Join(lines, "\n")
			}
		}

		part, err = mailReader.NextPart()
		if err != nil && err != io.EOF {
			return nil, "", "", fmt.Errorf("failed to read next mail part: %w", err)
		}
	}

	// If no HTML body found, use text body
	if htmlBody == "" {
		htmlBody = textBody
	}

	return &entity.Header, textBody, htmlBody, nil
}
