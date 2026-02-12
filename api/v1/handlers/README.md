# Handler Layer Documentation

The handler layer provides HTTP endpoints for the RSS Aggregator application. It follows RESTful design principles and integrates with the service layer.

## Structure

```
api/v1/
├── dto/                    # Data Transfer Objects
│   ├── user_dto.go        # User request/response types
│   ├── feed_dto.go        # Feed request/response types
│   ├── feed_follow_dto.go # Feed follow request/response types
│   └── (post DTOs in user_dto.go)
├── handlers/              # HTTP request handlers
│   ├── user_handler.go    # User endpoints
│   ├── feed_handler.go    # Feed endpoints
│   ├── feed_follow_handler.go # Feed follow endpoints
│   ├── post_handler.go    # Post endpoints
│   └── rss_handler.go     # RSS fetching endpoints
└── middleware/
    └── auth.go            # Authentication middleware
```

## Handlers

### UserHandler

**File**: `api/v1/handlers/user_handler.go`

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/v1/users` | No | Create a new user |
| GET | `/v1/users` | Yes | Get current authenticated user |

### FeedHandler

**File**: `api/v1/handlers/feed_handler.go`

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/v1/feeds` | Yes | Create a new feed |
| GET | `/v1/feeds` | No | Get all feeds |
| GET | `/v1/feeds?id={uuid}` | No | Get feed by ID |

### FeedFollowHandler

**File**: `api/v1/handlers/feed_follow_handler.go`

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/v1/feed_follows` | Yes | Follow a feed |
| GET | `/v1/feed_follows` | Yes | Get user's feed follows |
| DELETE | `/v1/feed_follows?id={uuid}` | Yes | Unfollow a feed |

### PostHandler

**File**: `api/v1/handlers/post_handler.go`

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/v1/posts?limit=10&offset=0` | Yes | Get posts for authenticated user |

### RSSHandler

**File**: `api/v1/handlers/rss_handler.go`

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/v1/rss/fetch?feed_id={uuid}` | Yes | Manually fetch a feed |

## Authentication

Authentication uses API keys via the `Authorization` header:

```
Authorization: ApiKey <your_api_key>
```

The auth middleware (`api/v1/middleware/auth.go`) validates the API key and injects the authenticated user into the request context.

## Request/Response Examples

### Create User

```bash
POST /v1/users
Content-Type: application/json

{
  "name": "John Doe"
}
```

Response:

```json
{
  "id": "uuid",
  "created_at": "2026-02-12T10:00:00Z",
  "updated_at": "2026-02-12T10:00:00Z",
  "name": "John Doe",
  "api_key": "generated_api_key"
}
```

### Create Feed

```bash
POST /v1/feeds
Authorization: ApiKey <your_api_key>
Content-Type: application/json

{
  "name": "Tech Blog",
  "url": "https://example.com/feed.xml"
}
```

Response:

```json
{
  "id": "uuid",
  "created_at": "2026-02-12T10:00:00Z",
  "updated_at": "2026-02-12T10:00:00Z",
  "name": "Tech Blog",
  "url": "https://example.com/feed.xml",
  "user_id": "uuid",
  "last_fetched_at": null
}
```

### Follow Feed

```bash
POST /v1/feed_follows
Authorization: ApiKey <your_api_key>
Content-Type: application/json

{
  "feed_id": "uuid"
}
```

### Get Posts

```bash
GET /v1/posts?limit=20&offset=0
Authorization: ApiKey <your_api_key>
```

Response:

```json
[
  {
    "id": "uuid",
    "created_at": "2026-02-12T10:00:00Z",
    "updated_at": "2026-02-12T10:00:00Z",
    "title": "Article Title",
    "description": "Article description...",
    "published_at": "2026-02-12T09:00:00Z",
    "url": "https://example.com/article",
    "feed_id": "uuid"
  }
]
```

## Error Handling

All handlers return errors in a consistent format:

```json
{
  "error": "Error message"
}
```

Common HTTP status codes:

- `200 OK` - Success
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Authenticated but not authorized
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource already exists
- `500 Internal Server Error` - Server error

## Helper Functions

Each handler file can use these helper functions (defined in `user_handler.go`):

- `respondWithJSON(w, code, payload)` - Send JSON response
- `respondWithError(w, code, msg)` - Send error response

## Dependencies

Handlers depend on:

- **Service Layer** - Business logic execution
- **Domain Layer** - Type definitions and validation
- **DTO Package** - Request/response type conversions
- **Auth Package** - API key extraction and validation

## Testing

To test handlers:

1. **Unit Tests**: Mock the service layer
2. **Integration Tests**: Use a test database
3. **API Tests**: Use tools like Postman or `curl`

Example `curl` commands:

```bash
# Create user
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Test User"}'

# Get feeds (no auth required)
curl http://localhost:8080/v1/feeds

# Create feed (auth required)
curl -X POST http://localhost:8080/v1/feeds \
  -H "Authorization: ApiKey YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"name": "My Feed", "url": "https://example.com/feed.xml"}'
```

## Architecture Integration

The handler layer fits into the overall architecture:

```
HTTP Request
     ↓
[Handler Layer] ← YOU ARE HERE
     ↓
[Service Layer]
     ↓
[Repository Layer]
     ↓
[Database]
```

Each handler:

1. Parses and validates HTTP requests
2. Extracts authenticated user (if required)
3. Calls service layer methods
4. Maps domain objects to DTOs
5. Returns JSON responses
