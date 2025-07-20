# Makefile
.PHONY: help up down logs test clean deps lint fmt init ngrok

# Color definitions
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
MAGENTA := \033[35m
CYAN := \033[36m
RED := \033[31m
BOLD := \033[1m
RESET := \033[0m

# Default target
help:
	@echo "$(BOLD)$(CYAN)"
	@echo "███████╗██╗    ██╗██╗███████╗████████╗ ██████╗ █████╗ ██╗     "
	@echo "██╔════╝██║    ██║██║██╔════╝╚══██╔══╝██╔════╝██╔══██╗██║     "
	@echo "███████╗██║ █╗ ██║██║███████╗   ██║   ██║     ███████║██║     "
	@echo "╚════██║██║███╗██║██║██║        ██║   ██║     ██╔══██║██║     "
	@echo "███████║╚███╔███╔╝██║██║        ██║   ╚██████╗██║  ██║███████╗"
	@echo "╚══════╝ ╚══╝╚══╝ ╚═╝╚═╝        ╚═╝    ╚═════╝╚═╝  ╚═╝╚══════╝"
	@echo "$(RESET)"
	@echo "$(BOLD)$(CYAN)Development Environment$(RESET)"
	@echo ""
	@echo "$(BOLD)Available targets:$(RESET)"
	@echo "  $(GREEN)up$(RESET)            Start development with hot reload"
	@echo "  $(RED)down$(RESET)          Stop all services"
	@echo "  $(BLUE)logs$(RESET)          Follow application logs"
	@echo "  $(YELLOW)test$(RESET)          Run tests"
	@echo "  $(MAGENTA)clean$(RESET)         Clean build artifacts and containers"
	@echo "  $(CYAN)deps$(RESET)          Download dependencies"
	@echo "  $(GREEN)lint$(RESET)          Run linter"
	@echo "  $(BLUE)fmt$(RESET)           Format code"
	@echo "  $(YELLOW)init$(RESET)          Initialize development tools"
	@echo "  $(MAGENTA)tunnel-start$(RESET)         Start ngrok tunnel"
	@echo "  $(MAGENTA)tunnel-stop$(RESET)         Stop ngrok tunnel"

# Start development environment with hot reload
up:
	@echo "$(BOLD)$(GREEN)Starting development environment...$(RESET)"
	docker compose --profile dev up -d
	@echo "$(YELLOW)Waiting for database to be ready...$(RESET)"
	@until docker compose exec -T postgres pg_isready -U swiftcal; do \
		echo "$(YELLOW)Waiting for postgres...$(RESET)"; \
		sleep 2; \
	done
	@echo "$(BOLD)$(GREEN)Development environment ready!$(RESET)"
	@echo "$(CYAN)App running at http://localhost:8081$(RESET)"
	@echo "$(CYAN)PostgreSQL available at localhost:5432$(RESET)"
	@echo "$(GREEN)Use 'make logs' to follow application logs$(RESET)"
	@echo "$(MAGENTA)Use 'make ngrok' to start ngrok tunnel$(RESET)"

# Start ngrok tunnel
tunnel-start:
	@echo "$(BOLD)$(MAGENTA)Starting ngrok tunnel...$(RESET)"
	@ngrok http 8081 --log=stdout &
	@sleep 3
	@echo "$(GREEN)Ngrok tunnel started!$(RESET)"
	@echo "$(CYAN)Check ngrok dashboard at http://localhost:4040$(RESET)"

# Stop ngrok tunnel
tunnel-stop:
	@echo "$(BOLD)$(MAGENTA)Stopping ngrok tunnel...$(RESET)"
	@pkill -f "ngrok http" || true
	@echo "$(GREEN)Ngrok tunnel stopped!$(RESET)"

# Stop all services
down:
	@echo "$(BOLD)$(RED)Stopping all services...$(RESET)"
	docker compose --profile dev down

# Follow application logs
logs:
	@echo "$(BOLD)$(BLUE)Following application logs...$(RESET)"
	docker compose logs -f app-dev

# Run tests
test:
	@echo "$(BOLD)$(YELLOW)Running tests...$(RESET)"
	go test -v ./...

# Clean everything
clean:
	@echo "$(BOLD)$(MAGENTA)Cleaning build artifacts and containers...$(RESET)"
	rm -rf bin/
	go clean
	@pkill -f "ngrok http" || true
	docker compose --profile dev down -v --rmi local --remove-orphans
	@echo "$(GREEN)Cleanup completed!$(RESET)"

# Download dependencies
deps:
	@echo "$(BOLD)$(CYAN)Downloading dependencies...$(RESET)"
	go mod download
	go mod tidy
	@echo "$(GREEN)Dependencies updated!$(RESET)"

# Run linter
lint:
	@echo "$(BOLD)$(GREEN)Running linter...$(RESET)"
	golangci-lint run

# Format code
fmt:
	@echo "$(BOLD)$(BLUE)Formatting code...$(RESET)"
	go fmt ./...
	goimports -w .
	@echo "$(GREEN)Code formatting completed!$(RESET)"

# Initialize development tools
init:
	@echo "$(BOLD)$(YELLOW)Installing development tools...$(RESET)"
	go install github.com/air-verse/air@latest
	go install golang.org/x/tools/cmd/goimports@latest
	brew install ngrok
	@echo "$(BOLD)$(GREEN)Development tools installed!$(RESET)"
