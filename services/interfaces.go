package services

import (
	"context"

	"github.com/leventeberry/goapi/models"
	"github.com/leventeberry/goapi/middleware"
)

// UserService defines the interface for user business logic
type UserService interface {
	CreateUser(ctx context.Context, input *CreateUserInput) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	UpdateUser(ctx context.Context, id int, input *UpdateUserInput) (*models.User, error)
	DeleteUser(ctx context.Context, id int) error
	ValidateRole(role string) bool
}

// AuthService defines the interface for authentication business logic
type AuthService interface {
	Login(email, password string) (*models.User, *middleware.Authentication, error)
	Register(input *RegisterInput) (*models.User, *middleware.Authentication, error)
	ValidateCredentials(email, password string) (*models.User, error)
}

