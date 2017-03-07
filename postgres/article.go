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
func (a *Article) CreateArticle(tx *sqlx.Tx) error {
	q := `Insert INTO Article(name,text,chapter_id,law_id, reviewed) 
	VALUES($1,$2,$3,$4,$5)`

	if tx != nil {
		if _, err := tx.Exec(q, a.Article.Name, a.Article.Text, a.Article.ChapterID, a.Article.LawID, a.Article.Reviewed); err != nil {
			log.Println(err)
			return err
		}
	} else {
		if _, err := a.DB.Exec(q, a.Article.Name, a.Article.Text, a.Article.ChapterID, a.Article.LawID, a.Article.Reviewed); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
