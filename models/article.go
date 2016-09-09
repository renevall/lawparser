package models

import (
	"log"

	"github.com/jmoiron/sqlx"
)

//Article Holds the article model and his methods
type Article struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Text      string `json:"text"`
	ChapterID int64  `json:"chapterID"`
	LawID     int64  `json:"lawID"`
	Reviewed  bool   `json:"reviewed"`
}

//CreateArticle Adds an Article to the DB
func (a *Article) CreateArticle(db *sqlx.DB, tx *sqlx.Tx) error {
	q := `Insert INTO Article(name,text,chapter_id,law_id, reviewed) 
	VALUES($1,$2,$3,$4,$5)`

	if tx != nil {
		if _, err := tx.Exec(q, a.Name, a.Text, a.ChapterID, a.LawID, a.Reviewed); err != nil {
			log.Println(err)
			return err
		}
	} else {
		if _, err := db.Exec(q, a.Name, a.Text, a.ChapterID, a.LawID, a.Reviewed); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

//GetArticles read all articles from DB
func (a *Article) GetArticles(db *sqlx.DB) ([]Article, error) {
	q := "SELECT article_id,name,text,chapter_id,law_id,reviewed` FROM Article"
	rows, err := db.Query(q)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var articles []Article
	for rows.Next() {
		if err := rows.Scan(&a.ID, &a.Name, &a.Text, &a.ChapterID,
			&a.LawID, &a.Reviewed); err != nil {
			log.Println(err)
			return nil, err
		}
		articles = append(articles, *a)
	}
	return articles, nil
}
