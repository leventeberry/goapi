package factories

import (
	"github.com/leventeberry/goapi/repositories"
	"gorm.io/gorm"
)

// RepositoryFactory creates repository instances
// Implements Factory Pattern for repository creation
type RepositoryFactory struct {
	db *gorm.DB
}

// NewRepositoryFactory creates a new repository factory
func NewRepositoryFactory(db *gorm.DB) *RepositoryFactory {
	return &RepositoryFactory{
		db: db,
	}
}

// CreateUserRepository creates a UserRepository instance
func (f *RepositoryFactory) CreateUserRepository() repositories.UserRepository {
	return repositories.NewUserRepository(f.db)
}

