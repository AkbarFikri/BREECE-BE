package repository

import (
	"gorm.io/gorm"

	"github.com/AkbarFikri/BREECE-BE/internal/app/entity"
)

type InvoiceRepository interface {
	FindById(id string) (entity.Invoice, error)
	FindByUserId(id string) ([]entity.Invoice, error)
	FindByEventId(id string) ([]entity.Invoice, error)
	Insert(invoice entity.Invoice) error
	Update(invoice entity.Invoice) error
}

type invoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(DB *gorm.DB) InvoiceRepository {
	return &invoiceRepository{
		db: DB,
	}
}

// FindById implements InvoiceRepository.
func (r *invoiceRepository) FindById(id string) (entity.Invoice, error) {
	var invoice entity.Invoice

	if err := r.db.Where("id = ?", id).First(&invoice).Error; err != nil {
		return invoice, err
	}

	return invoice, nil
}

// FindByEventId implements InvoiceRepository.
func (r *invoiceRepository) FindByEventId(id string) ([]entity.Invoice, error) {
	var invoices []entity.Invoice
	if err := r.db.Where("event_id = ?", id).Find(&invoices).Error; err != nil {
		return invoices, err
	}
	return invoices, nil
}

// FindByUserId implements InvoiceRepository.
func (r *invoiceRepository) FindByUserId(id string) ([]entity.Invoice, error) {
	var invoices []entity.Invoice
	if err := r.db.Preload("Event").Where("user_id = ?", id).Find(&invoices).Error; err != nil {
		return invoices, err
	}
	return invoices, nil
}

// Insert implements InvoiceRepository.
func (r *invoiceRepository) Insert(invoice entity.Invoice) error {
	if err := r.db.Create(&invoice).Error; err != nil {
		return err
	}
	return nil
}

// Update implements InvoiceRepository.
func (r *invoiceRepository) Update(invoice entity.Invoice) error {
	if err := r.db.Save(&invoice).Error; err != nil {
		return err
	}
	return nil
}
