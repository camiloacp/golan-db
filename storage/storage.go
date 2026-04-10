package storage

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	once sync.Once
)

func NewPostgresDB() {
	once.Do(func() {
		var err error

		db, err = sql.Open("postgres", "postgres://golang_db_user:golang_db_password@localhost:7530/godb?sslmode=disable")
		if err != nil {
			log.Fatalf("Error opening database: %v", err)
		}

		if err = db.Ping(); err != nil {
			log.Fatalf("Error pinging database: %v", err)
		}
		fmt.Println("Database Postgres connected successfully")
	})
}

// Pool returns a unique instance of db
func Pool() *sql.DB {
	return db
}

func stringToNull(s string) sql.NullString {
	null := sql.NullString{String: s}
	if null.String != "" {
		null.Valid = true
	}
	return null
}

func timeToNull(t time.Time) sql.NullTime {
	null := sql.NullTime{Time: t}
	if !null.Time.IsZero() {
		null.Valid = true
	}
	return null
}
