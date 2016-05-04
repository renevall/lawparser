package db

import (
	"log"

	"github.com/jmoiron/sqlx"

	"bitbucket.org/reneval/lawparser/config"
)

type orm interface {
	makeInsertSQl() string
}

//NewDB checks for a db, if not creates a new one.
func NewDB() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", "penshiru.sqlite")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, table := range config.InitSqls {
		_, err := db.Exec(table)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	return db, nil

}
