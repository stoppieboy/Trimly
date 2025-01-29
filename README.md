# Trimly

A URL shortening service using Golang, Go-Fiber, Redis and containerized with docker

The API service Implements Rate limiting, Input Sanitization, and Custom Short URLs.

These shortened URLs are especially valuable on platforms with character limits, like Twitter, or in marketing campaigns where a clean, branded link is more appealing.

## API Documentation

### Shorten URL

**URL** : `/api/v1/`

**Method** : `POST`

**Data Params** :

| Name | Type | Description |
| --- | --- | --- |
| url | string | The URL to be shortened |
| short | string | The custom short URL |
| expiry | int | The expiry time in hours |

**Success Status Code** : `201 Created`

**Success Response** :

```json
{
    "url": "https://www.google.com",
    "short": "https://trimly.com/abl3i4",
    "expiry": 24,
    "rate_limit": 10,
    "rate_limit_reset": 30
}
```

**Error Response** :

```json
{
    "error": "Invalid URL"
}
```
```json
{
    "error": "Domain error"
}
```
```json
{
    "error": "Rate limit exceeded"
}
```

### Resolve URL

**URL** : `/:url`

**Method** : `GET`

**Success Status Code** : `301 Moved Permanently`

**Success Response** :

```json
{
    "url": "https://www.google.com"
}
```

**Error Response** :

```json
{
    "error": "short not found in the database"
}
```
