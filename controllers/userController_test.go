package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Success Test for the GetUsers function
func TestGetUsers_Success(t *testing.T) {
	// Setup mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening stub database: %s", err)
	}
	defer db.Close()

	// Prepare mock rows
	rows := sqlmock.NewRows([]string{
		"user_id", "first_name", "last_name", "email", "password_hash", "phone_number", "role", "created_at", "updated_at",
	}).AddRow(
		"1", "John", "Doe", "test@test.com", "hashedpassword", "1234567890", "user", "2021-01-01", "2021-01-01",
	)

	// Expect query with stricter match
	mock.ExpectQuery("^SELECT \\* FROM users$").WillReturnRows(rows)

	// Define router and endpoint
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/users", GetUsers(db))

	// Simulate GET request
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatalf("failed to create HTTP request: %v", err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))


	var users []User
	if err := json.Unmarshal(w.Body.Bytes(), &users); err != nil {
		t.Fatalf("failed to parse response JSON: %v", err)
	}

	// Validate user data
	assert.Len(t, users, 1)
	assert.Equal(t, "John", users[0].FirstName)
	assert.Equal(t, "Doe", users[0].LastName)
	assert.Equal(t, "test@test.com", users[0].Email)
	assert.Equal(t, "1234567890", users[0].PhoneNum)
	assert.Equal(t, "user", users[0].Role)

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// Error Test for the GetUsers function
func TestGetUsers_ErrorCases(t *testing.T) {
	// Setup mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening stub database: %s", err)
	}
	defer db.Close()

	// Simulate a DB error
	mock.ExpectQuery("^SELECT \\* FROM users$").
		WillReturnError(fmt.Errorf("database connection error"))

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/users", GetUsers(db))

	// Create request and record response
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Validate response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "Failed to query the database")

	// Ensure mock expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

// Success Test for the GetUser function
func TestGetUser_Success(t *testing.T) {
	// Setup mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening stub database: %s", err)
	}
	defer db.Close()

	// Prepare mock rows
	rows := sqlmock.NewRows([]string{
		"user_id", "first_name", "last_name", "email", "password_hash", "phone_number", "role", "created_at", "updated_at",
	}).AddRow(
		1, "John", "Doe", "test@test.com", "hashedpassword", "1234567890", "user", "2021-01-01", "2021-01-01",
	)

	// Expect query with stricter match
	mock.ExpectQuery("^SELECT \\* FROM users WHERE user_id = \\?$").WithArgs(1).WillReturnRows(rows)

	// Define router and endpoint
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/users/:id", GetUser(db))

	// Simulate GET request
	req, err := http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatalf("failed to create HTTP request: %v", err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var user User
	if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
		t.Fatalf("failed to parse response JSON: %v", err)
	}

	// Validate user data
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, "test@test.com", user.Email)
	assert.Equal(t, "1234567890", user.PhoneNum)
	assert.Equal(t, "user", user.Role)
	assert.Equal(t, "2021-01-01", user.CreatedAt)
	assert.Equal(t, "2021-01-01", user.UpdateAt)

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// Error Test for the GetUser function
func TestGetUser_ErrorCases(t *testing.T) {
	// Setup mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening stub database: %s", err)
	}
	defer db.Close()

	// Expect query with int64 arg
	mock.ExpectQuery("^SELECT \\* FROM users WHERE user_id = \\?$").
		WithArgs(int64(3)).
		WillReturnError(fmt.Errorf("database connection error"))

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/users/:id", GetUser(db))

	// Create request and record response
	req, err := http.NewRequest("GET", "/users/3", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Validate response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "Failed to fetch user from the database")

	// Ensure mock expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

// Sucsess Test for the CreateUser function
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
        "first_name": "John",
        "last_name": "Doe",
        "email": "john@example.com",
        "password_hash": "hashedpassword",
        "phone_number": "1234567890",
        "role": "user"
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

