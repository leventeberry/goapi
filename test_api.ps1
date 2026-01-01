# PowerShell API Test Script for Windows
# Tests all endpoints of the GoAPI

$BASE_URL = "http://localhost:8080"
$PASSED = 0
$FAILED = 0

# Helper function to make API calls
function Invoke-ApiCall {
    param(
        [string]$Method,
        [string]$Endpoint,
        [string]$Data = $null,
        [string]$Token = $null
    )
    
    $headers = @{
        "Content-Type" = "application/json"
    }
    
    if ($Token) {
        $headers["Authorization"] = "Bearer $Token"
    }
    
    try {
        if ($Data) {
            $response = Invoke-RestMethod -Uri "$BASE_URL$Endpoint" -Method $Method -Headers $headers -Body $Data -ErrorAction Stop
            $statusCode = 200
        } else {
            $response = Invoke-RestMethod -Uri "$BASE_URL$Endpoint" -Method $Method -Headers $headers -ErrorAction Stop
            $statusCode = 200
        }
        return @{ StatusCode = $statusCode; Body = $response }
    } catch {
        $statusCode = $_.Exception.Response.StatusCode.value__
        $errorBody = $_.ErrorDetails.Message
        return @{ StatusCode = $statusCode; Body = $errorBody }
    }
}

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "  GoAPI Comprehensive Test Suite" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# Check if server is running
Write-Host "Checking if server is running..." -ForegroundColor Yellow
try {
    $check = Invoke-WebRequest -Uri "$BASE_URL/" -UseBasicParsing -ErrorAction Stop
    Write-Host "Server is running" -ForegroundColor Green
} catch {
    Write-Host "Error: Server is not running at $BASE_URL" -ForegroundColor Red
    Write-Host "Please start the server first: make run" -ForegroundColor Yellow
    exit 1
}
Write-Host ""

# Test 1: Health Check
Write-Host "Test 1: Health Check (GET /)" -ForegroundColor Yellow
$result = Invoke-ApiCall -Method "GET" -Endpoint "/"
if ($result.StatusCode -eq 200 -and $result.Body.message -eq "Welcome!") {
    Write-Host "✓ PASS: Health check endpoint" -ForegroundColor Green
    $PASSED++
} else {
    Write-Host "✗ FAIL: Health check endpoint (expected 200, got $($result.StatusCode))" -ForegroundColor Red
    $FAILED++
}
Write-Host ""

# Test 2: Register User
Write-Host "Test 2: Register User (user role)" -ForegroundColor Yellow
$registerData = @{
    first_name = "John"
    last_name = "Doe"
    email = "john.doe@test.com"
    password = "password123"
    phone_number = "+1234567890"
    role = "user"
} | ConvertTo-Json

$result = Invoke-ApiCall -Method "POST" -Endpoint "/register" -Data $registerData
if ($result.StatusCode -eq 200 -and $result.Body.token.jwt_token) {
    $USER_TOKEN = $result.Body.token.jwt_token
    $USER_ID = $result.Body.user.id
    Write-Host "✓ PASS: Register user endpoint" -ForegroundColor Green
    $PASSED++
} else {
    Write-Host "✗ FAIL: Register user endpoint (expected 200, got $($result.StatusCode))" -ForegroundColor Red
    $FAILED++
}
Write-Host ""

# Test 3: Register Admin
Write-Host "Test 3: Register Admin User" -ForegroundColor Yellow
$adminData = @{
    first_name = "Admin"
    last_name = "User"
    email = "admin@test.com"
    password = "adminpass123"
    phone_number = "+1234567891"
    role = "admin"
} | ConvertTo-Json

$result = Invoke-ApiCall -Method "POST" -Endpoint "/register" -Data $adminData
if ($result.StatusCode -eq 200 -and $result.Body.token.jwt_token) {
    $ADMIN_TOKEN = $result.Body.token.jwt_token
    $ADMIN_ID = $result.Body.user.id
    Write-Host "✓ PASS: Register admin user endpoint" -ForegroundColor Green
    $PASSED++
} else {
    Write-Host "✗ FAIL: Register admin user endpoint (expected 200, got $($result.StatusCode))" -ForegroundColor Red
    $FAILED++
}
Write-Host ""

# Test 4: Login
Write-Host "Test 4: Login" -ForegroundColor Yellow
$loginData = @{
    email = "john.doe@test.com"
    password = "password123"
} | ConvertTo-Json

