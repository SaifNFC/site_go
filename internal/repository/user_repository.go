package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"demo/internal/models"
)

var ErrUserNotFound = errors.New("utilisateur introuvable")

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return fmt.Errorf("création utilisateur: %w", err)
	}
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("recherche utilisateur par email: %w", err)
	}
	return &user, nil
}
