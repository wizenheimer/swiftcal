// internal/models/user.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	AccessToken  *string    `json:"-" db:"access_token"`
	RefreshToken *string    `json:"-" db:"refresh_token"`
	ExpiryDate   *time.Time `json:"-" db:"expiry_date"`
	TokenScope   *string    `json:"-" db:"token_scope"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}
