package repositories

import (
	"errors"
	"github.com/leventeberry/goapi/models"
	"gorm.io/gorm"
)

// userRepository implements UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
// Factory function for creating user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create inserts a new user into the database
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByID retrieves a user by their ID
func (r *userRepository) FindByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail retrieves a user by their email
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindAll retrieves all users from the database
func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Update updates an existing user in the database
// Uses Updates() instead of Save() to only update changed fields
func (r *userRepository) Update(user *models.User) error {
	return r.db.Model(user).Updates(user).Error
}

// Delete removes a user from the database
func (r *userRepository) Delete(id int) error {
	result := r.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

// ExistsByEmail checks if a user with the given email exists
func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

