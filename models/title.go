package models

//Title struc is the model for a law Title
type Title struct {
	Name     string    `json:"name"`
	Chapters []Chapter `json:"chapters"`
}

func (title *Title) AddChapter(chapter Chapter) []Chapter {
	title.Chapters = append(title.Chapters, chapter)
	return title.Chapters
}
