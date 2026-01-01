package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/leventeberry/goapi/container"
	"github.com/leventeberry/goapi/initializers"
	"github.com/leventeberry/goapi/routes"
)

var (
	testRouter    *gin.Engine
	testContainer *container.Container
	userToken     string
	adminToken    string
	userID        int
	adminID       int
)

// Setup test environment
func TestMain(m *testing.M) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Initialize database (use test database if available)
	initializers.Init()

	// Create container
	testContainer = container.NewContainer(initializers.DB, initializers.GetCacheClient())

	// Create router
	testRouter = gin.New()
	routes.SetupRoutes(testRouter, testContainer)

	// Run tests
	code := m.Run()

	// Cleanup if needed
	os.Exit(code)
}

// Helper function to make requests
func makeRequest(method, url string, body interface{}, token string) (*httptest.ResponseRecorder, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w, nil
}

// Test 1: Health Check
func TestHealthCheck(t *testing.T) {
	w, err := makeRequest("GET", "/", nil, "")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["message"] != "Welcome!" {
		t.Errorf("Expected 'Welcome!', got '%v'", response["message"])
	}

	fmt.Println("✓ Health check test passed")
}

// Test 2: Register User
func TestRegisterUser(t *testing.T) {
	registerData := map[string]interface{}{
		"first_name":   "John",
		"last_name":    "Doe",
		"email":        "john.doe@test.com",
		"password":     "password123",
		"phone_number": "+1234567890",
		"role":         "user",
	}

	w, err := makeRequest("POST", "/register", registerData, "")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		return
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	tokenData, ok := response["token"].(map[string]interface{})
	if !ok {
		t.Fatal("Token not found in response")
	}

	userToken = tokenData["jwt_token"].(string)
	userID = int(response["user"].(map[string]interface{})["id"].(float64))

	fmt.Println("✓ Register user test passed")
}

// Test 3: Register Admin
func TestRegisterAdmin(t *testing.T) {
	adminData := map[string]interface{}{
		"first_name":   "Admin",
		"last_name":    "User",
		"email":        "admin@test.com",
		"password":     "adminpass123",
		"phone_number": "+1234567891",
		"role":         "admin",
	}

	w, err := makeRequest("POST", "/register", adminData, "")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		return
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	tokenData := response["token"].(map[string]interface{})
	adminToken = tokenData["jwt_token"].(string)
	adminID = int(response["user"].(map[string]interface{})["id"].(float64))

	fmt.Println("✓ Register admin test passed")
}

// Test 4: Login
func TestLogin(t *testing.T) {
	loginData := map[string]interface{}{
		"email":    "john.doe@test.com",
		"password": "password123",
	}

	w, err := makeRequest("POST", "/login", loginData, "")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		return
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["token"] == nil {
		t.Error("Token not found in response")
	}

	fmt.Println("✓ Login test passed")
}

// Test 5: Login with Invalid Credentials
func TestLoginInvalidCredentials(t *testing.T) {
	loginData := map[string]interface{}{
		"email":    "john.doe@test.com",
		"password": "wrongpassword",
	}

	w, err := makeRequest("POST", "/login", loginData, "")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}

	fmt.Println("✓ Login with invalid credentials test passed")
}

// Test 6: Get All Users (Authenticated)
func TestGetAllUsers(t *testing.T) {
	if userToken == "" {
		t.Skip("User token not available")
	}

	w, err := makeRequest("GET", "/users", nil, userToken)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		return
	}

	var users []interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &users); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(users) == 0 {
		t.Error("Expected at least one user")
	}

	fmt.Println("✓ Get all users test passed")
}

// Test 7: Get All Users (Unauthenticated)
func TestGetAllUsersUnauthenticated(t *testing.T) {
	w, err := makeRequest("GET", "/users", nil, "")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}

	fmt.Println("✓ Get all users unauthenticated test passed")
}

// Test 8: Get User by ID
func TestGetUserByID(t *testing.T) {
	if userToken == "" || userID == 0 {
		t.Skip("User token or ID not available")
	}

	url := fmt.Sprintf("/users/%d", userID)
	w, err := makeRequest("GET", url, nil, userToken)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		return
	}

	var user map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if user["email"] != "john.doe@test.com" {
		t.Errorf("Expected email 'john.doe@test.com', got '%v'", user["email"])
	}

	fmt.Println("✓ Get user by ID test passed")
}

