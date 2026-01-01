# API Testing Guide

## Prerequisites

Before testing, ensure:
1. Database is running (PostgreSQL)
2. `.env` file is configured with database credentials
3. Server is running: `make run` or `go run main.go`
4. Redis (optional) - For testing Redis caching features

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
- **With Redis**: Rate limits are shared across all API instances
- **Without Redis**: Rate limits are per-instance (in-memory)

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

## Redis Caching Tests

### Test 1: Application with Redis Disabled

**Setup:**
```bash
# Set REDIS_ENABLED=false in .env or don't set it
REDIS_ENABLED=false
```

**Expected Behavior:**
- ✅ Application starts successfully
- ✅ No Redis connection errors
- ✅ User operations work normally
- ✅ Rate limiting uses in-memory storage
- ✅ No cache entries in Redis (if Redis is running but disabled)

**Validation:**
```bash
# Check logs for "Redis is disabled"
# All user operations should work normally
# Rate limiting should work per-instance
```

### Test 2: Application with Redis Enabled

**Setup:**
```bash
# Start Redis (if using Docker)
docker run -d -p 6379:6379 redis:7-alpine

# Or use docker-compose
make docker-up

# Set in .env
REDIS_ENABLED=true
REDIS_HOST=localhost
REDIS_PORT=6379
```

**Expected Behavior:**
- ✅ Application connects to Redis on startup
- ✅ Log shows "Redis connection established"
- ✅ User caching works (faster subsequent lookups)
- ✅ Rate limiting uses Redis (distributed)

**Validation:**
```bash
# Check Redis for cached keys
docker exec -it goapi_redis redis-cli
> KEYS *
> GET user:id:1
> GET ratelimit:127.0.0.1
```

### Test 3: User Caching - Cache Hits

**Steps:**
1. Get user by ID: `GET /users/1` (first request - cache miss)
2. Get same user again: `GET /users/1` (should be faster - cache hit)
3. Get user by email: `GET /users?email=test@example.com` (cache miss)
4. Get same user by email again (cache hit)

**Expected:**
- ✅ First request queries database
- ✅ Second request uses cache (faster response)
- ✅ Both ID and email keys are cached
- ✅ Cache TTL is 15 minutes

**Validation:**
```bash
# Check Redis keys
redis-cli KEYS "user:*"
# Should see: user:id:1 and user:email:test@example.com
```

### Test 4: Cache Invalidation on Update

**Steps:**
1. Get user: `GET /users/1` (caches user)
2. Update user: `PUT /users/1` with new data
3. Get user again: `GET /users/1` (should have updated data)

**Expected:**
- ✅ Cache is invalidated on update
- ✅ Updated user is stored in cache
- ✅ Old cache entries are deleted
- ✅ Response contains updated data

**Validation:**
```bash
# Before update
redis-cli GET "user:id:1"

# After update
redis-cli GET "user:id:1"
# Should show updated user data
```

### Test 5: Cache Invalidation on Delete

**Steps:**
1. Get user: `GET /users/1` (caches user)
2. Delete user: `DELETE /users/1` (admin only)
3. Try to get deleted user: `GET /users/1` (should return 404)

**Expected:**
- ✅ Cache entries are deleted when user is deleted
- ✅ Both ID and email keys are removed
- ✅ Subsequent requests return 404

**Validation:**
```bash
# After delete
redis-cli KEYS "user:*1*"
# Should return empty (no keys)
```

### Test 6: Cache Invalidation on Email Change

**Steps:**
1. Get user: `GET /users/1` (caches with old email)
2. Update email: `PUT /users/1` with new email
3. Get user by old email (should not find in cache)
4. Get user by new email (should find in cache)

**Expected:**
- ✅ Old email cache key is deleted
- ✅ New email cache key is created
- ✅ ID cache key is updated
- ✅ No stale data in cache

### Test 7: Distributed Rate Limiting

**Setup:**
- Start two API instances (different ports)
- Both connected to same Redis instance

**Steps:**
1. Make 30 requests to instance 1
2. Make 30 requests to instance 2
3. Make 1 more request to either instance

**Expected:**
- ✅ Total requests across both instances = 60
- ✅ 61st request should return 429 (rate limit exceeded)
- ✅ Rate limit is shared across instances

**Validation:**
```bash
# Check rate limit in Redis
redis-cli GET "ratelimit:127.0.0.1"
# Should show count across all instances
```

### Test 8: Graceful Degradation - Redis Unavailable

**Setup:**
1. Start application with Redis enabled
2. Stop Redis while application is running

**Expected:**
- ✅ Application continues to work
- ✅ Cache operations fail gracefully (no-op)
- ✅ Rate limiting falls back to in-memory
- ✅ No application crashes

**Validation:**
```bash
# Stop Redis
docker stop goapi_redis

# Make API requests - should still work
curl http://localhost:8080/users/1

# Check logs - should show cache errors but continue
```

### Test 9: Cache TTL Expiration

**Steps:**
1. Get user: `GET /users/1` (caches user)
2. Wait 16 minutes (or manually expire in Redis)
3. Get user again: `GET /users/1`

**Expected:**
- ✅ Cache expires after 15 minutes
- ✅ Request after expiration queries database
- ✅ User is re-cached after expiration

**Validation:**
```bash
# Check TTL
redis-cli TTL "user:id:1"
# Should show seconds remaining (max 900 = 15 minutes)

# Manually expire for testing
redis-cli EXPIRE "user:id:1" 0
```

### Test 10: Cache Performance Comparison

**Setup:**
- Test with Redis enabled vs disabled

**Steps:**
1. Make 100 requests to `GET /users/1` with Redis
2. Make 100 requests to `GET /users/1` without Redis
3. Compare response times

**Expected:**
- ✅ Cached requests are significantly faster
- ✅ First request (cache miss) is slower
- ✅ Subsequent requests (cache hits) are faster

**Validation:**
```bash
# With Redis (after first request)
time curl http://localhost:8080/users/1
# Should be < 10ms

# Without Redis
time curl http://localhost:8080/users/1
# Should be > 50ms (database query)
```

## Redis Testing Checklist

### Basic Functionality
- [ ] Application starts with Redis disabled (no-op cache)
- [ ] Application starts with Redis enabled
- [ ] Application starts with Redis unavailable (graceful degradation)
- [ ] User caching works (cache hits/misses)
- [ ] Cache invalidation on user update
- [ ] Cache invalidation on user delete
- [ ] Cache invalidation on email change
- [ ] Cache TTL expiration works

### Rate Limiting
- [ ] In-memory rate limiting works (Redis disabled)
- [ ] Redis rate limiting works (Redis enabled)
- [ ] Distributed rate limiting across instances
- [ ] Rate limit fallback when Redis unavailable

### Performance
- [ ] Cached requests are faster than database queries
- [ ] Cache reduces database load
- [ ] No performance degradation when Redis unavailable

### Error Handling
- [ ] Application continues when Redis connection fails
- [ ] Cache errors don't crash application
- [ ] Rate limiting works even if Redis fails

## Continuous Integration

Example CI/CD test command:

```yaml
# .github/workflows/test.yml
- name: Run tests
  run: |
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

- name: Test with Redis
  run: |
    docker run -d -p 6379:6379 redis:7-alpine
    export REDIS_ENABLED=true
    go test -v ./...
```

