package domain

type Book struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	LawID    int64   `json:"lawID"`
	Titles   []Title `json:"titles"`
	Reviewed bool    `json:"reviewed"`
}

// AddTitle adds parsed titles data to parsed law object
func (b *Book) AddTitle(title Title) []Title {
	b.Titles = append(b.Titles, title)
	return b.Titles
}
