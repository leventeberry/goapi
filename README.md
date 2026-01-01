# GoAPI

A RESTful API built with Go (Golang) using the Gin web framework. This API provides user management functionality with JWT-based authentication, role-based access control (RBAC), and comprehensive middleware for security and logging.

## Features

- 🔐 **JWT Authentication** - Secure token-based authentication with 60-day expiration
- 👥 **User Management** - Full CRUD operations for user accounts
- 🛡️ **Role-Based Access Control** - Support for `user` and `admin` roles
- 🔒 **Password Security** - Bcrypt password hashing with secure defaults
- ⚡ **Rate Limiting** - IP-based rate limiting (60 requests/minute with burst of 10)
- 📝 **Request Logging** - Comprehensive HTTP request logging with status codes
- 🗄️ **Database Migrations** - Automatic database schema migration using GORM
- 🏥 **Health Check** - Root endpoint for API status verification

## Tech Stack

- **Go 1.24.1** - Programming language
- **Gin** - HTTP web framework
- **GORM** - ORM library for database operations
- **MySQL** - Database (via GORM MySQL driver)
- **JWT (golang-jwt/jwt/v5)** - JSON Web Token implementation
- **Bcrypt (golang.org/x/crypto)** - Password hashing
- **godotenv** - Environment variable management

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
├── initializers/        # Application initialization
│   └── initializers.go     # Database connection and migration
├── schema.sql          # Database schema reference
├── main.go             # Application entry point
├── go.mod              # Go module dependencies
└── README.md           # This file
```

## Prerequisites

- Go 1.24.1 or higher
- MySQL database server
- Git (for cloning the repository)

## Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd goapi
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   
   Create a `.env` file in the root directory:
   ```env
   # Database Configuration
   DB_USER=your_db_user
   DB_PASS=your_db_password
   DB_HOST=localhost:3306
   DB_NAME=your_database_name

   # JWT Secret (use a strong random string)
   JWT_SECRET=your_super_secret_jwt_key_here

   # Server Port (optional, defaults to 8080)
   PORT=8080
   ```

4. **Set up the database**
   
   The application will automatically create the necessary tables using GORM AutoMigrate. Ensure your MySQL database exists and is accessible with the credentials provided in `.env`.

5. **Run the application**
   ```bash
   go run main.go
   ```

   The server will start on `http://localhost:8080` (or the port specified in `PORT` environment variable).

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

## Database

The application uses MySQL with GORM for database operations. The schema is automatically created via GORM AutoMigrate on startup.

### Manual Schema Setup

If you prefer to set up the database manually, refer to `schema.sql` for the table structure.

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

```bash
# Set GIN_MODE to development for verbose logging
export GIN_MODE=debug
go run main.go
```

### Building for Production

```bash
go build -o goapi main.go
./goapi
```

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
