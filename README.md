# Bookstore OAuth API

A Go-based microservice responsible for authentication and access token management within the "Bookstore" project ecosystem. It issues, retrieves, and manages access tokens backed by Cassandra as the primary data store.

> **Note:** This is version `v1.0.0`. A future `v2` release will migrate token generation to stateless JWT, removing the Cassandra dependency.

## Technology Stack

- **Framework:** [Gin](https://github.com/gin-gonic/gin) for HTTP routing.
- **Database:** [Cassandra](https://cassandra.apache.org/) via [gocql](https://github.com/gocql/gocql) for access token persistence.
- **HTTP Client:** [Resty](https://github.com/go-resty/resty) for communication with the Users API.
- **Token Generation:** MD5-based token generation (to be replaced with JWT in v2).
- **Configuration:** [GoDotEnv](https://github.com/joho/godotenv) for local environment variable loading.
- **API Documentation:** [Swaggo](https://github.com/swaggo/swag) with Gin integration.

## Architectural Notes

- **Clean Architecture:** Domain logic is separated from infrastructure. The `domain` package defines the `AccessToken` model and repository interfaces; the `repository/db` package provides the Cassandra implementation.
- **ACL Pattern:** The `users_client` communicates with the external Users API but returns only a `userId` (`int64`) to the service layer, keeping external models out of the domain.
- **Grant Types:** The `Create` endpoint currently supports the `password` grant type. Support for `client_credentials` is planned.
- **Swagger UI** is available at `/swagger/index.html` when the service is running.

## Prerequisites

- Go 1.18 or newer
- A running **Cassandra** instance with the required keyspace and table.
- A running instance of the **Bookstore Users API** for credential validation.

## Cassandra Schema

Before starting the service, ensure the following keyspace and table exist in Cassandra:

```cql
CREATE KEYSPACE IF NOT EXISTS oauth
    WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

USE oauth;

CREATE TABLE IF NOT EXISTS access_tokens (
    access_token text PRIMARY KEY,
    user_id      bigint,
    client_id    bigint,
    expires      bigint
);
```

## Configuration

Create a `.env` file in the project root using the following template:

```env
# Application
GIN_PORT=:8081
CTX_TIMEOUT=2s

# Logger
LEVEL=info
OUTPUT_PATHS=stdout

# Users API
USERS_API_BASE_URL=http://localhost:8080/users/login
RESTY_REQUEST_TIME=150        # Optional: HTTP client timeout in ms. Defaults to 150.

# Cassandra
DB_HOST=127.0.0.1
KEYSPACE=oauth
CONSISTENCY=Quorum            # Optional: Any | One | All | Quorum. Defaults to Quorum.
```

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

3. **Ensure Cassandra is running** and the schema from above has been applied.

4. **Ensure the Users API is running** and accessible at the URL specified in `USERS_API_BASE_URL`.

5. **Run the application:**
   ```bash
   go run src/main.go
   ```

## API Documentation (Swagger)

Interactive API documentation is available once the service is running:

```
http://localhost:8081/swagger/index.html
```

To regenerate the Swagger docs after modifying controller annotations:

```bash
swag init --parseDependency --parseInternal --generalInfo src/main.go --dir ./src
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

| Method  | Path                                      | Description                        | Auth Required |
|---------|-------------------------------------------|------------------------------------|---------------|
| `POST`  | `/oauth/access_token`                     | Create a new access token (login)  | No            |
| `GET`   | `/oauth/access_token/{access_token_id}`   | Get an access token by its ID      | No            |
| `PATCH` | `/oauth/access_token/{access_token_id}`   | Update a token's expiration time   | No            |

### `POST /oauth/access_token`

Validates user credentials against the Users API and issues a new access token stored in Cassandra.

**Request body (`password` grant type):**
```json
{
  "grant_type": "password",
  "username": "user@example.com",
  "password": "secret"
}
```

**Success response `201`:**
```json
{
  "access_token": "5f4dcc3b5aa765d61d8327deb882cf99",
  "user_id": 1,
  "client_id": 0,
  "expires": 1718000000
}
```

### `GET /oauth/access_token/{access_token_id}`

Returns the access token record for the given ID. Used internally by other microservices (e.g. `bookstore-oauth-go`) to validate tokens.

### `PATCH /oauth/access_token/{access_token_id}`

Updates the expiration time of an existing token. The request body must include the token string and the new `expires` Unix timestamp.

**Request body:**
```json
{
  "access_token": "5f4dcc3b5aa765d61d8327deb882cf99",
  "expires": 1720000000
}
```