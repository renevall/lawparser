package domain

import (
	"database/sql"
)

//Title struc is the model for a law Title
type Title struct {
	ID       int64         `json:"id" db:"title_id"`
	Name     string        `json:"name"`
	Chapters []Chapter     `json:"chapters"`
	LawID    int64         `json:"lawID" db:"law_id"`
	BookID   sql.NullInt64 `json:"bookID" db:"book_id"`
	Reviewed bool          `json:"reviewed"`
}

type TitleStore interface {
	CreateTitle() (int64, error)
}

//AddChapter adds parsed chapter data to parsed law object
func (t *Title) AddChapter(chapter Chapter) []Chapter {
	t.Chapters = append(t.Chapters, chapter)
	return t.Chapters
}
