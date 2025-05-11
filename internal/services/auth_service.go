// internal/services/auth_service.go
package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/wizenheimer/swiftcal/internal/config"
	"github.com/wizenheimer/swiftcal/internal/database"
	"github.com/wizenheimer/swiftcal/internal/models"
	"github.com/wizenheimer/swiftcal/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	googleoauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type AuthService struct {
	db          *database.DB
	config      *config.Config
	oauthConfig *oauth2.Config
}

func NewAuthService(db *database.DB, cfg *config.Config) *AuthService {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes: []string{
			calendar.CalendarScope,
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}

	return &AuthService{
		db:          db,
		config:      cfg,
		oauthConfig: oauthConfig,
	}
}

func (s *AuthService) GetAuthURL() string {
	return s.oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

func (s *AuthService) HandleCallback(ctx context.Context, code string) (*models.User, error) {
	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user info from Google using the official OAuth2 service
	oauth2Service, err := googleoauth2.NewService(ctx, option.WithTokenSource(s.oauthConfig.TokenSource(ctx, token)))
	if err != nil {
		return nil, fmt.Errorf("failed to create oauth2 service: %w", err)
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Check if user exists
	user, err := s.GetUserByEmail(ctx, userInfo.Email)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	if user == nil {
		// Create new user
		user, err = s.CreateUser(ctx, userInfo.Email, token)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		// Add default email address
		if err := s.AddEmailAddress(ctx, user.ID, userInfo.Email, true); err != nil {
			logger.GetLogger().Error("Failed to add default email address", zap.Error(err))
		}

		logger.GetLogger().Info("New user created", zap.String("user_id", user.ID.String()))
	} else {
		// Update existing user's tokens
		if err := s.UpdateUserTokens(ctx, user.ID, token); err != nil {
			return nil, fmt.Errorf("failed to update user tokens: %w", err)
		}
	}

	return user, nil
}

func (s *AuthService) CreateUser(ctx context.Context, email string, token *oauth2.Token) (*models.User, error) {
	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		AccessToken:  &token.AccessToken,
		RefreshToken: &token.RefreshToken,
		ExpiryDate:   &token.Expiry,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if len(token.Extra("scope").(string)) > 0 {
		scope := token.Extra("scope").(string)
		user.TokenScope = &scope
	}

	query := `
		INSERT INTO users (id, email, access_token, refresh_token, expiry_date, token_scope, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	err := s.db.Pool.QueryRow(ctx, query,
		user.ID, user.Email, user.AccessToken, user.RefreshToken,
		user.ExpiryDate, user.TokenScope, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT u.id, u.email, u.access_token, u.refresh_token, u.expiry_date, u.token_scope, u.created_at, u.updated_at
		FROM users u
		JOIN email_addresses ea ON u.id = ea.user_id
		WHERE ea.email = $1
	`

	user := &models.User{}
	err := s.db.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.AccessToken, &user.RefreshToken,
		&user.ExpiryDate, &user.TokenScope, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, access_token, refresh_token, expiry_date, token_scope, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := s.db.Pool.QueryRow(ctx, query, userID).Scan(
		&user.ID, &user.Email, &user.AccessToken, &user.RefreshToken,
		&user.ExpiryDate, &user.TokenScope, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) UpdateUserTokens(ctx context.Context, userID uuid.UUID, token *oauth2.Token) error {
	query := `
		UPDATE users
		SET access_token = $1, refresh_token = $2, expiry_date = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := s.db.Pool.Exec(ctx, query,
		token.AccessToken, token.RefreshToken, token.Expiry, time.Now(), userID,
	)

	return err
}

func (s *AuthService) RefreshAccessToken(ctx context.Context, userID uuid.UUID) (*oauth2.Token, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user.RefreshToken == nil {
		return nil, fmt.Errorf("no refresh token available")
	}

	token := &oauth2.Token{
		AccessToken:  *user.AccessToken,
		RefreshToken: *user.RefreshToken,
		Expiry:       *user.ExpiryDate,
	}

	tokenSource := s.oauthConfig.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	// Update database with new tokens
	if err := s.UpdateUserTokens(ctx, userID, newToken); err != nil {
		return nil, fmt.Errorf("failed to update tokens in database: %w", err)
	}

	logger.GetLogger().Info("Access token refreshed", zap.String("user_id", userID.String()))
	return newToken, nil
}

func (s *AuthService) GetOAuthClient(ctx context.Context, userID uuid.UUID) (*http.Client, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	token := &oauth2.Token{
		AccessToken:  *user.AccessToken,
		RefreshToken: *user.RefreshToken,
		Expiry:       *user.ExpiryDate,
	}

	// Check if token needs refresh
	if token.Expiry.Before(time.Now().Add(5 * time.Minute)) {
		refreshedToken, err := s.RefreshAccessToken(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}
		token = refreshedToken
	}

	return s.oauthConfig.Client(ctx, token), nil
}

func (s *AuthService) AddEmailAddress(ctx context.Context, userID uuid.UUID, email string, isDefault bool) error {
	query := `
		INSERT INTO email_addresses (email, user_id, is_default, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email) DO UPDATE SET
			user_id = EXCLUDED.user_id,
			is_default = EXCLUDED.is_default
	`

	_, err := s.db.Pool.Exec(ctx, query, email, userID, isDefault, time.Now())
	return err
}

func (s *AuthService) RemoveEmailAddress(ctx context.Context, email string) error {
	query := `DELETE FROM email_addresses WHERE email = $1`
	_, err := s.db.Pool.Exec(ctx, query, email)
	return err
}

func (s *AuthService) AddPendingEmailAddress(ctx context.Context, userID uuid.UUID, ownerEmail, pendingEmail string) (uuid.UUID, error) {
	verificationCode := uuid.New()
	query := `
		INSERT INTO pending_email_addresses (email, owner_user_id, owner_email, verification_code, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (email) DO UPDATE SET
			verification_code = EXCLUDED.verification_code,
			created_at = EXCLUDED.created_at,
			expires_at = EXCLUDED.expires_at
	`

	_, err := s.db.Pool.Exec(ctx, query,
		pendingEmail, userID, ownerEmail, verificationCode,
		time.Now(), time.Now().Add(24*time.Hour),
	)

	return verificationCode, err
}

func (s *AuthService) GetPendingEmailByCode(ctx context.Context, code uuid.UUID) (*models.PendingEmailAddress, error) {
	query := `
		SELECT email, owner_user_id, owner_email, verification_code, created_at, expires_at
		FROM pending_email_addresses
		WHERE verification_code = $1 AND expires_at > NOW()
	`

	pending := &models.PendingEmailAddress{}
	err := s.db.Pool.QueryRow(ctx, query, code).Scan(
		&pending.Email, &pending.OwnerUserID, &pending.OwnerEmail,
		&pending.VerificationCode, &pending.CreatedAt, &pending.ExpiresAt,
	)

	if err != nil {
		return nil, err
	}

	return pending, nil
}

func (s *AuthService) DeletePendingEmail(ctx context.Context, email string) error {
	query := `DELETE FROM pending_email_addresses WHERE email = $1`
	_, err := s.db.Pool.Exec(ctx, query, email)
	return err
}

func (s *AuthService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			logger.GetLogger().Error("Failed to rollback transaction", zap.Error(rollbackErr))
		}
	}()

	// Delete email addresses
	_, err = tx.Exec(ctx, `DELETE FROM email_addresses WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("failed to delete email addresses: %w", err)
	}

	// Delete pending email addresses
	_, err = tx.Exec(ctx, `DELETE FROM pending_email_addresses WHERE owner_user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("failed to delete pending email addresses: %w", err)
	}

	// Delete user
	_, err = tx.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *AuthService) FindUsersWithExpiringTokens(ctx context.Context) ([]*models.User, error) {
	query := `
		SELECT id, email, access_token, refresh_token, expiry_date, token_scope, created_at, updated_at
		FROM users
		WHERE expiry_date <= $1
	`

	twoHoursLater := time.Now().Add(2 * time.Hour)
	rows, err := s.db.Pool.Query(ctx, query, twoHoursLater)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.AccessToken, &user.RefreshToken,
			&user.ExpiryDate, &user.TokenScope, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
