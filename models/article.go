package models

//Article Holds the article model and his methods
type Article struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Text      string `json:"text"`
	ArticleID int    `json:"articleID"`
}
