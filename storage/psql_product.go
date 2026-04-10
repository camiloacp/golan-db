package storage

import (
	"database/sql"
	"fmt"
	"mian/pkg/product"
)

type scanner interface {
	Scan(dest ...interface{}) error
}

const (
	psqlMigrateProduct = `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL NOT NULL PRIMARY KEY,
		name VARCHAR(25) NOT NULL,
		observations varchar(100),
		price INT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		updated_at TIMESTAMP
	)`

	psqlCreateProduct = `
	INSERT INTO products (name, observations, price, created_at) VALUES ($1, $2, $3, $4) 
	RETURNING id
	`

	psqlGetAllProducts = `
	SELECT id, name, observations, price, created_at, updated_at FROM products
	`

	psqlGetByIDProduct = psqlGetAllProducts + ` WHERE id = $1`

	psqlUpdateProduct = `
	UPDATE products SET name = $1, observations = $2, price = $3, updated_at = $4 WHERE id = $5
	`

	psqlDeleteProduct = `DELETE FROM products WHERE id = $1`
)

// PsqlProduct used for work with postgrs - product
type PsqlProduct struct {
	db *sql.DB
}

// NewPsqlProduct return a new pointer of PsqlProduct
func NewPsqlProduct(db *sql.DB) *PsqlProduct {
	return &PsqlProduct{db: db}
}

// Migrate implement the product.Storage interface
func (p *PsqlProduct) Migration() error {
	stmt, err := p.db.Prepare(psqlMigrateProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	fmt.Println("Migration product success")
	return nil
}

// Create implement the product.Storage interface
func (p *PsqlProduct) Create(m *product.Model) error {
	stmt, err := p.db.Prepare(psqlCreateProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		m.Name,
		stringToNull(m.Observations),
		m.Price,
		m.CreatedAt,
	).Scan(&m.ID)
	if err != nil {
		return err
	}
	fmt.Println("Product created successfully")
	return nil
}

// GetAll implement the product.Storage interface
func (p *PsqlProduct) GetAll() (product.Models, error) {
	stmt, err := p.db.Prepare(psqlGetAllProducts)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ms := make(product.Models, 0)
	for rows.Next() {
		m, err := scanRowProduct(rows)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	fmt.Println("Products fetched successfully")
	return ms, nil
}

// GetById implement the product.Storage interface
func (p *PsqlProduct) GetByID(id uint) (*product.Model, error) {
	stmt, err := p.db.Prepare(psqlGetByIDProduct)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	m, err := scanRowProduct(stmt.QueryRow(id))
	if err != nil {
		return nil, err
	}

	fmt.Println("Product fetched by ID successfully")
	return m, nil
}

// Update implement the product.Storage interface
func (p *PsqlProduct) Update(m *product.Model) error {
	stmt, err := p.db.Prepare(psqlUpdateProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		m.Name,
		stringToNull(m.Observations),
		m.Price,
		timeToNull(m.UpdatedAt),
		m.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("Not exist product with the ID: %d", m.ID)
	}

	fmt.Println("Product updated successfully")
	return nil
}

// Delete implement the product.Storage interface
func (p *PsqlProduct) Delete(id uint) error {
	stmt, err := p.db.Prepare(psqlDeleteProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("Not exist product with the ID: %d", id)
	}

	fmt.Println("Product deleted successfully")
	return nil
}

func scanRowProduct(s scanner) (*product.Model, error) {
	m := &product.Model{}
	observationNull := sql.NullString{}
	updatedAtNull := sql.NullTime{}

	err := s.Scan(
		&m.ID,
		&m.Name,
		&observationNull,
		&m.Price,
		&m.CreatedAt,
		&updatedAtNull,
	)
	if err != nil {
		return &product.Model{}, err
	}

	m.Observations = observationNull.String
	m.UpdatedAt = updatedAtNull.Time

	return m, nil
}
