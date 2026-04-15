package storage

import (
	"database/sql"
	"fmt"
	"mian/pkg/invoice"
	"mian/pkg/invoiceheader"
	"mian/pkg/invoiceitem"
)

// PsqlInvoice used for work with postgrs - invoice
type PsqlInvoice struct {
	db           *sql.DB
	sorageHeader invoiceheader.Storage
	sorageItem   invoiceitem.Storage
}

// NewPsqlInvoice return a new pointer of PsqlInvoice
func NewPsqlInvoice(db *sql.DB, h invoiceheader.Storage, i invoiceitem.Storage) *PsqlInvoice {
	return &PsqlInvoice{
		db:           db,
		sorageHeader: h,
		sorageItem:   i,
	}
}

// Create implement the invoice.Storage interface
func (p *PsqlInvoice) Create(m *invoice.Model) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	if err := p.sorageHeader.CreateTx(tx, m.Header); err != nil {
		tx.Rollback()
		return fmt.Errorf("Header: %w", err)
	}
	fmt.Printf("Invoice header created successfully with ID: %d \n", m.Header.ID)

	if err := p.sorageItem.CreateTx(tx, m.Header.ID, m.Items); err != nil {
		tx.Rollback()
		return fmt.Errorf("Items: %w", err)
	}
	fmt.Printf("Items created: %d\n", len(m.Items))
	return tx.Commit()
}
