package postgres

import (
	"log"
	"time"

	"fmt"

	"bitbucket.org/reneval/lawparser/domain"
	"github.com/jmoiron/sqlx"
)

type Law struct {
	*DB
	Law *domain.Law
}

func (l *Law) AutoComplete(query string) ([]string, error) {
	fmt.Println(query)
	q := `SELECT name FROM law WHERE unaccent(name) ILIKE unaccent($1)`
	data := []string{}
	hits := []struct {
		Name string
	}{}
	err := l.Select(&hits, q, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	for _, hit := range hits {
		data = append(data, hit.Name)
	}

	fmt.Println(data)
	return data, nil

}

func (l *Law) GetLaws() ([]domain.Law, error) {
	q := `SELECT 
	law_id,name,number,approval_date,publish_date,journal,intro, reviewed, revision
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
		if err := rows.Scan(&law.ID, &law.Name, &law.Number, &law.ApprovalDate, &law.PublishDate,
			&law.Journal, &law.Intro, &law.Reviewed, &law.Revision); err != nil {
			log.Println(err)
			return nil, err
		}
		laws = append(laws, *law)
	}
	return laws, nil
}

//GetLaw returns one Law record found on the db
func (l *Law) GetLaw(id string) (domain.Law, error) {
	law := domain.Law{}
	err := l.Get(&law, "SELECT * FROM law WHERE law_id=$1", id)
	if err != nil {
		return law, err
	}

	//books
	books := []domain.Book{}
	err = l.Select(&books, "SELECT * FROM book WHERE law_id=$1 ORDER BY book_id", id)
	if err != nil {
		return law, err
	}
	if len(books) > 0 {
		law.Books = books
	}

	//titles
	titles := []domain.Title{}
	err = l.Select(&titles, "SELECT * FROM title WHERE law_id=$1 ORDER BY title_id", id)
	if err != nil {
		return law, err
	}
	if len(law.Books) > 0 {
		for bIndex, book := range law.Books {
			for _, title := range titles {
				if book.ID == title.BookID.Int64 {
					law.Books[bIndex].Titles = append(law.Books[bIndex].Titles, title)

				}
			}
		}
	} else {
		law.Titles = titles
	}

	//chapters
	chapters := []domain.Chapter{}
	err = l.Select(&chapters, "SELECT * FROM chapter WHERE law_id=$1 ORDER BY chapter_id", id)
	if err != nil {
		return law, err
	}
	if len(law.Books) > 0 {
		for bIndex, book := range law.Books {
			for tIndex, title := range book.Titles {
				for _, chapter := range chapters {
					if title.ID == chapter.TitleID {
						law.Books[bIndex].Titles[tIndex].Chapters =
							append(law.Books[bIndex].Titles[tIndex].Chapters, chapter)
					}
				}

			}
		}
	} else {
		if len(law.Titles) > 0 {
			for tIndex, title := range law.Titles {
				for _, chapter := range chapters {
					if title.ID == chapter.TitleID {
						law.Titles[tIndex].Chapters =
							append(law.Titles[tIndex].Chapters, chapter)
					}
				}
			}
		} else {
			law.Chapters = chapters
		}
	}

	//articles
	articles := []domain.Article{}
	err = l.Select(&articles, "SELECT * FROM article WHERE law_id=$1 ORDER BY article_id", id)
	if err != nil {
		return law, err
	}
	if len(law.Books) > 0 {
		for bIndex, book := range law.Books {
			for tIndex, title := range book.Titles {
				for cIndex, chapter := range title.Chapters {
					for _, article := range articles {
						if chapter.ID == article.ChapterID {
							law.Books[bIndex].Titles[tIndex].Chapters[cIndex].Articles =
								append(law.Books[bIndex].Titles[tIndex].Chapters[cIndex].Articles, article)
						}
					}
				}
			}
		}
	} else {
		if len(law.Titles) > 0 {
			for tIndex, title := range law.Titles {
				for cIndex, chapter := range title.Chapters {
					for _, article := range articles {
						if chapter.ID == article.ChapterID {
							law.Titles[tIndex].Chapters[cIndex].Articles =
								append(law.Titles[tIndex].Chapters[cIndex].Articles, article)
						}
					}
				}
			}
		} else {
			law.Articles = articles
		}
	}
	return law, nil
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

	if len(l.Law.Books) > 0 {
		for _, book := range l.Law.Books {
			bookID, err := fillBooks(&book, lawID, l.DB)
			if err != nil {
				return err
			}
			if len(book.Titles) > 0 {
				for _, title := range book.Titles {
					titleID, err := fillTitles(&title, lawID, bookID, l.DB)
					if err != nil {
						return err
					}
					if len(title.Chapters) > 0 {
						for _, chapter := range title.Chapters {
							chapterID, err := fillChapter(&chapter, lawID, titleID, l.DB)
							if err != nil {
								return err
							}
							if len(chapter.Articles) > 0 {
								tx, err := l.DB.Beginx()
								if err != nil {
									return err
								}
								for _, article := range chapter.Articles {
									_, err := fillArticle(&article, lawID, chapterID, l.DB, tx)
									if err != nil {
										return nil
									}
								}
								tx.Commit()
							}
						}
					}
				}
			}
		}
	} else if len(l.Law.Titles) > 0 {
		for _, title := range l.Law.Titles {
			titleID, err := fillTitles(&title, lawID, 0, l.DB)
			if err != nil {
				return err
			}
			if len(title.Chapters) > 0 {
				for _, chapter := range title.Chapters {
					chapterID, err := fillChapter(&chapter, lawID, titleID, l.DB)
					if err != nil {
						return err
					}
					if len(chapter.Articles) > 0 {
						tx, err := l.DB.Beginx()
						if err != nil {
							return err
						}
						for _, article := range chapter.Articles {
							_, err := fillArticle(&article, lawID, chapterID, l.DB, tx)
							if err != nil {
								return nil
							}
						}
						tx.Commit()
					}
				}
			}
		}
	} else if len(l.Law.Articles) > 0 {
		tx, err := l.DB.Beginx()
		if err != nil {
			return err
		}
		for _, article := range l.Law.Articles {
			_, err := fillArticle(&article, lawID, 0, l.DB, tx)
			if err != nil {
				return nil
			}
		}
		tx.Commit()
	}
	elapsed := time.Since(start)
	log.Println("Inserting data to db took: ", elapsed)
	return nil
}

//CreateLaw Adds a Law to the DB
func (l *Law) CreateLaw() (int64, error) {
	q := `INSERT INTO LAW
		(name,number,approval_date,publish_date,journal,intro,reviewed, revision) 
		VALUES($1,$2,$3,$4,$5,$6,$7,$8) RETURNING law_id`

	//TODO Parse Date from txt file
	law := l.Law
	var id int64

	law.PublishDate = time.Now()
	law.ApprovalDate = law.PublishDate
	err := l.DB.QueryRow(q, law.Name, law.Number, law.ApprovalDate, law.PublishDate,
		law.Journal, law.Intro, law.Reviewed, law.Revision).Scan(&id)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	return id, nil

}

func fillBooks(book *domain.Book, lawID int64, db *DB) (int64, error) {
	pgBook := &Book{DB: db}
	pgBook.Book = book
	pgBook.Book.LawID = lawID

	bookID, err := pgBook.createBook()
	if err != nil {
		return 0, err
	}

	return bookID, nil

}

func fillTitles(title *domain.Title, lawID int64, bookID int64, db *DB) (int64, error) {
	pqTitle := &Title{DB: db}
	pqTitle.Title = title
	pqTitle.Title.LawID = lawID
	pqTitle.Title.BookID.Int64 = bookID
	if bookID == 0 {
		pqTitle.Title.BookID.Valid = false
	} else {
		pqTitle.Title.BookID.Valid = true

	}

	titleID, err := pqTitle.CreateTitle()
	if err != nil {
		return 0, nil
	}

	return titleID, nil
}

func fillChapter(chapter *domain.Chapter, lawID int64, titleID int64, db *DB) (int64, error) {
	pqChapter := &Chapter{DB: db}
	pqChapter.Chapter = chapter
	pqChapter.Chapter.LawID = lawID
	pqChapter.Chapter.TitleID = titleID

	chapterID, err := pqChapter.CreateChapter()
	if err != nil {
		return 0, nil
	}

	return chapterID, nil
}

func fillArticle(article *domain.Article, lawID int64, chapterID int64, db *DB, tx *sqlx.Tx) (int64, error) {
	pqArticle := &Article{DB: db}
	pqArticle.Article = article
	pqArticle.Article.LawID = lawID
	pqArticle.Article.ChapterID = chapterID

	articleID, err := pqArticle.CreateArticle(tx)
	if err != nil {
		return 0, nil
	}

	return articleID, nil
}
