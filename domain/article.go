package domain

//Article Holds the article model and his methods
type Article struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Text      string `json:"text"`
	ChapterID int64  `json:"chapterID"`
	LawID     int64  `json:"lawID"`
	Reviewed  bool   `json:"reviewed"`
}
