package repository

import (
	"gorm.io/gorm"

	"github.com/AkbarFikri/BREECE-BE/internal/app/entity"

)

type TicketRepository interface {
	FindById(id string) (entity.Ticket, error)
	FindByUserIdAndEventId(userId, eventId string) (entity.Ticket, error)
	Insert(ticket entity.Ticket) error
	Update(ticket entity.Ticket) error
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{
		db: db,
	}
}

// FindById implements TicketRepository.
func (r *ticketRepository) FindById(id string) (entity.Ticket, error) {
	var ticket entity.Ticket
	if err := r.db.Where("id = ?", id).First(&ticket).Error; err != nil {
		return ticket, err
	}
	return ticket, nil
}

// FindByUserIdAndEventId implements TicketRepository.
func (*ticketRepository) FindByUserIdAndEventId(userId string, eventId string) (entity.Ticket, error) {
	panic("unimplemented")
}

// Insert implements TicketRepository.
func (r *ticketRepository) Insert(ticket entity.Ticket) error {
	if err := r.db.Create(&ticket).Error; err != nil {
		return err
	}
	return nil
}

// Update implements TicketRepository.
func (r *ticketRepository) Update(ticket entity.Ticket) error {
	if err := r.db.Save(&ticket).Error; err != nil {
		return err
	}
	return nil
}
