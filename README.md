# GoAPI

A RESTful API built with Go (Golang) using the Gin web framework. This API provides user management functionality with JWT-based authentication, role-based access control (RBAC), and comprehensive middleware for security and logging.

## Features

- 🔐 **JWT Authentication** - Secure token-based authentication with 60-day expiration
- 👥 **User Management** - Full CRUD operations for user accounts
- 🛡️ **Role-Based Access Control** - Support for `user` and `admin` roles
- 🔒 **Password Security** - Bcrypt password hashing with secure defaults
- ⚡ **Rate Limiting** - IP-based rate limiting (60 requests/minute with burst of 10), supports Redis for distributed rate limiting
- 🚀 **Redis Caching** - Optional Redis integration for user caching and distributed rate limiting
- 📝 **Request Logging** - Comprehensive HTTP request logging with status codes
- 🗄️ **Database Migrations** - Automatic database schema migration using GORM
- 🏥 **Health Check** - Root endpoint for API status verification
- 📚 **Swagger/OpenAPI Documentation** - Interactive API documentation with Swagger UI

## Tech Stack

- **Go 1.25.5** - Programming language
- **Gin** - HTTP web framework
- **GORM** - ORM library for database operations
- **PostgreSQL** - Database (via GORM PostgreSQL driver)
- **Redis** - Caching and distributed rate limiting (optional)
- **JWT (golang-jwt/jwt/v5)** - JSON Web Token implementation
- **Bcrypt (golang.org/x/crypto)** - Password hashing
- **godotenv** - Environment variable management
- **Swagger/OpenAPI (swaggo)** - API documentation and interactive UI

## Project Structure

```
goapi/
├── controllers/          # Request handlers
│   ├── authController.go    # Authentication endpoints (login, signup)
│   └── userController.go    # User CRUD operations
├── middleware/          # HTTP middleware
│   ├── auth.go             # JWT authentication middleware
│   ├── bcrypt.go           # Password hashing utilities
│   ├── logger.go            # Request logging middleware
│   └── ratelimit.go        # Rate limiting middleware
├── models/              # Data models
│   └── index.go            # User model definition
├── routes/              # Route definitions
│   ├── index.go            # Main route setup
│   └── userRoutes.go       # User-specific routes
├── cache/               # Cache abstraction layer
│   ├── interfaces.go        # Cache interface definition
│   ├── redis_cache.go      # Redis cache implementation
│   ├── noop_cache.go       # No-op cache (when Redis disabled)
│   ├── constants.go         # Cache key patterns and TTL values
│   └── errors.go           # Cache-specific errors
├── initializers/        # Application initialization
│   └── initializers.go     # Database and Redis connection, migration
├── docs/                # Swagger/OpenAPI documentation (generated)
│   ├── docs.go             # Generated Swagger docs
│   ├── swagger.json        # OpenAPI JSON specification
│   └── swagger.yaml        # OpenAPI YAML specification
├── schema.sql          # Database schema reference
├── Dockerfile          # Docker image definition
├── docker-compose.yml  # Docker Compose configuration
├── Makefile           # Build automation and common tasks
├── main.go            # Application entry point
├── go.mod             # Go module dependencies
├── go.sum             # Go module checksums
└── README.md          # This file
```

## Prerequisites

### For Local Development:
- Go 1.25.5 or higher
- PostgreSQL database server (14+)
- Redis server (optional, for caching and distributed rate limiting)
- Git (for cloning the repository)

### For Docker:
- Docker and Docker Compose installed
- Git (for cloning the repository)

## Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd goapi
   ```

2. **Install dependencies**
   
   Using Make (recommended):
   ```bash
   make install
   ```
   
   Or manually:
   ```bash
   go mod download
   go mod tidy
   ```

3. **Set up environment variables**
   
   Create a `.env` file in the root directory (or copy from `.env.example`):
   ```env
   # Database Configuration
   DB_USER=your_db_user
   DB_PASS=your_db_password
   DB_HOST=localhost
   DB_PORT=5432
   DB_NAME=your_database_name

   # JWT Secret (use a strong random string)
   JWT_SECRET=your_super_secret_jwt_key_here

   # Server Port (optional, defaults to 8080)
   PORT=8080

   # Redis Configuration (optional)
   # Set REDIS_ENABLED=true to enable Redis caching
   # If Redis is disabled, the application will use a no-op cache
   REDIS_ENABLED=false
   REDIS_HOST=localhost
   REDIS_PORT=6379
   REDIS_PASSWORD=
   ```

4. **Set up the database**
   
   The application will automatically create the necessary tables using GORM AutoMigrate. Ensure your PostgreSQL database exists and is accessible with the credentials provided in `.env`.

5. **Set up Redis (optional)**
   
   If you want to use Redis for caching and distributed rate limiting:
   - Install Redis locally or use Docker: `docker run -d -p 6379:6379 redis:7-alpine`
   - Set `REDIS_ENABLED=true` in your `.env` file
   - Configure `REDIS_HOST` and `REDIS_PORT` if different from defaults
   - If Redis is not available, the application will gracefully degrade to in-memory caching

6. **Run the application**
   
   Using Make (recommended):
   ```bash
   make run
   ```
   
   Or manually:
   ```bash
   go run main.go
   ```

   The server will start on `http://localhost:8080` (or the port specified in `PORT` environment variable).

7. **Generate Swagger documentation** (if you modify API endpoints)
   
   Using Make (recommended):
   ```bash
   make swagger
   ```
   
   Or manually:
   ```bash
   # Install swag CLI tool
   go install github.com/swaggo/swag/cmd/swag@latest
   
   # Generate Swagger docs from annotations
   swag init
   ```

## Makefile Commands

This project includes a Makefile with convenient commands for common tasks. Run `make help` to see all available commands.

### Quick Start Commands

```bash
# Show all available commands
make help

# Full setup (install deps + generate Swagger docs)
make setup

# Run locally
make run

# Start with Docker
make docker-up

# View Docker logs
make docker-logs-api
```

### Common Commands

**Local Development:**
- `make install` or `make deps` - Install Go dependencies
- `make run` - Run the application locally
- `make build` - Build the application binary
- `make test` - Run tests
- `make test-coverage` - Run tests with coverage report
- `make clean` - Clean build artifacts (binary, coverage files)

**Docker:**
- `make docker-build` - Build Docker images
- `make docker-up` - Start Docker containers in detached mode
- `make docker-down` - Stop Docker containers
- `make docker-down-volumes` - Stop containers and remove volumes (clears database)
- `make docker-logs` - View all container logs (follow mode)
- `make docker-logs-api` - View API container logs only
- `make docker-logs-db` - View database container logs only
- `make docker-logs-redis` - View Redis container logs only
- `make docker-restart` - Restart Docker containers
- `make docker-rebuild` - Rebuild and restart containers
- `make docker-ps` - Show running Docker containers
- `make docker-shell-api` - Open shell in API container
- `make docker-shell-db` - Open PostgreSQL shell in database container
- `make docker-shell-redis` - Open Redis CLI in Redis container
- `make docker-open-redis-commander` - Open Redis Commander web UI in browser
- `make docker-open-pgadmin` - Open pgAdmin web UI in browser

**Documentation:**
- `make swagger` - Generate Swagger documentation (auto-installs swag if needed)
- `make swag` - Install swag CLI tool

**Database:**
- `make db-migrate` - Run database migrations (local)
- `make db-seed` - Seed database with sample data (placeholder)

**All-in-one:**
- `make dev` - Install deps and run locally
- `make dev-docker` - Start Docker and follow API logs
- `make setup` - Full setup: install deps and generate Swagger docs
- `make all` - Clean, install, generate docs, and build
- `make prod-build` - Production build: clean and build
- `make docker-all` - Full Docker rebuild: down, build, up

