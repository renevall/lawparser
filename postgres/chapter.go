package postgres

import (
	"log"

	"bitbucket.org/reneval/lawparser/domain"
)

type Chapter struct {
	DB      *DB
	Chapter *domain.Chapter
}

//CreateChapter Adds a Chapter to the DB
func (c *Chapter) CreateChapter() (int64, error) {
	var id int64
	q := "INSERT INTO Chapter(name,title_id,law_id, reviewed) VALUES($1,$2,$3,$4) RETURNING chapter_id"
	err := c.DB.QueryRow(q, c.Chapter.Name, c.Chapter.TitleID, c.Chapter.LawID, c.Chapter.Reviewed).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return id, nil
}
