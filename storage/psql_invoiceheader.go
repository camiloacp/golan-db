package storage

import (
	"database/sql"
	"fmt"
	"mian/pkg/invoiceheader"
)

const (
	psqlMigrateInvoiceHeader = `
	CREATE TABLE IF NOT EXISTS invoice_headers (
		id SERIAL NOT NULL PRIMARY KEY,
		client VARCHAR(25) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		updated_at TIMESTAMP
	)`

	psqlCreateInvoiceHeader = `
	INSERT INTO invoice_headers (client) VALUES ($1)
	RETURNING id, created_at
	`
)

// PsqlInvoiceHeader used for work with postgrs - invoiceHeader
type PsqlInvoiceHeader struct {
	db *sql.DB
}

// NewPsqlInvoiceHeader return a new pointer of PsqlInvoiceHeader
func NewPsqlInvoiceHeader(db *sql.DB) *PsqlInvoiceHeader {
	return &PsqlInvoiceHeader{db: db}
}

// Migrate implement the invoiceheader.Storage interface
func (p *PsqlInvoiceHeader) Migrate() error {
	stmt, err := p.db.Prepare(psqlMigrateInvoiceHeader)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	fmt.Println("Migration invoiceheader success")
	return nil
}

func (p *PsqlInvoiceHeader) CreateTx(tx *sql.Tx, m *invoiceheader.Model) error {
	stmt, err := tx.Prepare(psqlCreateInvoiceHeader)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRow(m.Client).Scan(&m.ID, &m.CreatedAt)
}
