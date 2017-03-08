package domain

import "time"

type Law struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	ApprovalDate time.Time `json:"approvalDate"`
	PublishDate  time.Time `json:"publishDate"`
	Journal      string    `json:"journal"`
	Intro        string    `json:"intro"`
	Reviewed     bool      `json:"reviewed"`
	Revision     int       `json:"rev"`
	Books        []Book    `json:"books"`
	Titles       []Title   `json:"titles"`
	Chapters     []Chapter `json:"chapters"`
	Articles     []Article `json:"articles"`
	Init         string    `json:"init"`
}

type LawStore interface {
	GetLaws() ([]Law, error)
	InsertLawDB(law *Law) error
	CreateLaw() (int64, error)
}

//AddTitle adds parsed title data to parsed law object
func (law *Law) AddTitle(title Title) []Title {
	law.Titles = append(law.Titles, title)
	return law.Titles
}

func (law *Law) AddBook(book Book) []Book {
	law.Books = append(law.Books, book)
	return law.Books
}

//AddBook adds parsed Bookdata to law object

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
