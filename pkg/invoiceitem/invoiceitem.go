package invoiceitem

import (
	"database/sql"
	"time"
)

// Model of invoiceitem
type Model struct {
	ID              uint
	InvoiceHeaderID uint
	ProductID       uint
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Models slice of Model
type Models []*Model

// Storage interface that must be implemented by the storage layer
type Storage interface {
	Migrate() error
	CreateTx(*sql.Tx, uint, Models) error
}

// Service of invoiceitem
type Service struct {
	storage Storage
}

// NewService return a new pointer of Service
func NewService(s Storage) *Service {
	return &Service{storage: s}
}

// Migrate is used to migrate the product table
func (s *Service) Migrate() error {
	return s.storage.Migrate()
}
