package models

import (
	"database/sql"
	"log"
)

//Title struc is the model for a law Title
type Title struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Chapters []Chapter `json:"chapters"`
	LawID    int       `json:"lawID"`
}

//AddChapter adds parsed chapter data to parsed law object
func (t *Title) AddChapter(chapter Chapter) []Chapter {
	t.Chapters = append(t.Chapters, chapter)
	return t.Chapters
}

//CreateTitle Adds a Chapter to the DB
func (t *Title) CreateTitle(db *sql.DB) error {
	q := "INSERT INTO Title(name,law_id) VALUES($1,$2)"
		if _, err := db.Exec(q, t.Name,t.LawID); err != nil {
			log.Println(err)
			return err
		}
	return nil
}

//GetTitles read all Titles from DB
func (t *Title) GetTitles(db *sql.DB) ([]Title, error) {
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
		titles = append(titles,*t)
	}
	return titles, nil
}
