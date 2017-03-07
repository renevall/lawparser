package postgres

import (
	"log"
	"time"

	"fmt"

	"bitbucket.org/reneval/lawparser/domain"
)

type Law struct {
	*DB
	Law *domain.Law
}

func (l *Law) GetLaws() ([]domain.Law, error) {
	q := `SELECT 
	law_id,name,approval_date,publish_date,journal,intro, reviewed, revision
	FROM Law`
	rows, err := l.DB.Query(q)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	law := &domain.Law{}
	var laws []domain.Law
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

//InsertLawDB inserts all parsed law to DB
func (l *Law) InsertLawDB(newLaw *domain.Law) error {
	l.Law = newLaw
	start := time.Now()
	lawID, err := l.CreateLaw()
	fmt.Println(lawID)
	if err != nil {
		log.Println(err)
		return nil
	}
	pqTitle := &Title{DB: l.DB}
	pqChapter := &Chapter{DB: l.DB}
	pqArticle := &Article{DB: l.DB}
	for _, title := range l.Law.Titles {
		pqTitle.Title = &title
		pqTitle.Title.LawID = lawID

		titleID, err := pqTitle.CreateTitle()
		if err != nil {
			log.Println(err)
			return nil
		}
		for _, chapter := range title.Chapters {
			pqChapter.Chapter = &chapter
			pqChapter.Chapter.TitleID = titleID
			chapterID, err := pqChapter.CreateChapter()
			if err != nil {
				log.Println(err)
				return nil
			}
			tx, err := l.DB.Beginx()
			if err != nil {
				log.Fatal(err)
			}
			for _, article := range chapter.Articles {
				pqArticle.Article = &article
				pqArticle.Article.ChapterID = chapterID
				err := pqArticle.CreateArticle(tx)
				if err != nil {
					log.Println(err)
					return nil
				}
			}
			tx.Commit()

		}
	}
	elapsed := time.Since(start)
	log.Println("Inserting data to db took: ", elapsed)
	return nil
}

//CreateLaw Adds a Law to the DB
func (l *Law) CreateLaw() (int64, error) {
	q := `INSERT INTO LAW
		(name,approval_date,publish_date,journal,intro,reviewed, revision) 
		VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING law_id`

	//TODO Parse Date from txt file
	law := l.Law
	var id int64

	law.PublishDate = time.Now()
	law.ApprovalDate = law.PublishDate
	err := l.DB.QueryRow(q, law.Name, law.ApprovalDate, law.PublishDate,
		law.Journal, law.Intro, law.Reviewed, law.Revision).Scan(&id)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	return id, nil

}
