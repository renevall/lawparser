package postgres

import (
	"log"

	"bitbucket.org/reneval/lawparser/domain"
	"github.com/jmoiron/sqlx"
)

type Article struct {
	DB      *DB
	Article *domain.Article
}

//CreateArticle Adds an Article to the DB
func (a *Article) CreateArticle(tx *sqlx.Tx) (int64, error) {
	q := `Insert INTO Article(name,text,chapter_id,law_id, reviewed) 
	VALUES($1,$2,$3,$4,$5) RETURNING article_id`

	var id int64

	if tx != nil {
		if err := tx.QueryRow(q, a.Article.Name, a.Article.Text, a.Article.ChapterID, a.Article.LawID,
			a.Article.Reviewed).Scan(&id); err != nil {
			log.Println(err)
			return 0, err
		}
	} else {
		if err := a.DB.QueryRow(q, a.Article.Name, a.Article.Text, a.Article.ChapterID, a.Article.LawID,
			a.Article.Reviewed).Scan(&id); err != nil {
			log.Println(err)
			return 0, err
		}
	}
	return id, nil
}
