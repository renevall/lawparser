package models

//Law struct with most methods.
type Law struct {
	Name   string  `json:"name"`
	Titles []Title `json:"titles"`
}

func (law *Law) AddTitle(title Title) []Title {
	law.Titles = append(law.Titles, title)
	return law.Titles
}
