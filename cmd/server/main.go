// cmd/server/main.go
package main

import (
	"context"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wizenheimer/swiftcal/internal/config"
	"github.com/wizenheimer/swiftcal/internal/database"
	"github.com/wizenheimer/swiftcal/internal/handlers"
	"github.com/wizenheimer/swiftcal/internal/middleware"
	"github.com/wizenheimer/swiftcal/internal/services"
	"github.com/wizenheimer/swiftcal/pkg/logger"
	"github.com/wizenheimer/swiftcal/templates"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

// Server represents the main application server
type Server struct {
	app          *fiber.App
	config       *config.Config
	db           *database.DB
	authService  *services.AuthService
	emailService *services.EmailService
	cronService  *services.CronService
	shutdownChan chan os.Signal
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config, db *database.DB) *Server {
	// Initialize services
	authService := services.NewAuthService(db, cfg)
	openaiService := services.NewOpenAIService(cfg)
	calendarService := services.NewCalendarService(cfg, authService)
	emailService := services.NewEmailService(cfg, authService, calendarService, openaiService)
	cronService := services.NewCronService(db, cfg, authService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, cfg)
	emailHandler := handlers.NewEmailHandler(emailService, cfg)
	calendarHandler := handlers.NewCalendarHandler(calendarService, authService, cfg)

	// Initialize Fiber app
	app := createFiberApp()

	// Setup middleware
	setupMiddleware(app)

	// Setup routes
	setupRoutes(app, authHandler, emailHandler, calendarHandler, cfg)

	return &Server{
		app:          app,
		config:       cfg,
		db:           db,
		authService:  authService,
		emailService: emailService,
		cronService:  cronService,
		shutdownChan: make(chan os.Signal, 1),
	}
}

// Start starts the server and background tasks
func (s *Server) Start() error {
	// Start background tasks
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.cronService.StartBackgroundJobs(ctx)

	// Get port
	port := s.getPort()
	logger.GetLogger().Info("Server starting", zap.String("port", port))

	// Setup graceful shutdown
	signal.Notify(s.shutdownChan, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		if err := s.app.Listen(":" + port); err != nil {
			logger.GetLogger().Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for shutdown signal
	<-s.shutdownChan
	logger.GetLogger().Info("Shutting down server...")

	// Cancel background tasks
	cancel()

	// Shutdown server gracefully
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := s.app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.GetLogger().Error("Server shutdown error", zap.Error(err))
		return err
	}

	logger.GetLogger().Info("Server stopped")
	return nil
}

// Close closes the server resources
func (s *Server) Close() {
	s.db.Close()
}

func main() {
	// Initialize logger
	logger.Init()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.GetLogger().Fatal("Failed to load configuration", zap.Error(err))
	}

	logger.GetLogger().Info("Starting swiftcal server", zap.String("environment", cfg.Environment))

	// Initialize database
	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		logger.GetLogger().Fatal("Failed to connect to database", zap.Error(err))
	}

	// Create and start server
	server := NewServer(cfg, db)
	defer server.Close()

	if err := server.Start(); err != nil {
		logger.GetLogger().Fatal("Server failed", zap.Error(err))
	}
}

func createFiberApp() *fiber.App {
	return fiber.New(fiber.Config{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		BodyLimit:    10 * 1024 * 1024, // 10MB
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.GetLogger().Error("Fiber error", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		},
	})
}

func setupMiddleware(app *fiber.App) {
	// Recovery middleware
	app.Use(recover.New())

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Rate limiting
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		},
	}))

	// Custom middleware
	app.Use(middleware.RequestLogger())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	})
}

func setupRoutes(app *fiber.App, authHandler *handlers.AuthHandler, emailHandler *handlers.EmailHandler, calendarHandler *handlers.CalendarHandler, cfg *config.Config) {
	// Auth routes
	setupAuthRoutes(app, authHandler, calendarHandler)

	// Webhook routes
	setupWebhookRoutes(app, emailHandler, cfg)

	// Static pages
	setupStaticPages(app, cfg)

	// Catch-all route
	app.Use(func(c *fiber.Ctx) error {
		return c.Redirect("/404", fiber.StatusFound)
	})
}

func setupAuthRoutes(app *fiber.App, authHandler *handlers.AuthHandler, calendarHandler *handlers.CalendarHandler) {
	app.Get("/signup", authHandler.Signup)
	app.Get("/auth/callback", authHandler.Callback)
	app.Get("/auth/verifyAdditionalEmail", authHandler.VerifyAdditionalEmail)
	app.Get("/auth/inviteAdditionalAttendees", calendarHandler.InviteAdditionalAttendees)
}

func setupWebhookRoutes(app *fiber.App, emailHandler *handlers.EmailHandler, cfg *config.Config) {
	if cfg.MailgunWebhookSecret != "" {
		endpoint := "/webhooks/mailgun/" + cfg.MailgunWebhookSecret
		app.Post(endpoint, emailHandler.HandleMailgunWebhook)
		logger.GetLogger().Info("Mailgun webhook endpoint", zap.String("endpoint", endpoint))
	}
}

func setupStaticPages(app *fiber.App, cfg *config.Config) {
	// Home page redirect
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "swiftcal - AI Executive Assistant",
		})
	})

	// Success page
	app.Get("/thanks", func(c *fiber.Ctx) error {
		return c.SendString(templates.GetWelcomePageHTML(cfg.EmailDomain))
	})

	// 404 page
	app.Get("/404", func(c *fiber.Ctx) error {
		return c.Status(404).SendString(templates.Get404PageHTML(cfg.EmailDomain))
	})

	// Invited page
	app.Get("/invited", func(c *fiber.Ctx) error {
		return c.SendString(templates.GetInvitedPageHTML())
	})
}

func (s *Server) getPort() string {
	if s.config.Port != "" {
		return s.config.Port
	}
	return "8080"
}
