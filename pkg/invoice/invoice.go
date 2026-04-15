package invoice

import (
	"mian/pkg/invoiceheader"
	"mian/pkg/invoiceitem"
)

// Model of invoice
type Model struct {
	Header *invoiceheader.Model
	Items  invoiceitem.Models
}

// Storage interface that must be implemented by the storage layer
type Storage interface {
	Create(*Model) error
}

// Service of invoice
type Service struct {
	storage Storage
}

// NewService return a new pointer of Service
func NewService(s Storage) *Service {
	return &Service{s}
}

// Create a new invoice
func (s *Service) Create(m *Model) error {
	return s.storage.Create(m)
}
