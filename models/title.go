package models

import (
	"log"

	"github.com/jmoiron/sqlx"
)

//Title struc is the model for a law Title
type Title struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Chapters []Chapter `json:"chapters"`
	LawID    int64     `json:"lawID"`
}

//AddChapter adds parsed chapter data to parsed law object
func (t *Title) AddChapter(chapter Chapter) []Chapter {
	t.Chapters = append(t.Chapters, chapter)
	return t.Chapters
}

//CreateTitle Adds a Chapter to the DB
func (t *Title) CreateTitle(db *sqlx.DB) (int64, error) {
	q := "INSERT INTO Title(name,law_id) VALUES($1,$2)"
	result, err := db.Exec(q, t.Name, t.LawID)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, nil
	}
	return lastInsertedID, nil
}

//GetTitles read all Titles from DB
func (t *Title) GetTitles(db *sqlx.DB) ([]Title, error) {
	q := "SELECT ID,name, law_id FROM Title"
	rows, err := db.Query(q)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var titles []Title
	for rows.Next() {
		if err := rows.Scan(&t.ID, &t.Name, &t.LawID); err != nil {
			log.Println(err)
			return nil, err
		}
		titles = append(titles, *t)
	}
	return titles, nil
}
