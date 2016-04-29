package models

//Chapter is the model for a Law chapter
type Chapter struct {
	Name     string    `json:"name"`
	Articles []Article `json:"articles"`
}

func (chapter *Chapter) AddArticle(article Article) []Article {
	chapter.Articles = append(chapter.Articles, article)
	return chapter.Articles
}