$result = Invoke-ApiCall -Method "POST" -Endpoint "/login" -Data $loginData
if ($result.StatusCode -eq 200 -and $result.Body.token.jwt_token) {
    $LOGIN_TOKEN = $result.Body.token.jwt_token
    Write-Host "✓ PASS: Login endpoint" -ForegroundColor Green
    $PASSED++
} else {
    Write-Host "✗ FAIL: Login endpoint (expected 200, got $($result.StatusCode))" -ForegroundColor Red
    $FAILED++
}
Write-Host ""

# Test 5: Login with Invalid Credentials
Write-Host "Test 5: Login with Invalid Credentials" -ForegroundColor Yellow
$invalidLogin = @{
    email = "john.doe@test.com"
    password = "wrongpassword"
} | ConvertTo-Json

$result = Invoke-ApiCall -Method "POST" -Endpoint "/login" -Data $invalidLogin
if ($result.StatusCode -eq 401) {
    Write-Host "✓ PASS: Login with invalid credentials (should return 401)" -ForegroundColor Green
    $PASSED++
} else {
    Write-Host "✗ FAIL: Login with invalid credentials (expected 401, got $($result.StatusCode))" -ForegroundColor Red
    $FAILED++
}
Write-Host ""

# Test 6: Get All Users (Authenticated)
Write-Host "Test 6: Get All Users (Authenticated)" -ForegroundColor Yellow
$result = Invoke-ApiCall -Method "GET" -Endpoint "/users" -Token $USER_TOKEN
if ($result.StatusCode -eq 200 -and $result.Body.Count -gt 0) {
    Write-Host "✓ PASS: Get all users endpoint" -ForegroundColor Green
    $PASSED++
} else {
    Write-Host "✗ FAIL: Get all users endpoint (expected 200, got $($result.StatusCode))" -ForegroundColor Red
    $FAILED++
}
Write-Host ""

# Test 7: Get All Users (Unauthenticated)
Write-Host "Test 7: Get All Users (Unauthenticated - should fail)" -ForegroundColor Yellow
$result = Invoke-ApiCall -Method "GET" -Endpoint "/users"
if ($result.StatusCode -eq 401) {
    Write-Host "✓ PASS: Get all users without authentication (should return 401)" -ForegroundColor Green
    $PASSED++
} else {
    Write-Host "✗ FAIL: Get all users without authentication (expected 401, got $($result.StatusCode))" -ForegroundColor Red
    $FAILED++
}
Write-Host ""

# Test 8: Get User by ID
Write-Host "Test 8: Get User by ID" -ForegroundColor Yellow
$result = Invoke-ApiCall -Method "GET" -Endpoint "/users/$USER_ID" -Token $USER_TOKEN
if ($result.StatusCode -eq 200 -and $result.Body.email -eq "john.doe@test.com") {
    Write-Host "✓ PASS: Get user by ID endpoint" -ForegroundColor Green
    $PASSED++
} else {
    Write-Host "✗ FAIL: Get user by ID endpoint (expected 200, got $($result.StatusCode))" -ForegroundColor Red
    $FAILED++
}
Write-Host ""

# Test 9: Get Non-Existent User
Write-Host "Test 9: Get Non-Existent User" -ForegroundColor Yellow
$result = Invoke-ApiCall -Method "GET" -Endpoint "/users/99999" -Token $USER_TOKEN
if ($result.StatusCode -eq 404) {
    Write-Host "✓ PASS: Get non-existent user (should return 404)" -ForegroundColor Green
    $PASSED++
} else {
    Write-Host "✗ FAIL: Get non-existent user (expected 404, got $($result.StatusCode))" -ForegroundColor Red
    $FAILED++
}
Write-Host ""

# Test 10: Create User (Authenticated)
Write-Host "Test 10: Create User (Authenticated)" -ForegroundColor Yellow
$createData = @{
    first_name = "Jane"
    last_name = "Smith"
    email = "jane.smith@test.com"
    password = "password123"
    phone_number = "+1234567892"
    role = "user"
} | ConvertTo-Json

$result = Invoke-ApiCall -Method "POST" -Endpoint "/users" -Data $createData -Token $USER_TOKEN
if ($result.StatusCode -eq 201 -and $result.Body.email -eq "jane.smith@test.com") {
    $NEW_USER_ID = $result.Body.user_id
    Write-Host "✓ PASS: Create user endpoint" -ForegroundColor Green
    $PASSED++
} else {
    Write-Host "✗ FAIL: Create user endpoint (expected 201, got $($result.StatusCode))" -ForegroundColor Red
    $FAILED++
}
Write-Host ""

