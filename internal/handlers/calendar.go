// internal/handlers/calendar.go
package handlers

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wizenheimer/swiftcal/internal/config"
	"github.com/wizenheimer/swiftcal/internal/services"
	"github.com/wizenheimer/swiftcal/pkg/logger"
	"go.uber.org/zap"
)

type CalendarHandler struct {
	calendarService *services.CalendarService
	authService     *services.AuthService
	config          *config.Config
}

func NewCalendarHandler(calendarService *services.CalendarService, authService *services.AuthService, cfg *config.Config) *CalendarHandler {
	return &CalendarHandler{
		calendarService: calendarService,
		authService:     authService,
		config:          cfg,
	}
}

func (h *CalendarHandler) InviteAdditionalAttendees(c *fiber.Ctx) error {
	userIDStr := c.Query("uid")
	eventID := c.Query("eventId")
	calendarID := c.Query("calendarId")
	attendeesStr := c.Query("attendees")

	if userIDStr == "" || eventID == "" || calendarID == "" || attendeesStr == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required parameters",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Parse attendees (comma-separated emails)
	attendees := strings.Split(attendeesStr, ",")
	var validAttendees []string
	for _, attendee := range attendees {
		email := strings.TrimSpace(attendee)
		if email != "" {
			validAttendees = append(validAttendees, email)
		}
	}

	if len(validAttendees) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "No valid attendees provided",
		})
	}

	// Invite attendees
	err = h.calendarService.InviteAdditionalAttendees(c.Context(), userID, eventID, calendarID, validAttendees)
	if err != nil {
		logger.GetLogger().Error("Failed to invite additional attendees", zap.Error(err))
		return c.Redirect(h.config.GetWebURL("/404"), http.StatusFound)
	}

	// Redirect to the event (we don't have the HTML link here, so redirect to a success page)
	return c.Redirect(h.config.GetWebURL("/invited"), http.StatusFound)
}
