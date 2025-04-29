# API Documentation

## Authentication Endpoints

### POST /api/auth/register

Register a new user account.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "securepassword",
  "name": "John Doe"
}
```

**Response:**

```json
{
  "success": true,
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe"
  },
  "token": "jwt_token_here"
}
```

### POST /api/auth/login

Authenticate user and get access token.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

## Calendar Endpoints

### GET /api/calendar/events

Get all events for the authenticated user.

**Query Parameters:**

- `start_date`: Start date filter (YYYY-MM-DD)
- `end_date`: End date filter (YYYY-MM-DD)

### POST /api/calendar/events

Create a new calendar event.

**Request Body:**

```json
{
  "title": "Team Meeting",
  "description": "Weekly team sync",
  "start_time": "2025-05-01T10:00:00Z",
  "end_time": "2025-05-01T11:00:00Z",
  "location": "Conference Room A"
}
```

## Email Endpoints

### POST /api/email/parse

Parse email content and extract events.

**Request Body:**

```json
{
  "email_content": "Meeting tomorrow at 2 PM in the office"
}
```

**Response:**

```json
{
  "events": [
    {
      "title": "Meeting",
      "start_time": "2025-05-02T14:00:00Z",
      "location": "office"
    }
  ]
}
```
