package domain

//Title struc is the model for a law Title
type Title struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Chapters []Chapter `json:"chapters"`
	LawID    int64     `json:"lawID"`
	Reviewed bool      `json:"reviewed"`
}

//AddChapter adds parsed chapter data to parsed law object
func (t *Title) AddChapter(chapter Chapter) []Chapter {
	t.Chapters = append(t.Chapters, chapter)
	return t.Chapters
}
