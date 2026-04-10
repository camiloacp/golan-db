package storage

import (
	"database/sql"
	"fmt"
)

const (
	psqlMigrateInvoiceItem = `
	CREATE TABLE IF NOT EXISTS invoice_items (
		id SERIAL NOT NULL PRIMARY KEY,
		invoice_header_id INT NOT NULL,
		product_id INT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		updated_at TIMESTAMP,
		CONSTRAINT invoice_items_invoice_header_id_fk FOREIGN KEY 
		(invoice_header_id) REFERENCES invoice_headers(id) ON UPDATE RESTRICT ON DELETE RESTRICT,
		CONSTRAINT invoice_items_product_id_fk FOREIGN KEY 
		(product_id) REFERENCES products(id) ON UPDATE RESTRICT ON DELETE RESTRICT
	)`
)

// PsqlInvoiceItem used for work with postgrs - invoiceItem
type PsqlInvoiceItem struct {
	db *sql.DB
}

// NewPsqlInvoiceItem return a new pointer of PsqlInvoiceItem
func NewPsqlInvoiceItem(db *sql.DB) *PsqlInvoiceItem {
	return &PsqlInvoiceItem{db: db}
}

// Migrate implement the invoiceitem.Storage interface
func (p *PsqlInvoiceItem) Migration() error {
	stmt, err := p.db.Prepare(psqlMigrateInvoiceItem)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	fmt.Println("Migration invoiceitem success")
	return nil
}
