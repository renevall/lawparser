package models

import (
	"database/sql"
	"log"
)

//Article Holds the article model and his methods
type Article struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Text      string `json:"text"`
	ChapterID int    `json:"chapterID"`
}

//CreateArticle Adds an Article to the DB
func (a *Article) CreateArticle(db *sql.DB, tx *sql.Tx) error {
	q := "Insert INTO Article(name,text,chapter_id) VALUES($1,$2,$3)"

	if tx != nil {
		if _, err := tx.Exec(q, a.Name, a.Text, a.ChapterID); err != nil {
			log.Println(err)
			return err
		}
	} else {
		if _, err := db.Exec(q, a.Name, a.Text, a.ChapterID); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

//GetArticles read all articles from DB
func (a *Article) GetArticles(db *sql.DB) ([]Article, error) {
	q := "SELECT ID,name, text, chapter_id FROM Article"
	rows, err := db.Query(q)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var articles []Article
	for rows.Next() {
		if err := rows.Scan(&a.ID, &a.Name, &a.Text, &a.ChapterID); err != nil {
			log.Println(err)
			return nil, err
		}
		articles = append(articles, *a)
	}
	return articles, nil
}
