package invoiceheader

import (
	"database/sql"
	"time"
)

// Model of invoiceheader
type Model struct {
	ID        uint
	Client    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Storage interface {
	Migration() error
	CreateTx(*sql.Tx, *Model) error
}

// Service of invoiceheader
type Service struct {
	storage Storage
}

// NewService return a new pointer of Service
func NewService(s Storage) *Service {
	return &Service{storage: s}
}

// Migrate is used to migrate the product table
func (s *Service) Migrate() error {
	return s.storage.Migration()
}
