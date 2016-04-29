package models

//Chapter is the model for a Law chapter
type Chapter struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Articles []Article `json:"articles"`
	TitleID  int       `json:"titleID"`
}

//AddArticle adds parsed article data to parsed law object
func (chapter *Chapter) AddArticle(article Article) []Article {
	chapter.Articles = append(chapter.Articles, article)
	return chapter.Articles
}
