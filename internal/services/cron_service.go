// internal/services/cron_service.go
package services

import (
	"context"
	"time"

	"github.com/wizenheimer/swiftcal/internal/config"
	"github.com/wizenheimer/swiftcal/internal/database"
	"github.com/wizenheimer/swiftcal/pkg/logger"

	"go.uber.org/zap"
)

type CronService struct {
	db          *database.DB
	config      *config.Config
	authService *AuthService
}

func NewCronService(db *database.DB, cfg *config.Config, authService *AuthService) *CronService {
	return &CronService{
		db:          db,
		config:      cfg,
		authService: authService,
	}
}

// StartBackgroundJobs starts all background tasks
func (s *CronService) StartBackgroundJobs(ctx context.Context) {
	// Token refresh job - runs every hour
	go s.runTokenRefreshJob(ctx)

	// Cleanup expired pending emails - runs every 6 hours
	go s.runCleanupJob(ctx)

	logger.GetLogger().Info("Background jobs started")
}

func (s *CronService) runTokenRefreshJob(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	// Run immediately on startup
	s.refreshExpiringTokens(ctx)

	for {
		select {
		case <-ctx.Done():
			logger.GetLogger().Info("Token refresh job stopped")
			return
		case <-ticker.C:
			s.refreshExpiringTokens(ctx)
		}
	}
}

func (s *CronService) runCleanupJob(ctx context.Context) {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	// Run immediately on startup
	s.cleanupExpiredData(ctx)

	for {
		select {
		case <-ctx.Done():
			logger.GetLogger().Info("Cleanup job stopped")
			return
		case <-ticker.C:
			s.cleanupExpiredData(ctx)
		}
	}
}

func (s *CronService) refreshExpiringTokens(ctx context.Context) {
	logger.GetLogger().Debug("Starting token refresh job")

	users, err := s.authService.FindUsersWithExpiringTokens(ctx)
	if err != nil {
		logger.GetLogger().Error("Failed to find users with expiring tokens", zap.Error(err))
		return
	}

	if len(users) == 0 {
		logger.GetLogger().Debug("No users with expiring tokens found")
		return
	}

	logger.GetLogger().Info("Refreshing tokens for users with expiring tokens", zap.Int("count", len(users)))

	successCount := 0
	for _, user := range users {
		refreshCtx, cancel := context.WithTimeout(ctx, 30*time.Second)

		if _, err := s.authService.RefreshAccessToken(refreshCtx, user.ID); err != nil {
			logger.GetLogger().Error("Failed to refresh token", zap.Error(err), zap.String("user_id", user.ID.String()))
		} else {
			successCount++
		}

		cancel()
	}

	logger.GetLogger().Info("Token refresh job completed",
		zap.Int("total", len(users)),
		zap.Int("succeeded", successCount),
		zap.Int("failed", len(users)-successCount),
	)
}

func (s *CronService) cleanupExpiredData(ctx context.Context) {
	logger.GetLogger().Debug("Starting cleanup job")

	// Clean up expired pending email addresses
	query := `DELETE FROM pending_email_addresses WHERE expires_at < NOW()`
	result, err := s.db.Pool.Exec(ctx, query)
	if err != nil {
		logger.GetLogger().Error("Failed to cleanup expired pending emails", zap.Error(err))
		return
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected > 0 {
		logger.GetLogger().Info("Cleaned up expired pending email addresses", zap.Int64("count", rowsAffected))
	} else {
		logger.GetLogger().Debug("No expired pending email addresses to clean up")
	}
}
