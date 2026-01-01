package container

import (
	"github.com/leventeberry/goapi/factories"
	"github.com/leventeberry/goapi/repositories"
	"github.com/leventeberry/goapi/services"
	"gorm.io/gorm"
)

// Container holds all application dependencies
// Implements Dependency Injection Container pattern
type Container struct {
	DB                *gorm.DB
	RepositoryFactory *factories.RepositoryFactory
	ServiceFactory    *factories.ServiceFactory
	UserRepository    repositories.UserRepository
	UserService       services.UserService
	AuthService       services.AuthService
}

// NewContainer creates and initializes a new dependency injection container
// Uses Factory Pattern to create all dependencies
func NewContainer(db *gorm.DB) *Container {
	// Create repository factory
	repoFactory := factories.NewRepositoryFactory(db)

	// Create repositories
	userRepo := repoFactory.CreateUserRepository()

	// Create service factory
	serviceFactory := factories.NewServiceFactory(userRepo)

	// Create services
	userService := serviceFactory.CreateUserService()
	authService := serviceFactory.CreateAuthService()

	return &Container{
		DB:                db,
		RepositoryFactory: repoFactory,
		ServiceFactory:    serviceFactory,
		UserRepository:    userRepo,
		UserService:       userService,
		AuthService:       authService,
	}
}