# Test 11: Update User
Write-Host "Test 11: Update User" -ForegroundColor Yellow
$updateData = @{
    first_name = "John Updated"
    last_name = "Doe Updated"
} | ConvertTo-Json

$result = Invoke-ApiCall -Method "PUT" -Endpoint "/users/$USER_ID" -Data $updateData -Token $USER_TOKEN
if ($result.StatusCode -eq 200 -and $result.Body.first_name -eq "John Updated") {
    Write-Host "✓ PASS: Update user endpoint" -ForegroundColor Green
    $PASSED++
} else {
    Write-Host "✗ FAIL: Update user endpoint (expected 200, got $($result.StatusCode))" -ForegroundColor Red
    $FAILED++
}
Write-Host ""

# Test 12: Update User with Invalid Data
Write-Host "Test 12: Update User with Invalid Data (empty body)" -ForegroundColor Yellow
$result = Invoke-ApiCall -Method "PUT" -Endpoint "/users/$USER_ID" -Data "{}" -Token $USER_TOKEN
if ($result.StatusCode -eq 400) {
    Write-Host "✓ PASS: Update user with invalid data (should return 400)" -ForegroundColor Green
    $PASSED++
} else {
    Write-Host "✗ FAIL: Update user with invalid data (expected 400, got $($result.StatusCode))" -ForegroundColor Red
    $FAILED++
}
Write-Host ""

# Test 13: Delete User (Admin Only)
Write-Host "Test 13: Delete User (Admin Only - should succeed)" -ForegroundColor Yellow
$result = Invoke-ApiCall -Method "DELETE" -Endpoint "/users/$NEW_USER_ID" -Token $ADMIN_TOKEN
if ($result.StatusCode -eq 200) {
    Write-Host "✓ PASS: Delete user as admin" -ForegroundColor Green
    $PASSED++
} else {
    Write-Host "✗ FAIL: Delete user as admin (expected 200, got $($result.StatusCode))" -ForegroundColor Red
    $FAILED++
}
Write-Host ""

# Test 14: Delete User (Regular User - should fail)
Write-Host "Test 14: Delete User (Regular User - should fail with 403)" -ForegroundColor Yellow
$result = Invoke-ApiCall -Method "DELETE" -Endpoint "/users/$USER_ID" -Token $USER_TOKEN
if ($result.StatusCode -eq 403) {
    Write-Host "✓ PASS: Delete user as regular user (should return 403)" -ForegroundColor Green
    $PASSED++
} else {
    Write-Host "✗ FAIL: Delete user as regular user (expected 403, got $($result.StatusCode))" -ForegroundColor Red
    $FAILED++
}
Write-Host ""

# Test 15: Register Duplicate Email
Write-Host "Test 15: Register Duplicate Email (should fail)" -ForegroundColor Yellow
$result = Invoke-ApiCall -Method "POST" -Endpoint "/register" -Data $registerData
if ($result.StatusCode -eq 409) {
    Write-Host "✓ PASS: Register duplicate email (should return 409)" -ForegroundColor Green
    $PASSED++
} else {
    Write-Host "✗ FAIL: Register duplicate email (expected 409, got $($result.StatusCode))" -ForegroundColor Red
    $FAILED++
}
Write-Host ""

# Test 16: Register with Invalid Role
Write-Host "Test 16: Register with Invalid Role" -ForegroundColor Yellow
$invalidRoleData = @{
    first_name = "Test"
    last_name = "User"
    email = "invalid@test.com"
    password = "password123"
    phone_number = "+1234567893"
    role = "invalid_role"
} | ConvertTo-Json

$result = Invoke-ApiCall -Method "POST" -Endpoint "/register" -Data $invalidRoleData
if ($result.StatusCode -eq 400) {
    Write-Host "✓ PASS: Register with invalid role (should return 400)" -ForegroundColor Green
    $PASSED++
} else {
    Write-Host "✗ FAIL: Register with invalid role (expected 400, got $($result.StatusCode))" -ForegroundColor Red
    $FAILED++
}
Write-Host ""

# Summary
Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "  Test Summary" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Passed: $PASSED" -ForegroundColor Green
Write-Host "Failed: $FAILED" -ForegroundColor Red
Write-Host "Total: $($PASSED + $FAILED)"
Write-Host ""

if ($FAILED -eq 0) {
    Write-Host "All tests passed! ✓" -ForegroundColor Green
    exit 0
} else {
    Write-Host "Some tests failed. Please review the output above." -ForegroundColor Red
    exit 1
}

