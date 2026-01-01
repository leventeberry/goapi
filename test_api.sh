#!/bin/bash

# API Test Script
# Tests all endpoints of the GoAPI

BASE_URL="http://localhost:8080"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
PASSED=0
FAILED=0

# Helper function to print test results
test_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ PASS${NC}: $2"
        ((PASSED++))
    else
        echo -e "${RED}✗ FAIL${NC}: $2"
        ((FAILED++))
    fi
}

# Helper function to make API calls
api_call() {
    METHOD=$1
    ENDPOINT=$2
    DATA=$3
    TOKEN=$4
    
    if [ -z "$TOKEN" ]; then
        if [ -z "$DATA" ]; then
            curl -s -w "\n%{http_code}" -X $METHOD "$BASE_URL$ENDPOINT" -H "Content-Type: application/json"
        else
            curl -s -w "\n%{http_code}" -X $METHOD "$BASE_URL$ENDPOINT" -H "Content-Type: application/json" -d "$DATA"
        fi
    else
        if [ -z "$DATA" ]; then
            curl -s -w "\n%{http_code}" -X $METHOD "$BASE_URL$ENDPOINT" -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN"
        else
            curl -s -w "\n%{http_code}" -X $METHOD "$BASE_URL$ENDPOINT" -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d "$DATA"
        fi
    fi
}

echo "=========================================="
echo "  GoAPI Comprehensive Test Suite"
echo "=========================================="
echo ""

# Check if server is running
echo "Checking if server is running..."
SERVER_CHECK=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/")
if [ "$SERVER_CHECK" != "200" ]; then
    echo -e "${RED}Error: Server is not running at $BASE_URL${NC}"
    echo "Please start the server first: make run"
    exit 1
fi
echo -e "${GREEN}Server is running${NC}"
echo ""

# Test 1: Health Check
echo "Test 1: Health Check (GET /)"
RESPONSE=$(api_call "GET" "/")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')
if [[ "$HTTP_CODE" == "200" ]] && [[ "$BODY" == *"Welcome"* ]]; then
    test_result 0 "Health check endpoint"
else
    test_result 1 "Health check endpoint (expected 200, got $HTTP_CODE)"
fi
echo ""

# Test 2: Register User (Regular User)
echo "Test 2: Register User (user role)"
REGISTER_DATA='{"first_name":"John","last_name":"Doe","email":"john.doe@test.com","password":"password123","phone_number":"+1234567890","role":"user"}'
RESPONSE=$(api_call "POST" "/register" "$REGISTER_DATA")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')
if [[ "$HTTP_CODE" == "200" ]] && [[ "$BODY" == *"jwt_token"* ]]; then
    USER_TOKEN=$(echo "$BODY" | grep -o '"jwt_token":"[^"]*' | cut -d'"' -f4)
    USER_ID=$(echo "$BODY" | grep -o '"id":[0-9]*' | cut -d':' -f2)
    test_result 0 "Register user endpoint"
else
    test_result 1 "Register user endpoint (expected 200, got $HTTP_CODE)"
    echo "Response: $BODY"
fi
echo ""

# Test 3: Register Admin User
echo "Test 3: Register Admin User"
ADMIN_DATA='{"first_name":"Admin","last_name":"User","email":"admin@test.com","password":"adminpass123","phone_number":"+1234567891","role":"admin"}'
RESPONSE=$(api_call "POST" "/register" "$ADMIN_DATA")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')
if [[ "$HTTP_CODE" == "200" ]] && [[ "$BODY" == *"jwt_token"* ]]; then
    ADMIN_TOKEN=$(echo "$BODY" | grep -o '"jwt_token":"[^"]*' | cut -d'"' -f4)
    ADMIN_ID=$(echo "$BODY" | grep -o '"id":[0-9]*' | cut -d':' -f2)
    test_result 0 "Register admin user endpoint"
else
    test_result 1 "Register admin user endpoint (expected 200, got $HTTP_CODE)"
    echo "Response: $BODY"
fi
echo ""

