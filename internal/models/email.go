// internal/models/email.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type EmailAddress struct {
	Email     string    `json:"email" db:"email"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	IsDefault bool      `json:"is_default" db:"is_default"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type PendingEmailAddress struct {
	Email            string    `json:"email" db:"email"`
	OwnerUserID      uuid.UUID `json:"owner_user_id" db:"owner_user_id"`
	OwnerEmail       string    `json:"owner_email" db:"owner_email"`
	VerificationCode uuid.UUID `json:"verification_code" db:"verification_code"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	ExpiresAt        time.Time `json:"expires_at" db:"expires_at"`
}
