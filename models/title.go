package models

//Title struc is the model for a law Title
type Title struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Chapters []Chapter `json:"chapters"`
	LawID    int       `json:"lawID"`
}

//AddChapter adds parsed chapter data to parsed law object
func (title *Title) AddChapter(chapter Chapter) []Chapter {
	title.Chapters = append(title.Chapters, chapter)
	return title.Chapters
}
