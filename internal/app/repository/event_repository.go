package repository

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/AkbarFikri/BREECE-BE/internal/app/entity"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/model"

)

type EventRepository interface {
	FindById(id string) (entity.Event, error)
	FindAllPublic(page int) ([]entity.Event, error)
	FindWithFilter(params model.FilterParam) ([]entity.Event, error)
	FindByOrganizer(id string) ([]entity.Event, error)
	FindForBooking(id string) (entity.Event, error)
	Update(event entity.Event) error
	UpdateFailurePayment(id string) error
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

// FindById implements EventRepository.
func (r *eventRepository) FindById(id string) (entity.Event, error) {
	var event entity.Event
	if err := r.db.Where("id = ?", id).First(&event).Error; err != nil {
		return event, err
	}
	return event, nil
}

func (r *eventRepository) FindByOrganizer(id string) ([]entity.Event, error) {
	var events []entity.Event

	if err := r.db.Where("organize_by = ?", id).Find(&events).Error; err != nil {
		return events, err
	}

	return events, nil
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

	perPage := 10

	sql = fmt.Sprintf("%s OFFSET %d LIMIT %d", sql, (params.Page-1)*perPage, perPage)

	if err := r.db.Raw(sql).Scan(&events).Error; err != nil {
		return events, err
	}

	return events, nil
}

func (r *eventRepository) FindForBooking(id string) (entity.Event, error) {
	var event entity.Event

	tx := r.db.Begin()

	//TODO ???
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&event).Error; err != nil {
		return event, err
	}

	// TODO Perbaikin Logic tidak pada tempatnya...
	if event.TicketQty == 0 {
		tx.Rollback()
		return event, errors.New("ticket is sold out")
	}

	event.TicketQty = event.TicketQty - 1

	if err := tx.Save(&event).Error; err != nil {
		tx.Rollback()
		return event, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return event, errors.New("failed to commit transaction")
	}

	return event, nil
}

// Insert implements EventRepository.
func (r *eventRepository) Insert(event entity.Event) error {
	if err := r.db.Create(&event).Error; err != nil {
		return err
	}
	return nil
}

func (r *eventRepository) UpdateFailurePayment(id string) error {
	var event entity.Event

	tx := r.db.Begin()

	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&event).Error; err != nil {
		tx.Rollback()
		return err
	}

	event.TicketQty = event.TicketQty + 1

	if err := tx.Save(&event).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
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
