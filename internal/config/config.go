// internal/config/config.go
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/wizenheimer/swiftcal/pkg/logger"
)

type Config struct {
	// Server
	Port        string
	Environment string
	APIURL      string

	// Database
	DatabaseURL      string
	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string

	// Google OAuth2
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	// OpenAI
	OpenAIAPIKey string

	// Email Providers
	MailgunAPIKey    string
	MailgunDomain    string
	MainEmailAddress string

	// Webhook Security
	MailgunWebhookSecret string

	// JWT
	JWTSecret string

	// Domain Configuration
	BaseDomain  string
	AppDomain   string
	EmailDomain string
}

func Load() (*Config, error) {
	// Load .env file in development
	if os.Getenv("ENVIRONMENT") != "production" {
		if err := godotenv.Load(); err != nil {
			logger.GetLogger().Warn("No .env file found")
		}
	}

	config := &Config{
		// Server
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		APIURL:      getEnv("API_URL", "http://localhost:8080"),

		// Database
		DatabaseURL:      getEnv("DATABASE_URL", ""),
		DatabaseHost:     getEnv("DATABASE_HOST", "localhost"),
		DatabasePort:     getEnv("DATABASE_PORT", "5432"),
		DatabaseUser:     getEnv("DATABASE_USER", "username"),
		DatabasePassword: getEnv("DATABASE_PASSWORD", "password"),
		DatabaseName:     getEnv("DATABASE_NAME", "swiftcal"),

		// Google OAuth2
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),

		// OpenAI
		OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),

		// Email Providers
		MailgunAPIKey:    getEnv("MAILGUN_API_KEY", ""),
		MailgunDomain:    getEnv("MAILGUN_DOMAIN", ""),
		MainEmailAddress: getEnv("MAIN_EMAIL_ADDRESS", "swiftcal@"+getEnv("EMAIL_DOMAIN", "swiftcallabs.com")),

		// Webhook Security
		MailgunWebhookSecret: getEnv("MAILGUN_WEBHOOK_SECRET", "incoming"),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", ""),

		// Domain Configuration
		BaseDomain:  getEnv("BASE_DOMAIN", "swiftcallabs.com"),
		AppDomain:   getEnv("APP_DOMAIN", "app.swiftcallabs.com"),
		EmailDomain: getEnv("EMAIL_DOMAIN", "swiftcallabs.com"),
	}

	// Build database URL if not provided
	if config.DatabaseURL == "" {
		config.DatabaseURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			config.DatabaseUser,
			config.DatabasePassword,
			config.DatabaseHost,
			config.DatabasePort,
			config.DatabaseName,
		)
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Validate() error {
	required := map[string]string{
		"GOOGLE_CLIENT_ID":     c.GoogleClientID,
		"GOOGLE_CLIENT_SECRET": c.GoogleClientSecret,
		"GOOGLE_REDIRECT_URL":  c.GoogleRedirectURL,
		"OPENAI_API_KEY":       c.OpenAIAPIKey,
		"JWT_SECRET":           c.JWTSecret,
		"MAIN_EMAIL_ADDRESS":   c.MainEmailAddress,
	}

	for key, value := range required {
		if value == "" {
			return fmt.Errorf("required environment variable %s is not set", key)
		}
	}

	// At least one email provider should be configured
	if c.MailgunAPIKey == "" {
		return fmt.Errorf("mailgun must be configured")
	}

	return nil
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// GetWebURL returns the full web URL for the given path
func (c *Config) GetWebURL(path string) string {
	return "https://www." + c.BaseDomain + path
}

// GetAppURL returns the full app URL for the given path
func (c *Config) GetAppURL(path string) string {
	return "https://" + c.AppDomain + path
}

// GetSupportEmail returns the support email address
func (c *Config) GetSupportEmail() string {
	return "hey@" + c.EmailDomain
}

// GetswiftcalEmail returns the main swiftcal email address
func (c *Config) GetswiftcalEmail() string {
	return "swiftcal@" + c.EmailDomain
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
