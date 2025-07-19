<div align="center">
  <h1>SwiftCal</h1>
  <p>Turn Emails Into Calendar Events</p>
  <p>
    <a href="https://github.com/wizenheimer/swiftcal/tree/main/docs"><strong>Explore the docs »</strong></a>
  </p>
</div>

SwiftCal lets you forward emails to a special address and have them automatically added to your Google Calendar.

No forms, no UI, no clicking around — just forward an email and you're done.

## What It Does

- Takes any email you forward
- Pulls out the date, time, location, and people
- Adds it to your Google Calendar

## Who It's For

Anyone who:

- Gets meeting requests or invites by email
- Forwards things to themselves to remember later
- Doesn’t want to manually create calendar events

## How It Works

1. You send or forward an email to `swiftcal@yourdomain.com`
2. SwiftCal checks the email content
3. If it finds a date/time, it adds an event to your Google Calendar

That’s it.

---

## Quick Start

### Requirements

- Go 1.24.2 or later
- PostgreSQL 15+
- OpenAI account (for text processing)
- Google Cloud project (for Calendar access)
- Mailgun account (to receive forwarded emails)

### Using Docker (recommended)

```bash
git clone https://github.com/wizenheimer/swiftcal.git
cd swiftcal
cp .env.example .env
# Edit .env with your secrets
make up
```

Open your browser to: [http://localhost:8081](http://localhost:8081)

---

## Configuration

Create a `.env` file like this:

```env
# Server
PORT=8080
ENVIRONMENT=development

# Database
DATABASE_URL=postgres://swiftcal:password@localhost:5432/swiftcal?sslmode=disable

# Google Calendar
GOOGLE_CLIENT_ID=your_client_id
GOOGLE_CLIENT_SECRET=your_client_secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/callback

# OpenAI
OPENAI_API_KEY=your_openai_key

# Mailgun
MAILGUN_API_KEY=your_mailgun_key
MAILGUN_DOMAIN=your_domain.com
MAILGUN_WEBHOOK_SECRET=your_webhook_secret

# App Email
MAIN_EMAIL_ADDRESS=swiftcal@your_domain.com

# Auth
JWT_SECRET=your_jwt_secret
```

---

## API Overview

These are the main endpoints:

- `GET /signup` – Starts Google Calendar setup
- `GET /auth/callback` – Handles OAuth return
- `POST /webhooks/mailgun/{secret}` – Handles forwarded emails from Mailgun

You can also configure multiple email addresses and invite attendees via links.

---

## Make Commands

```bash
make help         # List all commands
make up           # Start development environment
make down         # Stop containers
make test         # Run tests
make deps         # Install Go dependencies
make lint         # Run linter
make fmt          # Format code
```

---

## Folder Structure

```
swiftcal/
├── cmd/server/          # Main entry point
├── internal/            # App logic
│   ├── config/          # Config loading
│   ├── handlers/        # HTTP routes
│   ├── services/        # Core functions
│   ├── database/        # DB access
│   └── models/          # Structs and data models
├── pkg/                 # Shared utilities
├── templates/           # Email + calendar prompt templates
├── docs/                # Developer docs
└── docker-compose.yml   # Local dev environment
```

---

## Testing

```bash
make test
```

---

## Deployment

To build and run in production:

```bash
docker build -t swiftcal .
docker run -d --env-file .env -p 8080:8080 swiftcal
```

Or use Docker Compose:

```bash
docker compose up -d
```

---

## Support & Contributions

- Found a bug? [Open an issue](https://github.com/wizenheimer/swiftcal/issues)
- Want a feature? Fork the repo and send a pull request
- Contact: [hey@crusolabs.com](mailto:hey@crusolabs.com)

---

## About Cruso Labs

This project was originally built as an early MVP for [Cruso Labs](https://crusolabs.com), an AI assistant you can email. We're open-sourcing this early version.

## License

MIT – See [LICENSE](LICENSE)
