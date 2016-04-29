package models

import (
	"time"
	"log"
	"database/sql"
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
func (law *Law) CreateLaw(db *sql.DB) error{
	q:="INSERT INTO LAW(name,approval_date,publish_date,journal,intro) VALUES($1,$2,$3,$4,$5)"
	
	if _,err := db.Exec(q,law.Name,law.ApprovalDate,law.PublishDate,law.Journal,law.Intro); err!=nil{
		log.Println(err)
		return err
	}
	
	return nil
	
}

//GetLaws read all articles from DB
func (law *Law) GetLaws(db *sql.DB) ([]Law, error) {
	q := "SELECT ID,name,approval_date,publish_date,journal,intro FROM Law"
	rows, err := db.Query(q)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var laws []Law
	for rows.Next() {
		if err := rows.Scan(&law.ID, &law.Name, &law.ApprovalDate, &law.PublishDate,&law.PublishDate,
		&law.Journal, &law.Intro); err != nil {
			log.Println(err)
			return nil, err
		}
		laws = append(laws,*law)
	}
	return laws, nil
}