## Docker Setup

### Quick Start with Docker Compose

The easiest way to run the entire application stack (API + PostgreSQL) is using Docker Compose:

#### Using Make (Recommended)

1. **Clone the repository** (if you haven't already)
   ```bash
   git clone <repository-url>
   cd goapi
   ```

2. **Set JWT Secret** (optional, but recommended)
   
   Create a `.env` file or set the `JWT_SECRET` environment variable:
   ```bash
   export JWT_SECRET=your_super_secret_jwt_key_here
   ```
   
   Or create a `.env` file:
   ```env
   JWT_SECRET=your_super_secret_jwt_key_here
   ```

3. **Build and start services**
   ```bash
   make docker-up
   ```
   
   Or for a complete rebuild:
   ```bash
   make docker-all
   ```

4. **View logs**
   ```bash
   make docker-logs-api
   ```

5. **Access the API and Admin Tools**
   - API: `http://localhost:8080`
   - Swagger UI: `http://localhost:8080/swagger/index.html`
   - PostgreSQL: `localhost:5432`
   - **Redis Commander**: `http://localhost:8081` (Username: `admin`, Password: `admin`)
   - **pgAdmin**: `http://localhost:5050` (Email: `admin@goapi.com`, Password: `admin`)

6. **Stop services**
   ```bash
   make docker-down
   ```

7. **Stop and remove volumes** (clears database data)
   ```bash
   make docker-down-volumes
   ```

#### Using Docker Compose Directly

If you prefer using Docker Compose directly:

```bash
# Build and start
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

### Docker Compose Services

- **`api`**: Go API application (port 8080)
- **`db`**: PostgreSQL 16 database (port 5432)
- **`redis`**: Redis 7 cache server (port 6379)
- **`redis-commander`**: Redis Commander web UI (port 8081) - Admin interface for Redis
- **`pgadmin`**: pgAdmin 4 web UI (port 5050) - Admin interface for PostgreSQL

### Default Database Credentials (Docker)

When using Docker Compose, the database is automatically configured with:
- **Database**: `goapi`
- **User**: `goapi_user`
- **Password**: `goapi_password`

These credentials are set in `docker-compose.yml` and can be customized if needed.

### Docker Image Details

The Dockerfile uses a multi-stage build:
- **Builder stage**: Uses `golang:1.25-alpine` to compile the application
- **Final stage**: Uses `alpine:latest` for a minimal production image (~10MB)

The image includes:
- Compiled Go binary
- Swagger documentation (if generated)
- CA certificates for HTTPS requests

### Building Docker Image Manually

If you want to build just the API Docker image:

```bash
# Build the image
docker build -t goapi:latest .

# Run the container
docker run -p 8080:8080 \
  -e DB_USER=goapi_user \
  -e DB_PASS=goapi_password \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_NAME=goapi \
  -e JWT_SECRET=your_secret_key \
  goapi:latest
```

Or using Make:
```bash
make docker-build
```

## API Documentation

### Swagger UI

The API includes interactive Swagger/OpenAPI documentation accessible at:

**http://localhost:8080/swagger/index.html**

The Swagger UI provides:
- Interactive API testing interface
- Complete endpoint documentation
- Request/response examples
- Authentication testing with JWT tokens
- Schema definitions for all models

### Generating Swagger Documentation

After modifying API endpoints or adding new ones, regenerate the Swagger documentation:

```bash
swag init
```

This command scans your code for Swagger annotations (comments starting with `@Summary`, `@Description`, `@Tags`, etc.) and generates the documentation files in the `docs/` directory.

## API Endpoints

### Public Endpoints

#### Health Check
- **GET** `/`
  - Returns API status
  - **Response:**
    ```json
    {
      "message": "Welcome!",
      "status": 200
    }
    ```

#### Authentication

- **POST** `/register`
  - Register a new user
  - **Request Body:**
    ```json
    {
      "first_name": "John",
      "last_name": "Doe",
      "email": "john.doe@example.com",
      "password": "password123",
      "phone_number": "+1234567890",
      "role": "user"
    }
    ```
  - **Response (201):**
    ```json
    {
      "token": {
        "api_key": "uuid-string",
        "jwt_token": "jwt-token-string"
      },
      "user": {
        "id": 1,
        "email": "john.doe@example.com"
      }
    }
    ```
  - **Valid Roles:** `user`, `admin`

- **POST** `/login`
  - Authenticate and receive JWT token
  - **Request Body:**
    ```json
    {
      "email": "john.doe@example.com",
      "password": "password123"
    }
    ```
  - **Response (200):**
    ```json
    {
      "token": {
        "api_key": "uuid-string",
        "jwt_token": "jwt-token-string"
      },
      "user": {
        "id": 1,
        "email": "john.doe@example.com"
      }
    }
    ```

### Protected Endpoints (Require Authentication)

All user endpoints require a valid JWT token in the `Authorization` header:
```
Authorization: Bearer <jwt_token>
```

#### User Management

- **GET** `/users`
  - Get all users
  - **Headers:** `Authorization: Bearer <token>`
  - **Response (200):** Array of user objects

- **GET** `/users/:id`
  - Get a specific user by ID
  - **Headers:** `Authorization: Bearer <token>`
  - **Response (200):** User object
  - **Response (404):** `{"error": "User not found"}`

- **POST** `/users`
  - Create a new user (authenticated users only)
  - **Headers:** `Authorization: Bearer <token>`
  - **Request Body:**
    ```json
    {
      "first_name": "Jane",
      "last_name": "Smith",
      "email": "jane.smith@example.com",
      "password": "password123",
      "phone_number": "+1234567891",
      "role": "user"
    }
    ```
  - **Note:** `phone_number` and `role` are optional. Default role is `user`.
  - **Response (201):** Created user object

- **PUT** `/users/:id`
  - Update a user (partial updates supported)
  - **Headers:** `Authorization: Bearer <token>`
  - **Request Body:** (all fields optional, but at least one required)
    ```json
    {
      "first_name": "Jane",
      "last_name": "Smith",
      "email": "jane.smith@example.com",
      "password": "newpassword123",
      "phone_number": "+1234567891",
      "role": "admin"
    }
    ```
  - **Response (200):** Updated user object

- **DELETE** `/users/:id`
  - Delete a user (Admin only)
  - **Headers:** `Authorization: Bearer <token>`
  - **Response (200):** `{"message": "User deleted successfully"}`
  - **Response (403):** `{"error": "Insufficient permissions"}` (if not admin)

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. Tokens are valid for 60 days and include:
- User ID (subject)
- API Key (UUID)
- Issued and expiration timestamps

### Using Authentication

Include the JWT token in the `Authorization` header for protected endpoints:
```bash
curl -H "Authorization: Bearer <your_jwt_token>" http://localhost:8080/users
```

## User Model

```go
type User struct {
    ID        int       `json:"user_id"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Email     string    `json:"email"`        // Unique
    PassHash  string    `json:"-"`            // Never returned in JSON
    PhoneNum  string    `json:"phone_number"`
    Role      string    `json:"role"`         // "user" or "admin"
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

## Middleware

### Rate Limiting
- **Default:** 60 requests per minute per IP
- **Burst:** 10 requests
- **Response (429):** `{"error": "Rate limit exceeded. Please try again later."}`
- **Redis Support:** When Redis is enabled, rate limiting is distributed across all API instances
- **Fallback:** If Redis is unavailable, automatically falls back to in-memory rate limiting

### Request Logging
Logs all HTTP requests with:
- HTTP method
- Request path
- Status code
- Response time
- Client IP
- User agent

Log levels:
- **INFO:** Status codes < 400
- **WARN:** Status codes 400-499
- **ERROR:** Status codes ≥ 500

### Authentication Middleware
- Validates JWT tokens from `Authorization` header
- Extracts user ID and API key from token claims
- Stores user information in request context

### Role-Based Access Control
- `RequireRole("admin")` middleware restricts endpoints to admin users
- Currently used for user deletion endpoint

## Caching

The application supports optional Redis caching for improved performance:

### Cache Features
- **User Caching**: Caches user lookups by ID and email (15-minute TTL)
- **Cache-Aside Pattern**: Checks cache first, falls back to database on miss
- **Automatic Invalidation**: Cache is invalidated on user updates and deletes
- **Distributed Rate Limiting**: Redis enables shared rate limits across multiple API instances
- **Graceful Degradation**: If Redis is unavailable, uses no-op cache (app continues to work)

### Cache Configuration
- **User Cache TTL**: 15 minutes (configurable in `cache/constants.go`)
- **Rate Limit Window**: 1 minute (configurable in `cache/constants.go`)
- **Key Patterns**:
  - User by ID: `user:id:{id}`
  - User by Email: `user:email:{email}`
  - Rate Limit: `ratelimit:{ip}`

### Cache Invalidation Strategy
- **On User Update**: All cached entries for the user are invalidated
- **On User Delete**: All cached entries for the user are removed
- **On User Create**: New user is stored in cache
- **Email Changes**: Old email cache key is deleted when email is updated

### Enabling Redis
Set `REDIS_ENABLED=true` in your `.env` file. The application will automatically:
- Connect to Redis on startup
- Use Redis for caching and rate limiting
- Fall back to in-memory if Redis connection fails

## Database

The application uses PostgreSQL with GORM for database operations. The schema is automatically created via GORM AutoMigrate on startup.

### Manual Schema Setup

If you prefer to set up the database manually, refer to `schema.sql` for the table structure. The schema includes:
- Automatic `updated_at` timestamp updates via PostgreSQL trigger
- Index on email column for faster lookups
- Proper PostgreSQL data types (SERIAL for auto-increment IDs)

## Error Responses

The API returns standard HTTP status codes:

- **200 OK** - Successful request
- **201 Created** - Resource created successfully
- **400 Bad Request** - Invalid request data
- **401 Unauthorized** - Missing or invalid authentication
- **403 Forbidden** - Insufficient permissions
- **404 Not Found** - Resource not found
- **409 Conflict** - Resource already exists (e.g., duplicate email)
- **429 Too Many Requests** - Rate limit exceeded
- **500 Internal Server Error** - Server error

Error responses follow this format:
```json
{
  "error": "Error message description"
}
```

## Security Features

- **Password Hashing:** Bcrypt with default cost (10 rounds)
- **JWT Tokens:** HMAC-SHA256 signed tokens
- **Rate Limiting:** Prevents abuse and DoS attacks
- **Input Validation:** Request body validation using Gin's binding
- **SQL Injection Protection:** GORM parameterized queries
- **Password Requirements:** Minimum 8 characters for passwords

## Development

### Running in Development Mode

Using Make:
```bash
# Set GIN_MODE to development for verbose logging
export GIN_MODE=debug
make run
```

Or manually:
```bash
export GIN_MODE=debug
go run main.go
```

### Building for Production

Using Make:
```bash
make prod-build
# Binary will be created as 'goapi'
./goapi
```

Or manually:
```bash
go build -o goapi main.go
./goapi
```

### Testing

Run all tests:
```bash
make test
```

Run tests with coverage:
```bash
make test-coverage
# Opens coverage.html in your default browser
```

### Generating Swagger Documentation

When you add or modify API endpoints, update the Swagger annotations in your controller functions and regenerate the docs:

```bash
# Install swag CLI (if not already installed)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
swag init
```

The Swagger annotations use the following format:
```go
// @Summary      Brief summary of the endpoint
// @Description  Detailed description
// @Tags         tag-name
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Parameter description"
// @Success      200  {object}  models.User
// @Failure      400  {object}  map[string]string
// @Router       /users/{id} [get]
```

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
