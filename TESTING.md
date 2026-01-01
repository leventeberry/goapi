# API Testing Guide

## Prerequisites

Before testing, ensure:
1. Database is running (PostgreSQL)
2. `.env` file is configured with database credentials
3. Server is running: `make run` or `go run main.go`

## Quick Test Commands

### Using cURL (Linux/Mac/Git Bash)

```bash
# 1. Health Check
curl http://localhost:8080/

# 2. Register User
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@test.com",
    "password": "password123",
    "phone_number": "+1234567890",
    "role": "user"
  }'

# 3. Register Admin
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Admin",
    "last_name": "User",
    "email": "admin@test.com",
    "password": "adminpass123",
    "phone_number": "+1234567891",
    "role": "admin"
  }'

# 4. Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@test.com",
    "password": "password123"
  }'

# Save the JWT token from response, then:

# 5. Get All Users (replace TOKEN with actual token)
curl -X GET http://localhost:8080/users \
  -H "Authorization: Bearer TOKEN"

# 6. Get User by ID
curl -X GET http://localhost:8080/users/1 \
  -H "Authorization: Bearer TOKEN"

# 7. Create User
curl -X POST http://localhost:8080/users \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jane",
    "last_name": "Smith",
    "email": "jane.smith@test.com",
    "password": "password123",
    "phone_number": "+1234567892",
    "role": "user"
  }'

# 8. Update User
curl -X PUT http://localhost:8080/users/1 \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John Updated",
    "last_name": "Doe Updated"
  }'

# 9. Delete User (Admin only - replace ADMIN_TOKEN)
curl -X DELETE http://localhost:8080/users/2 \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

### Using PowerShell (Windows)

```powershell
# 1. Health Check
Invoke-RestMethod -Uri "http://localhost:8080/" -Method GET

# 2. Register User
$registerData = @{
    first_name = "John"
    last_name = "Doe"
    email = "john.doe@test.com"
    password = "password123"
    phone_number = "+1234567890"
    role = "user"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/register" -Method POST -Body $registerData -ContentType "application/json"

# 3. Login
$loginData = @{
    email = "john.doe@test.com"
    password = "password123"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:8080/login" -Method POST -Body $loginData -ContentType "application/json"
$token = $response.token.jwt_token

# 4. Get All Users
$headers = @{
    "Authorization" = "Bearer $token"
}
Invoke-RestMethod -Uri "http://localhost:8080/users" -Method GET -Headers $headers
```

## Test Coverage

### ✅ Endpoints to Test

#### Public Endpoints
- [x] `GET /` - Health check
- [x] `POST /register` - User registration
- [x] `POST /login` - User authentication

#### Protected Endpoints (Require JWT)
- [x] `GET /users` - Get all users
- [x] `GET /users/:id` - Get user by ID
- [x] `POST /users` - Create user
- [x] `PUT /users/:id` - Update user
- [x] `DELETE /users/:id` - Delete user (Admin only)

### ✅ Test Scenarios

#### Authentication Tests
1. ✅ Register new user (user role)
2. ✅ Register new admin (admin role)
3. ✅ Login with valid credentials
4. ✅ Login with invalid credentials (should return 401)
5. ✅ Register with duplicate email (should return 409)
6. ✅ Register with invalid role (should return 400)
7. ✅ Register with invalid email format (should return 400)
8. ✅ Register with short password (should return 400)

#### Authorization Tests
1. ✅ Access protected endpoint without token (should return 401)
2. ✅ Access protected endpoint with invalid token (should return 401)
3. ✅ Access protected endpoint with valid token (should succeed)
4. ✅ Delete user as regular user (should return 403)
5. ✅ Delete user as admin (should succeed)

#### User Management Tests
1. ✅ Get all users (authenticated)
2. ✅ Get user by ID (authenticated)
3. ✅ Get non-existent user (should return 404)
4. ✅ Create user (authenticated)
5. ✅ Update user (partial update)
6. ✅ Update user with empty body (should return 400)
7. ✅ Delete user (admin only)

#### Error Handling Tests
1. ✅ Invalid request body
2. ✅ Missing required fields
3. ✅ Invalid data types
4. ✅ Resource not found (404)
5. ✅ Unauthorized access (401)
6. ✅ Forbidden access (403)
7. ✅ Conflict (409 - duplicate email)

### ✅ Middleware Tests

#### Rate Limiting
- Test: Make 70 requests in quick succession
- Expected: First 60 succeed, then 429 Too Many Requests

#### Request Logging
- Check server logs for request details
- Verify: Method, path, status code, response time, IP, user agent

#### Authentication Middleware
- Test: Access protected route without token
- Expected: 401 Unauthorized

#### RBAC Middleware
- Test: Regular user tries to delete user
- Expected: 403 Forbidden

## Automated Testing

### Option 1: Use the Test Scripts

**Bash (Linux/Mac/Git Bash):**
```bash
chmod +x test_api.sh
./test_api.sh
```

**PowerShell (Windows):**
```powershell
powershell -ExecutionPolicy Bypass -File test_api.ps1
```

### Option 2: Use Go Tests

```bash
# Run all tests
go test -v ./...

# Run specific test
go test -v -run TestHealthCheck
```

### Option 3: Use Swagger UI

1. Start the server: `make run`
2. Open browser: `http://localhost:8080/swagger/index.html`
3. Use the interactive Swagger UI to test endpoints
4. Click "Authorize" and enter your JWT token
5. Test each endpoint interactively

## Expected Test Results

When all tests pass, you should see:

```
✓ Health check endpoint
✓ Register user endpoint
✓ Register admin user endpoint
✓ Login endpoint
✓ Login with invalid credentials
✓ Get all users endpoint
✓ Get all users without authentication
✓ Get user by ID endpoint
✓ Get non-existent user
✓ Create user endpoint
✓ Update user endpoint
✓ Update user with invalid data
✓ Delete user as admin
✓ Delete user as regular user
✓ Register duplicate email
✓ Register with invalid role
✓ Register with invalid email format
✓ Register with short password

Test Summary
Passed: 18
Failed: 0
Total: 18
```

## Troubleshooting

### Server won't start
- Check database connection in `.env`
- Ensure PostgreSQL is running
- Check port 8080 is not in use

### Tests fail with 401
- Verify JWT token is valid
- Check token hasn't expired
- Ensure token is in Authorization header: `Bearer <token>`

### Tests fail with 500
- Check database connection
- Review server logs
- Verify all environment variables are set

### Rate limiting issues
- Wait 1 minute between test runs
- Or restart server to reset rate limiter

## Performance Testing

### Load Testing with Apache Bench

```bash
# Test health endpoint with 1000 requests, 10 concurrent
ab -n 1000 -c 10 http://localhost:8080/

# Test authenticated endpoint (requires token)
ab -n 100 -c 5 -H "Authorization: Bearer TOKEN" http://localhost:8080/users
```

### Load Testing with wrk

```bash
# Install: brew install wrk (Mac) or download from GitHub
wrk -t4 -c100 -d30s http://localhost:8080/
```

## Integration Testing

For full integration testing with a test database:

1. Create a separate test database
2. Update `.env.test` with test database credentials
3. Run migrations on test database
4. Execute test suite
5. Clean up test data

## Continuous Integration

Example CI/CD test command:

```yaml
# .github/workflows/test.yml
- name: Run tests
  run: |
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
```

