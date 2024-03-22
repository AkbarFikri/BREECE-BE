package repository

import (
	"gorm.io/gorm"

	"github.com/AkbarFikri/BREECE-BE/internal/app/entity"

)

type CategoryRepository interface {
	FindAll() ([]entity.Category, error)
	Insert(category entity.Category) error
	Update(category entity.Category) error
	Delete(category entity.Category) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

// Delete implements CategoryRepository.
func (r *categoryRepository) Delete(category entity.Category) error {
	if err := r.db.Where("id = ?", category.ID).Delete(&category).Error; err != nil {
		return err
	}

	return nil
}

// FindAll implements CategoryRepository.
func (r *categoryRepository) FindAll() ([]entity.Category, error) {
	var categories []entity.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return categories, err
	}

	return categories, nil
}

// Insert implements CategoryRepository.
func (r *categoryRepository) Insert(category entity.Category) error {
	if err := r.db.Create(&category).Error; err != nil {
		return err
	}
	return nil
}

// Update implements CategoryRepository.
func (r *categoryRepository) Update(category entity.Category) error {
	if err := r.db.Where("id = ?", category.ID).Updates(&category).Error; err != nil {
		return err
	}
	return nil
}