// Error Test for the CreateUser function
func TestCreateUser_ErrorCases(t *testing.T) {
    tests := []struct {
        name           string
        payload        string
        setupMock      func(mock sqlmock.Sqlmock)
        expectedStatus int
        expectedBody   string
    }{
        {
            name:           "invalid JSON",
            payload:        `{"invalid": "json",`,
            setupMock:      func(mock sqlmock.Sqlmock) {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   "Failed to parse the request body",
        },
        {
            name:    "database error",
            payload: `{"FirstName": "Jane", "LastName": "Doe", "Email": "jane@example.com", "PassHash": "hash", "PhoneNum": "0987654321", "Role": "user"}`,
            setupMock: func(mock sqlmock.Sqlmock) {
                mock.ExpectExec("INSERT INTO users").WillReturnError(sql.ErrConnDone)
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   "Failed to insert the user into the database",
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            // Setup mock DB
            db, mock, err := sqlmock.New()
            if err != nil {
                t.Fatalf("error opening stub database: %v", err)
            }
            defer db.Close()

            tc.setupMock(mock)

            gin.SetMode(gin.TestMode)
            router := gin.Default()
            router.POST("/users", CreateUser(db))

            req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(tc.payload))
            if err != nil {
                t.Fatalf("Could not create request: %v", err)
            }
            req.Header.Set("Content-Type", "application/json")

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatus, w.Code)
            assert.Contains(t, w.Body.String(), tc.expectedBody)

            if err := mock.ExpectationsWereMet(); err != nil {
                t.Errorf("unfulfilled expectations: %s", err)
            }
        })
    }
}

// Success Test for the UpdateUser function
// func TestUpdateUser_Success(t *testing.T) {
// 	// Setup: Create a sqlmock database
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	// Set expected DB behavior for success case
// 	mock.ExpectExec("UPDATE users").WithArgs("John", "Doe", "testty@test.com", "hashedpassword", "1234567890", "user", 1).
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	// Create a Gin engine in test mode
// 	gin.SetMode(gin.TestMode)
// 	router := gin.Default()
// 	router.PUT("/users/:id", UpdateUser(db))

// 	// Create a sample user payload
// 	jsonPayload := `{
// 		"first_name": "John",
// 		"last_name": "Doe",
// 		"email": "testty@test.com",
// 		"password_hash": "hashedpassword",
// 		"phone_number": "1234567890",
// 		"role": "user"
// 	}`
// 	req, err := http.NewRequest("PUT", "/users/1", bytes.NewBufferString(jsonPayload))
// 	if err != nil {
// 		t.Fatalf("Could not create request: %v", err)
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	// Record the response
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	// Assert on the response
// 	assert.Equal(t, http.StatusOK, w.Code)
// 	assert.Contains(t, w.Body.String(), "User updated successfully")

// 	// Ensure all expectations were met
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 	}
// }

// Error Test for the UpdateUser function
// func TestUpdateUser_ErrorCases(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		payload        string
// 		setupMock      func(mock sqlmock.Sqlmock)
// 		expectedStatus int
// 		expectedBody   string
// 	}{
// 		{
// 			name:           "invalid JSON",
// 			payload:        `{"invalid": "json"}`,
// 			setupMock:      func(mock sqlmock.Sqlmock) {},
// 			expectedStatus: http.StatusBadRequest,
// 			expectedBody:   "Failed to parse the request body",
// 		},
// 		{
// 			name:    "database error",
// 			payload: `{"FirstName": "Jane", "LastName": "Doe", "Email": "test@test.com", "PassHash": "hash", "PhoneNum": "0987654321", "Role": "user"}`,
// 			setupMock: func(mock sqlmock.Sqlmock) {
// 				mock.ExpectExec("UPDATE users").WillReturnError(sql.ErrConnDone)
// 			},
// 			expectedStatus: http.StatusInternalServerError,
// 			expectedBody:   "Failed to update the user in the database",
// 		},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Setup mock DB
// 			db, mock, err := sqlmock.New()
// 			if err != nil {
// 				t.Fatalf("error opening stub database: %v", err)
// 			}
// 			defer db.Close()

// 			tc.setupMock(mock)

// 			gin.SetMode(gin.TestMode)
// 			router := gin.Default()
// 			router.PUT("/users/:id", UpdateUser(db))

// 			req, err := http.NewRequest("PUT", "/users/1", bytes.NewBufferString(tc.payload))
// 			if err != nil {
// 				t.Fatalf("Could not create request: %v", err)
// 			}
// 			req.Header.Set("Content-Type", "application/json")

// 			w := httptest.NewRecorder()
// 			router.ServeHTTP(w, req)

// 			assert.Equal(t, tc.expectedStatus, w.Code)
// 			assert.Contains(t, w.Body.String(), tc.expectedBody)

// 			if err := mock.ExpectationsWereMet(); err != nil {
// 				t.Errorf("unfulfilled expectations: %s", err)
// 			}
// 		})
// 	}
// }

