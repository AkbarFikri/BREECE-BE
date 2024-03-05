package repository

import (
	"gorm.io/gorm"

	"github.com/AkbarFikri/BREECE-BE/internal/app/entity"

)

type UserRepository interface {
	FindById(id string) (entity.User, error)
	FindByEmail(email string) (entity.User, error)
	Insert(user entity.User) error
	Update(user entity.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(DB *gorm.DB) UserRepository {
	return &userRepository{
		db: DB,
	}
}

// FindByEmail implements UserRepository.
func (r *userRepository) FindByEmail(email string) (entity.User, error) {
	var user entity.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

// FindById implements UserRepository.
func (r *userRepository) FindById(id string) (entity.User, error) {
	var user entity.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

// Insert implements UserRepository.
func (r *userRepository) Insert(user entity.User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

// Update implements UserRepository.
func (r *userRepository) Update(user entity.User) error {
	if err := r.db.Where("id = ?", user.ID).Save(&user).Error; err != nil {
		return err
	}
	return nil
}
