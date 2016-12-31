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
	Titles       []Title   `json:"titles"`
	Chapters     []Chapter `json:"chapters"`
	Articles     []Article `json:"articles"`
}

type LawStore interface {
	GetLaws() ([]Law, error)
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
