package postgres

import (
	"log"

	"bitbucket.org/reneval/lawparser/domain"
)

type Title struct {
	DB    *DB
	Title *domain.Title
}

//CreateTitle Adds a Chapter to the DB
func (t *Title) CreateTitle() (int64, error) {
	q := "INSERT INTO Title(name,law_id,reviewed) VALUES($1,$2,$3) RETURNING title_id"
	var id int64
	err := t.DB.QueryRow(q, t.Title.Name, t.Title.LawID, t.Title.Reviewed).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return id, nil
}

// //GetTitles read all Titles from DB
// func (t *Title) GetTitles() ([]domain.Title, error) {
// 	q := "SELECT ID,name, law_id, reviewed FROM Title"
// 	rows, err := t.DB.Query(q)
// 	defer rows.Close()
// 	if err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}

// 	var titles []Title
// 	for rows.Next() {
// 		if err := rows.Scan(&t.ID, &t.Name, &t.LawID, &t.Reviewed); err != nil {
// 			log.Println(err)
// 			return nil, err
// 		}
// 		titles = append(titles, *t)
// 	}
// 	return titles, nil
// }
