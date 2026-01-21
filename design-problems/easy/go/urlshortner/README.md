
# URL Shortener

A simple URL shortening service built with Go that allows you to create short URLs, redirect to original URLs, and track analytics.

## Features

- **Shorten URLs**: Convert long URLs into short codes
- **Redirect**: Access original URLs using short codes
- **Analytics**: Track visit counts and creation dates for shortened URLs

## API Endpoints

### POST /shorten
Shorten a URL.

**Request:**
```json
{
    "url": "https://example.com/very/long/url"
}
```

**Response:**
```json
{
    "shortened_url": "https://example.com/very/long/url",
    "short_code": "abc12345",
    "redirect_url": "http://localhost:8080/abc12345"
}
```

### GET /{shortCode}
Redirect to the original URL.

### GET /analytics/{shortCode}
Get analytics for a shortened URL.

**Response:**
```json
{
    "url": "https://example.com/very/long/url",
    "visit_count": 42,
    "created_at": "2024-01-15T10:30:00Z"
}
```

## Running the Server

```bash
go run main.go
```

The server starts on `http://localhost:8080`
