# Deployment Guide

## Prerequisites

- Docker and Docker Compose
- PostgreSQL database
- SMTP server credentials
- OpenAI API key

## Environment Variables

Create a `.env` file with the following variables:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=swiftcal
DB_USER=postgres
DB_PASSWORD=your_password

# JWT Secret
JWT_SECRET=your_jwt_secret_key

# OpenAI
OPENAI_API_KEY=your_openai_api_key

# SMTP Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your_email@gmail.com
SMTP_PASSWORD=your_app_password

# Server Configuration
SERVER_PORT=8080
ENVIRONMENT=production
```

## Docker Deployment

### Build and Run

```bash
# Build the application
make docker-build

# Run with Docker Compose
make docker-run
```

### Production Deployment

1. Set up a production PostgreSQL database
2. Configure environment variables for production
3. Build and push Docker image to registry
4. Deploy using Docker Compose or Kubernetes

## Manual Deployment

### Database Setup

```sql
-- Run the initialization script
\i scripts/init.sql
```

### Application Setup

```bash
# Install dependencies
go mod download

# Build the application
make build

# Run the server
make run
```

## Monitoring and Logging

- Application logs are written to stdout/stderr
- Use Docker logs or log aggregation service
- Monitor database connections and API response times
- Set up health checks for `/health` endpoint

## Security Considerations

- Use HTTPS in production
- Rotate JWT secrets regularly
- Implement rate limiting
- Use environment-specific configurations
- Regular security updates for dependencies
