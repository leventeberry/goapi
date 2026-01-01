package factories

import (
	"github.com/leventeberry/goapi/cache"
	"github.com/leventeberry/goapi/repositories"
	"github.com/leventeberry/goapi/services"
)

// ServiceFactory creates service instances
// Implements Factory Pattern for service creation
type ServiceFactory struct {
	userRepo repositories.UserRepository
	cache    cache.Cache
}

// NewServiceFactory creates a new service factory
func NewServiceFactory(userRepo repositories.UserRepository, cacheClient cache.Cache) *ServiceFactory {
	return &ServiceFactory{
		userRepo: userRepo,
		cache:    cacheClient,
	}
}

// CreateUserService creates a UserService instance
func (f *ServiceFactory) CreateUserService() services.UserService {
	return services.NewUserService(f.userRepo, f.cache)
}

// CreateAuthService creates an AuthService instance
func (f *ServiceFactory) CreateAuthService() services.AuthService {
	return services.NewAuthService(f.userRepo)
}

