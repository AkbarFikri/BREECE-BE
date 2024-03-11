package repository

import (
	"gorm.io/gorm"

	"github.com/AkbarFikri/BREECE-BE/internal/app/entity"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/model"
)

type EventRepository interface {
	FindAllPublic(page int) ([]entity.Event, error)
	FindWithFilter(params model.FilterParam) ([]entity.Event, error)
	Update(event entity.Event) error
	Insert(event entity.Event) error
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{
		db: db,
	}
}

// FindAllPublic implements EventRepository.
func (r *eventRepository) FindAllPublic(page int) ([]entity.Event, error) {
	var events []entity.Event
	perPage := 10
	if err := r.db.Where("is_public = true").Offset(page * perPage).Limit(perPage).Find(&events).Error; err != nil {
		return events, err
	}
	return events, nil
}

// FindWithFilter implements EventRepository.
func (*eventRepository) FindWithFilter(params model.FilterParam) ([]entity.Event, error) {
	panic("unimplemented")
}

// Insert implements EventRepository.
func (*eventRepository) Insert(event entity.Event) error {
	panic("unimplemented")
}

// Update implements EventRepository.
func (*eventRepository) Update(event entity.Event) error {
	panic("unimplemented")
}