# Test 4: Login
echo "Test 4: Login"
LOGIN_DATA='{"email":"john.doe@test.com","password":"password123"}'
RESPONSE=$(api_call "POST" "/login" "$LOGIN_DATA")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')
if [[ "$HTTP_CODE" == "200" ]] && [[ "$BODY" == *"jwt_token"* ]]; then
    LOGIN_TOKEN=$(echo "$BODY" | grep -o '"jwt_token":"[^"]*' | cut -d'"' -f4)
    test_result 0 "Login endpoint"
else
    test_result 1 "Login endpoint (expected 200, got $HTTP_CODE)"
    echo "Response: $BODY"
fi
echo ""

# Test 5: Login with Invalid Credentials
echo "Test 5: Login with Invalid Credentials"
INVALID_LOGIN='{"email":"john.doe@test.com","password":"wrongpassword"}'
RESPONSE=$(api_call "POST" "/login" "$INVALID_LOGIN")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')
if [[ "$HTTP_CODE" == "401" ]]; then
    test_result 0 "Login with invalid credentials (should return 401)"
else
    test_result 1 "Login with invalid credentials (expected 401, got $HTTP_CODE)"
fi
echo ""

# Test 6: Get All Users (Authenticated)
echo "Test 6: Get All Users (Authenticated)"
RESPONSE=$(api_call "GET" "/users" "" "$USER_TOKEN")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')
if [[ "$HTTP_CODE" == "200" ]] && [[ "$BODY" == *"john.doe@test.com"* ]]; then
    test_result 0 "Get all users endpoint"
else
    test_result 1 "Get all users endpoint (expected 200, got $HTTP_CODE)"
fi
echo ""

# Test 7: Get All Users (Unauthenticated)
echo "Test 7: Get All Users (Unauthenticated - should fail)"
RESPONSE=$(api_call "GET" "/users")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
if [[ "$HTTP_CODE" == "401" ]]; then
    test_result 0 "Get all users without authentication (should return 401)"
else
    test_result 1 "Get all users without authentication (expected 401, got $HTTP_CODE)"
fi
echo ""

# Test 8: Get User by ID
echo "Test 8: Get User by ID"
RESPONSE=$(api_call "GET" "/users/$USER_ID" "" "$USER_TOKEN")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')
if [[ "$HTTP_CODE" == "200" ]] && [[ "$BODY" == *"john.doe@test.com"* ]]; then
    test_result 0 "Get user by ID endpoint"
else
    test_result 1 "Get user by ID endpoint (expected 200, got $HTTP_CODE)"
fi
echo ""

# Test 9: Get Non-Existent User
echo "Test 9: Get Non-Existent User"
RESPONSE=$(api_call "GET" "/users/99999" "" "$USER_TOKEN")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
if [[ "$HTTP_CODE" == "404" ]]; then
    test_result 0 "Get non-existent user (should return 404)"
else
    test_result 1 "Get non-existent user (expected 404, got $HTTP_CODE)"
fi
echo ""

# Test 10: Create User (Authenticated)
echo "Test 10: Create User (Authenticated)"
CREATE_DATA='{"first_name":"Jane","last_name":"Smith","email":"jane.smith@test.com","password":"password123","phone_number":"+1234567892","role":"user"}'
RESPONSE=$(api_call "POST" "/users" "$CREATE_DATA" "$USER_TOKEN")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')
if [[ "$HTTP_CODE" == "201" ]] && [[ "$BODY" == *"jane.smith@test.com"* ]]; then
    NEW_USER_ID=$(echo "$BODY" | grep -o '"user_id":[0-9]*' | cut -d':' -f2)
    test_result 0 "Create user endpoint"
else
    test_result 1 "Create user endpoint (expected 201, got $HTTP_CODE)"
    echo "Response: $BODY"
fi
echo ""

# Test 11: Update User
echo "Test 11: Update User"
UPDATE_DATA='{"first_name":"John Updated","last_name":"Doe Updated"}'
RESPONSE=$(api_call "PUT" "/users/$USER_ID" "$UPDATE_DATA" "$USER_TOKEN")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')
if [[ "$HTTP_CODE" == "200" ]] && [[ "$BODY" == *"John Updated"* ]]; then
    test_result 0 "Update user endpoint"
