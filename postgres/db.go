package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

//TODO: USE ENV IN PRODUCTION
const (
	DBUser     = "Penshiru"
	DBPassword = "DevPass123"
	DBName     = "Penshiru_Dev"
	DBPort     = "5432"
	DBHost     = "localhost"
)

// DB is a wrapper for *sqlx.DB
type DB struct {
	*sqlx.DB
}

// InitDB initializes the DB
func InitDB() (*DB, error) {
	// TODO: Use config files
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		DBUser, DBPassword, DBName, DBHost, DBPort)
	db, err := sqlx.Connect("postgres", dbinfo)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
