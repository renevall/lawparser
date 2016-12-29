package postgres

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

//TODO: USE ENV IN PRODUCTION
const (
	DBUser     = "OfficeAdmin"
	DBPassword = "office123"
	DBName     = "Office_Dev"
	DBPort     = "5432"
	DBHost     = "localhost"
)

// InitDB initializes the DB
func InitDB() *sqlx.DB {
	// TODO: Use config files
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		DBUser, DBPassword, DBName, DBHost, DBPort)
	db, err := sqlx.Connect("postgres", dbinfo)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}
