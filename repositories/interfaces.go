package repositories

import "github.com/leventeberry/goapi/models"

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *models.User) error
	FindByID(id int) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindAll() ([]models.User, error)
	Update(user *models.User) error
	Delete(id int) error
	ExistsByEmail(email string) (bool, error)
}

