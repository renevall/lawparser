package db

import (
	"database/sql"
	"log"
    "bitbucket.com/reneval/lawparser/config"
)

//NewDB checks for a db, if not creates a new one.
func NewDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "penshiru.sqlite")
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
