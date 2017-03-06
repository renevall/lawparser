package postgres

import (
	"log"

	"bitbucket.org/reneval/lawparser/domain"
)

type Book struct {
	DB   *DB
	Book *domain.Book
}

func (b *Book) createBook() (int64, error) {
	q := "INSERT INTO Book(name,law_id,reviewed) VALUES($1,$2,$3) RETURNING book_id"
	var id int64
	err := b.DB.QueryRow(q, b.Book.Name, b.Book.LawID, b.Book.Reviewed).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return id, nil
}
