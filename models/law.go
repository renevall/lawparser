package models

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

//Law struct with most methods.
type Law struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	ApprovalDate time.Time `json:"approvalDate"`
	PublishDate  time.Time `json:"publishDate"`
	Journal      string    `json:"journal"`
	Intro        string    `json:"intro"`
	Reviewed     bool      `json:"reviewed"`
	Revision     int       `json:"rev"`
	Titles       []Title   `json:"titles"`
	Chapters     []Chapter `json:"chapters"`
	Articles     []Article `json:"articles"`
}

//TmpLaw hold basic data to access files
type TmpLaw struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

//AddTitle adds parsed title data to parsed law object
func (law *Law) AddTitle(title Title) []Title {
	law.Titles = append(law.Titles, title)
	return law.Titles
}

//AddChapter adds parsed article data to parsed law object
//when there is no title
func (law *Law) AddChapter(chapter Chapter) []Chapter {
	law.Chapters = append(law.Chapters, chapter)
	return law.Chapters
}

//AddArticle adds parsed article data to parsed law object
//when there is no title and no chapter
func (law *Law) AddArticle(article Article) []Article {
	law.Articles = append(law.Articles, article)
	return law.Articles
}

//CreateLaw Adds a Law to the DB
func (law *Law) CreateLaw(db *sqlx.DB) (int64, error) {
	q := `INSERT INTO LAW
		(name,approval_date,publish_date,journal,intro,reviewed, revision) 
		VALUES($1,$2,$3,$4,$5,$6,$7)`

	//TODO Parse Date from txt file
	law.PublishDate = time.Now()
	law.ApprovalDate = law.PublishDate
	result, err := db.Exec(q, law.Name, law.ApprovalDate, law.PublishDate,
		law.Journal, law.Intro, law.Reviewed, law.Revision)

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

//GetLaws read all laws from DB
func (law *Law) GetLaws(db *sqlx.DB) ([]Law, error) {
	q := `SELECT 
	law_id,name,approval_date,publish_date,journal,intro, reviewed, revision
	FROM Law`
	rows, err := db.Query(q)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var laws []Law
	for rows.Next() {
		if err := rows.Scan(&law.ID, &law.Name, &law.ApprovalDate, &law.PublishDate,
			&law.Journal, &law.Intro, &law.Reviewed, &law.Revision); err != nil {
			log.Println(err)
			return nil, err
		}
		laws = append(laws, *law)
	}
	return laws, nil
}

//GetFullLaw return a mapped law object with all the other associations
func (law *Law) GetFullLaw(db *sqlx.DB, id int) error {
	q := `SELECT law_id,name,approval_date,publish_date,journal,intro, reviewed, revision
	FROM Law WHERE law_id=?`
	err := db.QueryRow(q, id).Scan(&law.ID, &law.Name, &law.ApprovalDate, &law.PublishDate,
		&law.Journal, &law.Intro, &law.Reviewed, &law.Revision)
	if err != nil {
		log.Println(err)
		return err
	}

	//titles
	q = "SELECT title_id,name, law_id, reviewed FROM Title WHERE law_id=?"
	rows, err := db.Query(q, law.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	var TitleIDs []int
	for rows.Next() {
		var t Title
		if err := rows.Scan(&t.ID, &t.Name, &t.LawID, &t.Reviewed); err != nil {
			log.Println(err)
			return err
		}
		//add to main object
		law.AddTitle(t)
		TitleIDs = append(TitleIDs, t.ID)
	}

	//chapters
	q = "SELECT chapter_id,name, title_id, law_id, reviewed FROM Chapter WHERE title_id IN (?)"
	var chapters []Chapter
	var chapterIDs []int
	query, args, err := sqlx.In(q, TitleIDs)
	query = db.Rebind(query)
	rows, err = db.Query(query, args...)
	if err != nil {
		log.Println(err)
		return err
	}
	for rows.Next() {
		var c Chapter
		if err := rows.Scan(&c.ID, &c.Name, &c.TitleID, &c.LawID, &c.Reviewed); err != nil {
			log.Println(err)
			return err
		}
		chapters = append(chapters, c)
		chapterIDs = append(chapterIDs, c.ID)
	}
	//articles
	q = `SELECT article_id,name, text, chapter_id, law_id, reviewed 
	FROM Article WHERE chapter_id in (?)`
	var articles []Article
	query, args, err = sqlx.In(q, chapterIDs)
	query = db.Rebind(query)
	rows, err = db.Query(query, args...)
	for rows.Next() {
		var a Article
		if err := rows.Scan(&a.ID, &a.Name, &a.Text, &a.ChapterID,
			&a.LawID, &a.Reviewed); err != nil {
			log.Println(err)
			return err
		}
		articles = append(articles, a)
	}

	//making the final structure
	for t, title := range law.Titles {
		for _, chapter := range chapters {
			if title.ID == int(chapter.TitleID) {
				law.Titles[t].AddChapter(chapter)
			}
		}
	}

	//TODO improve this
	for t, _ := range law.Titles {
		for c, chapter := range law.Titles[t].Chapters {
			for _, article := range articles {
				if int(article.ChapterID) == chapter.ID {
					law.Titles[t].Chapters[c].AddArticle(article)
				}
			}
		}
	}
	return nil
}
