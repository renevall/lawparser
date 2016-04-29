package models

import "time"

//Law struct with most methods.
type Law struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Titles       []Title   `json:"titles"`
	ApprovalDate time.Time `json:"approvalDate"`
	PublishDate  time.Time `json:"publishDate"`
	Journal      string    `json:"journal"`
	Intro        string    `json:"intro"`
}

//AddTitle adds parsed title data to parsed law object
func (law *Law) AddTitle(title Title) []Title {
	law.Titles = append(law.Titles, title)
	return law.Titles
}
