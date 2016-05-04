package models

import (
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

//Law struct with most methods.
type Law struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Titles       []Title   `json:"titles"`
	ApprovalDate time.Time `json:"approvalDate"`
	PublishDate  time.Time `json:"publishDate"`
	Journal      string    `json:"journal"`
	Intro        string    `json:"intro"`
}

//AddTitle adds parsed title data to parsed law object
func (law *Law) AddTitle(title Title) []Title {
	law.Titles = append(law.Titles, title)
	return law.Titles
}

//CreateLaw Adds a Law to the DB
func (law *Law) CreateLaw(db *sqlx.DB) (int64, error) {
	q := "INSERT INTO LAW(name,approval_date,publish_date,journal,intro) VALUES($1,$2,$3,$4,$5)"

	result, err := db.Exec(q, law.Name, law.ApprovalDate, law.PublishDate, law.Journal, law.Intro)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastInsertedID, nil

}

//GetLaws read all articles from DB
func (law *Law) GetLaws(db *sqlx.DB) ([]Law, error) {
	q := "SELECT ID,name,approval_date,publish_date,journal,intro FROM Law"
	rows, err := db.Query(q)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var laws []Law
	for rows.Next() {
		if err := rows.Scan(&law.ID, &law.Name, &law.ApprovalDate, &law.PublishDate, &law.PublishDate,
			&law.Journal, &law.Intro); err != nil {
			log.Println(err)
			return nil, err
		}
		laws = append(laws, *law)
	}
	return laws, nil
}

//GetFullLaw return a mapped law object with all the other associations
func (law *Law) GetFullLaw(db *sqlx.DB, id int) error {
	q := "SELECT ID,name,approval_date,publish_date,journal,intro FROM Law WHERE id=?"
	err := db.QueryRow(q, id).Scan(&law.ID, &law.Name, &law.ApprovalDate, &law.PublishDate, &law.PublishDate,
		&law.Journal, &law.Intro)
	if err != nil {
		log.Println(err)
		return err
	}
	q = "SELECT ID,name, law_id FROM Title WHERE law_id=?"
	rows, err := db.Query(q, law.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	var titles []int
	for rows.Next() {
		var t Title
		if err := rows.Scan(&t.ID, &t.Name, &t.LawID); err != nil {
			log.Println(err)
			return err
		}
		//add to main object
		law.AddTitle(t)
	}
	q = "SELECT ID,name, title_id FROM Chapters WHERE title_id IN (?" + strings.Repeat(",?", len(titles)-1) + ")"
	var chapters []Chapter

	return nil
}
