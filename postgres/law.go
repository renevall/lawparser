package postgres

import (
	"log"

	"bitbucket.org/reneval/lawparser/domain"
)

type Law struct {
	*DB
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
