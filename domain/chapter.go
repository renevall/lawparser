package domain

//Chapter is the model for a Law chapter
type Chapter struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Articles []Article `json:"articles"`
	TitleID  int64     `json:"titleID"`
	LawID    int64     `json:"lawID"`
	Reviewed bool      `json:"reviewed"`
}

//AddArticle adds parsed article data to parsed law object
func (c *Chapter) AddArticle(article Article) []Article {
	c.Articles = append(c.Articles, article)
	return c.Articles
}
