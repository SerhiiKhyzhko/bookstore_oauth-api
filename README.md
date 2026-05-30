# Bookstore OAuth API

A Go-based microservice responsible for authentication and token management within the "Bookstore" project ecosystem. It issues and refreshes **stateless JWT tokens**, eliminating the need for any persistent token storage.

In production, other services (Users API, Items API) validate access tokens locally via the **`bookstore-oauth-go` SDK** — they do not call this service over HTTP for verification.

> **Note:** This is version `v2`. The previous `v1.0.0` used MD5 tokens backed by Cassandra. See the `v1.0.0` tag for that implementation.

## Technology Stack

- **Framework:** [Gin](https://github.com/gin-gonic/gin) for HTTP routing.
- **Token Standard:** Stateless JWT via [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt).
- **Signing Algorithm:** HS256 (HMAC-SHA256).
- **HTTP Client:** [Resty](https://github.com/go-resty/resty) for communication with the Users API.
- **Configuration:** [GoDotEnv](https://github.com/joho/godotenv) for local environment variable loading.
- **API Documentation:** [Swaggo](https://github.com/swaggo/swag) with Gin integration.

## Architectural Notes

- **Clean Architecture:** Domain logic is separated from infrastructure. The `domain/token` package defines request/response models; `internal/jwtutil` handles all JWT cryptography.
- **Stateless Tokens:** No token storage on the server side. Token validity is verified entirely by checking the JWT signature and claims — no database lookups required.
- **Two-Token Strategy:** The service issues both an **access token** (short-lived) and a **refresh token** (long-lived). When the access token expires, the client uses the refresh token to obtain a new one without re-authenticating.
- **ACL Pattern:** The `users_client` communicates with the external Users API but returns only a `userId` (`int64`) to the service layer, keeping external models out of the domain.
- **Token Type Enforcement:** Each token carries a `token_type` claim (`"access"` or `"refresh"`). The refresh endpoint explicitly rejects access tokens, and vice versa.
- **Algorithm Enforcement:** Token verification uses `jwt.WithValidMethods([]string{"HS256"})` to prevent algorithm confusion attacks.
- **Production vs development routes:** Only `POST /oauth/create` and `POST /oauth/refresh` are registered when `APP_ENV` is not `development`. Swagger UI and `POST /oauth/verify` are enabled only in development for manual testing and API exploration.
- **Microservice auth:** The `bookstore-oauth-go` SDK validates JWTs in-process (signature, expiry, claims). It does not perform HTTP calls to `/oauth/verify` on this service.

## Prerequisites

- Go 1.18 or newer
- A running instance of the **Bookstore Users API** for credential validation

## Configuration

Create a `.env` file in the project root using the following template:

```env
# Application
APP_ENV=development          # Use "development" for Swagger + /oauth/verify; any other value for production-like routing
GIN_PORT=:8081
CTX_TIMEOUT=2s

# Logger
LEVEL=info
OUTPUT_PATHS=stdout

# Users API
USERS_API_BASE_URL=http://localhost:8080/users/login

# JWT
SIGN_KEY=your-super-secret-key-min-32-characters-long
ACCESS_EXPIRATION=15m         # Optional: Defaults to 15m
REFRESH_EXPIRATION=24h        # Optional: Defaults to 24h
```

> **Security note:** The `SIGN_KEY` must be at least 32 characters long for HS256. Never commit it to version control.

## Getting Started

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd bookstore_oauth-api
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Ensure the Users API is running** and accessible at `USERS_API_BASE_URL`.

4. **Run the application:**
   ```bash
   go run src/main.go
   ```

   Set `APP_ENV=development` in `.env` if you need Swagger and the `/oauth/verify` test endpoint locally.

## API Documentation (Swagger)

Available only when `APP_ENV=development`. Once the service is running:

```
http://localhost:8081/swagger/index.html
```

To regenerate Swagger docs after modifying controller annotations:

```bash
swag init --parseDependency --parseInternal --generalInfo main.go --dir ./src --output ./src/docs
```

## Running Tests

```bash
go test ./...
```

To view test coverage:

```bash
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```

## API Endpoints

### Always available

| Method | Path             | Description                                  | Auth Required |
|--------|------------------|----------------------------------------------|---------------|
| `POST` | `/oauth/create`  | Authenticate user, issue access + refresh token | No         |
| `POST` | `/oauth/refresh` | Issue new access token using a refresh token | No            |

### Development only (`APP_ENV=development`)

| Method | Path             | Description                                  |
|--------|------------------|----------------------------------------------|
| `POST` | `/oauth/verify`  | Decode and validate a JWT; return claims (manual testing) |
| `GET`  | `/swagger/*`     | Swagger UI                                   |

### `POST /oauth/create`

Validates user credentials against the Users API and returns a JWT access token and refresh token pair.

**Request body:**
```json
{
  "user_email": "user@example.com",
  "password": "secret"
}
```

**Success response `201`:**
```json
{
  "access_token": "eyJhbGci...",
  "refresh_token": "eyJhbGci...",
  "user_id": 1
}
```

---

### `POST /oauth/refresh`

Accepts a valid, non-expired refresh token and issues a new access token. Returns `401` if the token is expired, invalid, or is not a refresh token.

**Request body:**
```json
{
  "refresh_token": "eyJhbGci..."
}
```

**Success response `200`:**
```json
{
  "access_token": "eyJhbGci..."
}
```

---

### `POST /oauth/verify` (development only)

> **Not used in production.** Users API and Items API rely on the `bookstore-oauth-go` SDK for JWT validation and do not call this endpoint.

Available when `APP_ENV=development`. Useful for manual checks, Swagger, and debugging token contents without running a consumer service.

The token can be supplied in two ways (checked in this order):

1. **`Authorization` header (preferred):**
   ```http
   Authorization: Bearer eyJhbGci...
   ```
2. **JSON body (fallback)** when the header is missing or empty:
   ```json
   {
     "token": "eyJhbGci..."
   }
   ```

**Example (Bearer):**
```bash
curl -X POST http://localhost:8081/oauth/verify \
  -H "Authorization: Bearer eyJhbGci..."
```

**Example (JSON body):**
```bash
curl -X POST http://localhost:8081/oauth/verify \
  -H "Content-Type: application/json" \
  -d '{"token":"eyJhbGci..."}'
```

**Success response `200`:**
```json
{
  "user_id": 1,
  "token_type": "access",
  "issuer": "bookstore-oauth-api",
  "expires_at": 1718000000
}
```

**Error responses:**

| Status | When |
|--------|------|
| `400` | Token is missing or could not be read from the header or body (`missing or invalid token in Authorization header or request body`). Detailed causes are logged server-side. |
| `401` | Token is present but invalid, expired, or has a bad signature. |
| `500` | Unexpected server error during verification. |

## Token Lifecycle

```
POST /oauth/create
  └─ validates credentials via Users API
  └─ issues access token  (exp: ACCESS_EXPIRATION,  type: "access")
  └─ issues refresh token (exp: REFRESH_EXPIRATION, type: "refresh")
  
Returns `404` if credentials are invalid, `408` if the Users API times out.
Client stores both tokens.

Every protected request:
  └─ sends access token in Authorization: Bearer <token> header
  └─ if 401 → calls POST /oauth/refresh with refresh token
     └─ if refresh valid → new access token issued
     └─ if refresh expired → 401, user must re-authenticate via POST /oauth/create

Microservice token validation (Users API, Items API) — production:
  └─ client sends Authorization: Bearer <access_token>
  └─ bookstore-oauth-go middleware validates JWT locally (same SIGN_KEY / HS256 rules)
  └─ claims (e.g. user_id) attached to request context → handler authorizes or rejects
  └─ no HTTP call to oauth-api /oauth/verify

Development only:
  └─ optional POST /oauth/verify on oauth-api to inspect a token via curl or Swagger
```