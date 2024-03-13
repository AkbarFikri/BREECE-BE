package repository

import (
	"fmt"
	"strings"

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
	if err := r.db.Where("is_public = true").Limit(perPage).Offset((page - 1) * perPage).Find(&events).Error; err != nil {
		return events, err
	}
	return events, nil
}

// FindWithFilter implements EventRepository.
func (r *eventRepository) FindWithFilter(params model.FilterParam) ([]entity.Event, error) {
	var events []entity.Event

	sql := "SELECT * FROM events"

	if params.Search != "" {
		sql = fmt.Sprintf("%s WHERE title LIKE '%%%s%%' OR description LIKE '%%%s%%'", sql, params.Search, params.Search)
	}

	if params.Category != "" {
		if strings.Contains(sql, "WHERE") {
			sql = fmt.Sprintf("%s AND category_id = '%s'", sql, params.Category)
		} else {
			sql = fmt.Sprintf("%s WHERE category_id = '%s'", sql, params.Category)
		}

	}

	if params.Place != "" {
		if strings.Contains(sql, "WHERE") {
			sql = fmt.Sprintf("%s OR tempat LIKE '%%%s%%'", sql, params.Place)
		} else {
			sql = fmt.Sprintf("%s WHERE tempat LIKE '%%%s%%'", sql, params.Place)
		}
	}

	if params.Date != "" {
		// data, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", params.Date)
		if strings.Contains(sql, "WHERE") {
			sql = fmt.Sprintf("%s AND CAST(date as DATE) = '%s'", sql, params.Date)
		} else {
			sql = fmt.Sprintf("%s WHERE CAST(date as DATE) = '%s'", sql, params.Date)
		}

	}

	if params.IsPublic {
		if strings.Contains(sql, "WHERE") {
			sql = fmt.Sprintf("%s AND is_public = true", sql)
		} else {
			sql = fmt.Sprintf("%s WHERE is_public = true", sql)
		}
	}

	// if params.Sort != "" {
	// 	sql = fmt.Sprintf("%s ORDER BY date %s", sql, params.Sort)
	// }

	perPage := 10

	// if sql != "" {
	// 	if err := r.db.Where(sql).Limit(perPage).Offset((params.Page - 1) * perPage).Find(&events).Error; err != nil {
	// 		return events, err
	// 	}
	// } else {
	// 	if err := r.db.Limit(perPage).Offset((params.Page - 1) * perPage).Find(&events).Error; err != nil {
	// 		return events, err
	// 	}
	// }

	sql = fmt.Sprintf("%s OFFSET %d LIMIT %d", sql, (params.Page-1)*perPage, perPage)

	if err := r.db.Raw(sql).Scan(&events).Error; err != nil {
		return events, err
	}

	return events, nil
}

// Insert implements EventRepository.
func (r *eventRepository) Insert(event entity.Event) error {
	if err := r.db.Create(&event).Error; err != nil {
		return err
	}
	return nil
}

// Update implements EventRepository.
func (r *eventRepository) Update(event entity.Event) error {
	if err := r.db.Save(event).Error; err != nil {
		return err
	}
	return nil
}
