// internal/handlers/auth.go
package handlers

import (
	"net/http"

	"github.com/wizenheimer/swiftcal/internal/config"
	"github.com/wizenheimer/swiftcal/internal/services"
	"github.com/wizenheimer/swiftcal/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService *services.AuthService
	config      *config.Config
}

func NewAuthHandler(authService *services.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		config:      cfg,
	}
}

func (h *AuthHandler) Signup(c *fiber.Ctx) error {
	authURL := h.authService.GetAuthURL()
	return c.Redirect(authURL, http.StatusFound)
}

func (h *AuthHandler) Callback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		logger.GetLogger().Error("No authorization code received")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "No authorization code received",
		})
	}

	user, err := h.authService.HandleCallback(c.Context(), code)
	if err != nil {
		logger.GetLogger().Error("OAuth callback failed", zap.Error(err))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Authentication failed",
		})
	}

	logger.GetLogger().Info("User authenticated successfully", zap.String("user_id", user.ID.String()))
	return c.Redirect(h.config.GetWebURL("/thanks"), http.StatusFound)
}

func (h *AuthHandler) VerifyAdditionalEmail(c *fiber.Ctx) error {
	uuidParam := c.Query("uuid")
	if uuidParam == "" {
		return c.Redirect(h.config.GetWebURL("/not-found"), http.StatusFound)
	}

	verificationCode, err := uuid.Parse(uuidParam)
	if err != nil {
		logger.GetLogger().Error("Invalid UUID format", zap.Error(err))
		return c.Redirect(h.config.GetWebURL("/not-found"), http.StatusFound)
	}

	pending, err := h.authService.GetPendingEmailByCode(c.Context(), verificationCode)
	if err != nil {
		logger.GetLogger().Error("Failed to get pending email", zap.Error(err))
		return c.Redirect(h.config.GetWebURL("/not-found"), http.StatusFound)
	}

	// Add email to user account
	if err := h.authService.AddEmailAddress(c.Context(), pending.OwnerUserID, pending.Email, false); err != nil {
		logger.GetLogger().Error("Failed to add email address", zap.Error(err))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add email address",
		})
	}

	// Clean up pending email
	if err := h.authService.DeletePendingEmail(c.Context(), pending.Email); err != nil {
		logger.GetLogger().Warn("Failed to delete pending email", zap.Error(err))
	}

	logger.GetLogger().Info("Email address verified and added",
		zap.String("email", pending.Email),
		zap.String("user_id", pending.OwnerUserID.String()))

	return c.JSON(fiber.Map{
		"message": "Email address added successfully",
		"email":   pending.Email,
	})
}
