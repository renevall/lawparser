package database

import (
	"fmt"
	"log"
	"time"

	"bitbucket.org/reneval/lawparser/models"

	"github.com/jmoiron/sqlx"
)

type orm interface {
	makeInsertSQl() string
}

//TODO: USE ENV IN PRODUCTION
const (
	DB_USER     = "Penshiru"
	DB_PASSWORD = "DevPass123"
	DB_NAME     = "Penshiru_Dev"
	DB_PORT     = "5432"
	DB_HOST     = "localhost"
)

//NewDB starts db
func NewDB() (*sqlx.DB, error) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME, DB_HOST, DB_PORT)
	db, err := sqlx.Open("postgres", dbinfo)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return db, nil

}

//InsertLawToDB inserts all parsed law to DB
func InsertLawToDB(db *sqlx.DB, law *models.Law) error {
	start := time.Now()
	lawID, err := law.CreateLaw(db)
	if err != nil {
		log.Println(err)
		return nil
	}
	for _, title := range law.Titles {
		title.LawID = lawID
		titleID, err := title.CreateTitle(db)
		if err != nil {
			log.Println(err)
			return nil
		}
		for _, chapter := range title.Chapters {
			chapter.TitleID = titleID
			chapterID, err := chapter.CreateChapter(db)
			if err != nil {
				log.Println(err)
				return nil
			}

			tx, err := db.Beginx()
			if err != nil {
				log.Fatal(err)
			}
			for _, article := range chapter.Articles {
				article.ChapterID = chapterID
				err := article.CreateArticle(db, tx)
				if err != nil {
					log.Println(err)
					return nil
				}
			}
			tx.Commit()

		}
	}
	elapsed := time.Since(start)
	log.Println("Inserting data to db took: ", elapsed)
	return nil
}
