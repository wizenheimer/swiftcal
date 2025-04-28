# Architecture Overview

## System Design

SwiftCal follows a clean architecture pattern with the following layers:

### Core Layers

- **Models**: Data structures and domain entities
- **Services**: Business logic and external integrations
- **Handlers**: HTTP request/response handling
- **Middleware**: Cross-cutting concerns

### Key Components

#### Authentication Service

- JWT-based authentication
- User session management
- Role-based access control

#### Calendar Service

- Event creation and management
- Calendar integration
- Scheduling algorithms

#### Email Service

- Email parsing and processing
- Template rendering
- Multi-provider support

#### OpenAI Integration

- Natural language processing
- Email content analysis
- Smart event extraction

## Database Schema

The system uses PostgreSQL with the following main tables:

- users
- events
- emails
- sessions

## API Design

RESTful API with the following endpoints:

- `/api/auth/*` - Authentication endpoints
- `/api/calendar/*` - Calendar management
- `/api/email/*` - Email operations
- `/api/events/*` - Event CRUD operations