// Test 9: Get Non-Existent User
func TestGetNonExistentUser(t *testing.T) {
	if userToken == "" {
		t.Skip("User token not available")
	}

	w, err := makeRequest("GET", "/users/99999", nil, userToken)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	fmt.Println("✓ Get non-existent user test passed")
}

// Test 10: Create User (Authenticated)
func TestCreateUser(t *testing.T) {
	if userToken == "" {
		t.Skip("User token not available")
	}

	createData := map[string]interface{}{
		"first_name":   "Jane",
		"last_name":    "Smith",
		"email":        "jane.smith@test.com",
		"password":     "password123",
		"phone_number": "+1234567892",
		"role":         "user",
	}

	w, err := makeRequest("POST", "/users", createData, userToken)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
		return
	}

	var user map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if user["email"] != "jane.smith@test.com" {
		t.Errorf("Expected email 'jane.smith@test.com', got '%v'", user["email"])
	}

	fmt.Println("✓ Create user test passed")
}

// Test 11: Update User
func TestUpdateUser(t *testing.T) {
	if userToken == "" || userID == 0 {
		t.Skip("User token or ID not available")
	}

	updateData := map[string]interface{}{
		"first_name": "John Updated",
		"last_name":  "Doe Updated",
	}

	url := fmt.Sprintf("/users/%d", userID)
	w, err := makeRequest("PUT", url, updateData, userToken)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		return
	}

	var user map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if user["first_name"] != "John Updated" {
		t.Errorf("Expected first_name 'John Updated', got '%v'", user["first_name"])
	}

	fmt.Println("✓ Update user test passed")
}

// Test 12: Delete User (Admin Only)
func TestDeleteUserAsAdmin(t *testing.T) {
	if adminToken == "" {
		t.Skip("Admin token not available")
	}

	// First create a user to delete
	createData := map[string]interface{}{
		"first_name":   "ToDelete",
		"last_name":    "User",
		"email":        "todelete@test.com",
		"password":     "password123",
		"phone_number": "+1234567893",
		"role":         "user",
	}

	w, err := makeRequest("POST", "/users", createData, adminToken)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if w.Code != http.StatusCreated {
		t.Skipf("Could not create user for deletion test: %d", w.Code)
		return
	}

	var user map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	deleteID := int(user["user_id"].(float64))
	url := fmt.Sprintf("/users/%d", deleteID)

	w, err = makeRequest("DELETE", url, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	fmt.Println("✓ Delete user as admin test passed")
}

// Test 13: Delete User (Regular User - should fail)
func TestDeleteUserAsRegularUser(t *testing.T) {
	if userToken == "" || userID == 0 {
		t.Skip("User token or ID not available")
	}

	url := fmt.Sprintf("/users/%d", userID)
	w, err := makeRequest("DELETE", url, nil, userToken)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d. Body: %s", w.Code, w.Body.String())
	}

	fmt.Println("✓ Delete user as regular user test passed")
}

// Test 14: Register Duplicate Email
func TestRegisterDuplicateEmail(t *testing.T) {
	registerData := map[string]interface{}{
		"first_name":   "Duplicate",
		"last_name":    "User",
		"email":        "john.doe@test.com", // Already registered
		"password":     "password123",
		"phone_number": "+1234567894",
		"role":         "user",
	}

	w, err := makeRequest("POST", "/register", registerData, "")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409, got %d. Body: %s", w.Code, w.Body.String())
	}

	fmt.Println("✓ Register duplicate email test passed")
}

// Test 15: Register with Invalid Role
func TestRegisterInvalidRole(t *testing.T) {
	registerData := map[string]interface{}{
		"first_name":   "Invalid",
		"last_name":    "Role",
		"email":        "invalidrole@test.com",
		"password":     "password123",
		"phone_number": "+1234567895",
		"role":         "invalid_role",
	}

	w, err := makeRequest("POST", "/register", registerData, "")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d. Body: %s", w.Code, w.Body.String())
	}

	fmt.Println("✓ Register with invalid role test passed")
}
