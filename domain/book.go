package domain

type Book struct {
	ID       int64   `json:"id" db:"book_id"`
	Name     string  `json:"name"`
	Text     string  `json:"text" db:"text"`
	LawID    int64   `json:"lawID" db:"law_id"`
	Titles   []Title `json:"titles"`
	Reviewed bool    `json:"reviewed"`
}

// AddTitle adds parsed titles data to parsed law object
func (b *Book) AddTitle(title Title) []Title {
	b.Titles = append(b.Titles, title)
	return b.Titles
}