else
    test_result 1 "Update user endpoint (expected 200, got $HTTP_CODE)"
fi
echo ""

# Test 12: Update User with Invalid Data
echo "Test 12: Update User with Invalid Data (empty body)"
RESPONSE=$(api_call "PUT" "/users/$USER_ID" "{}" "$USER_TOKEN")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
if [[ "$HTTP_CODE" == "400" ]]; then
    test_result 0 "Update user with invalid data (should return 400)"
else
    test_result 1 "Update user with invalid data (expected 400, got $HTTP_CODE)"
fi
echo ""

# Test 13: Delete User (Admin Only)
echo "Test 13: Delete User (Admin Only - should succeed)"
RESPONSE=$(api_call "DELETE" "/users/$NEW_USER_ID" "" "$ADMIN_TOKEN")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')
if [[ "$HTTP_CODE" == "200" ]] && [[ "$BODY" == *"deleted successfully"* ]]; then
    test_result 0 "Delete user as admin"
else
    test_result 1 "Delete user as admin (expected 200, got $HTTP_CODE)"
fi
echo ""

# Test 14: Delete User (Regular User - should fail)
echo "Test 14: Delete User (Regular User - should fail with 403)"
RESPONSE=$(api_call "DELETE" "/users/$USER_ID" "" "$USER_TOKEN")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
if [[ "$HTTP_CODE" == "403" ]]; then
    test_result 0 "Delete user as regular user (should return 403)"
else
    test_result 1 "Delete user as regular user (expected 403, got $HTTP_CODE)"
fi
echo ""

# Test 15: Register Duplicate Email
echo "Test 15: Register Duplicate Email (should fail)"
RESPONSE=$(api_call "POST" "/register" "$REGISTER_DATA")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
if [[ "$HTTP_CODE" == "409" ]]; then
    test_result 0 "Register duplicate email (should return 409)"
else
    test_result 1 "Register duplicate email (expected 409, got $HTTP_CODE)"
fi
echo ""

# Test 16: Register with Invalid Role
echo "Test 16: Register with Invalid Role"
INVALID_ROLE_DATA='{"first_name":"Test","last_name":"User","email":"invalid@test.com","password":"password123","phone_number":"+1234567893","role":"invalid_role"}'
RESPONSE=$(api_call "POST" "/register" "$INVALID_ROLE_DATA")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
if [[ "$HTTP_CODE" == "400" ]]; then
    test_result 0 "Register with invalid role (should return 400)"
else
    test_result 1 "Register with invalid role (expected 400, got $HTTP_CODE)"
fi
echo ""

# Test 17: Register with Invalid Email Format
echo "Test 17: Register with Invalid Email Format"
INVALID_EMAIL_DATA='{"first_name":"Test","last_name":"User","email":"invalid-email","password":"password123","phone_number":"+1234567894","role":"user"}'
RESPONSE=$(api_call "POST" "/register" "$INVALID_EMAIL_DATA")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
if [[ "$HTTP_CODE" == "400" ]]; then
    test_result 0 "Register with invalid email format (should return 400)"
else
    test_result 1 "Register with invalid email format (expected 400, got $HTTP_CODE)"
fi
echo ""

# Test 18: Register with Short Password
echo "Test 18: Register with Short Password"
SHORT_PASSWORD_DATA='{"first_name":"Test","last_name":"User","email":"shortpass@test.com","password":"short","phone_number":"+1234567895","role":"user"}'
RESPONSE=$(api_call "POST" "/register" "$SHORT_PASSWORD_DATA")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
if [[ "$HTTP_CODE" == "400" ]]; then
    test_result 0 "Register with short password (should return 400)"
else
    test_result 1 "Register with short password (expected 400, got $HTTP_CODE)"
fi
echo ""

# Summary
echo ""
echo "=========================================="
echo "  Test Summary"
echo "=========================================="
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"
echo "Total: $((PASSED + FAILED))"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed! ✓${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed. Please review the output above.${NC}"
    exit 1
fi

