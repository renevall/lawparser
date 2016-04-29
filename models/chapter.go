package models

import (
	"database/sql"
	"log"
)

//Chapter is the model for a Law chapter
type Chapter struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Articles []Article `json:"articles"`
	TitleID  int       `json:"titleID"`
}

//AddArticle adds parsed article data to parsed law object
func (c *Chapter) AddArticle(article Article) []Article {
	c.Articles = append(c.Articles, article)
	return c.Articles
}

//CreateChapter Adds a Chapter to the DB
func (c *Chapter) CreateChapter(db *sql.DB) error {
	q := "INSERT INTO Chapter(name,title_id) VALUES($1,$2)"
		if _, err := db.Exec(q, c.Name,c.TitleID); err != nil {
			log.Println(err)
			return err
		}
	return nil
}

//GetChapters read all Chapters from DB
func (c *Chapter) GetChapters(db *sql.DB) ([]Chapter, error) {
	q := "SELECT ID,name, title_id FROM Chapters"
	rows, err := db.Query(q)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var chapters []Chapter
	for rows.Next() {
		if err := rows.Scan(&c.ID, &c.Name, &c.TitleID); err != nil {
			log.Println(err)
			return nil, err
		}
		chapters = append(chapters,*c)
	}
	return chapters, nil
}
