package controllers

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestCreateUser_Success(t *testing.T) {
    // Setup: Create a sqlmock database
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    // Set expected DB behavior for success case
    mock.ExpectExec("INSERT INTO users").WithArgs("John", "Doe", "john@example.com", "hashedpassword", "1234567890", "user").
        WillReturnResult(sqlmock.NewResult(1, 1))

    // Create a Gin engine in test mode
    gin.SetMode(gin.TestMode)
    router := gin.Default()
    router.POST("/users", CreateUser(db))

    // Create a sample user payload
    jsonPayload := `{
        "FirstName": "John",
        "LastName": "Doe",
        "Email": "john@example.com",
        "PassHash": "hashedpassword",
        "PhoneNum": "1234567890",
        "Role": "user"
    }`
    req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(jsonPayload))
    if err != nil {
        t.Fatalf("Could not create request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Record the response
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Assert on the response
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Body.String(), "User created successfully")

    // Ensure all expectations were met
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectations: %s", err)
    }
}
