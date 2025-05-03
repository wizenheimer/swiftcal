# Changelog

All notable changes to SwiftCal will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial project structure
- Basic documentation
- Docker support
- Database schema

### Changed

- Project name from Crusogo to SwiftCal

### Fixed

- None

## [0.1.0] - 2025-05-11

### Added

- Initial release of SwiftCal
- Core calendar functionality
- Email parsing and integration
- User authentication system
- OpenAI integration for smart event extraction
- RESTful API endpoints
- PostgreSQL database support
- Docker containerization
- Comprehensive documentation
- Testing framework
- Deployment guides

### Features

- **Calendar Management**: Create, read, update, and delete calendar events
- **Email Integration**: Parse emails and extract calendar events automatically
- **User Authentication**: JWT-based authentication with user management
- **Smart Event Extraction**: Use OpenAI to intelligently parse email content
- **API First**: RESTful API for all operations
- **Docker Support**: Easy deployment with Docker and Docker Compose
- **Database**: PostgreSQL with proper schema and migrations

### Technical Stack

- **Backend**: Go 1.21+
- **Database**: PostgreSQL 13+
- **Authentication**: JWT
- **AI Integration**: OpenAI API
- **Email**: SMTP with multi-provider support
- **Containerization**: Docker & Docker Compose
- **Documentation**: Markdown with comprehensive guides

### Documentation

- README with project overview
- Architecture documentation
- API documentation with examples
- Deployment guide for production
- Testing guide with examples
- Contributing guidelines
- This changelog

### Infrastructure

- Docker configuration for development and production
- Database initialization scripts
- Environment configuration management
- Makefile for common operations
- Git hooks for code quality
