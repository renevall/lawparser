package parser

type Law struct {
	Name   string  `json:name`
	Titles []Title `json:titles`
}

type Title struct {
	Name     string    `json:"name"`
	Chapters []Chapter `json: "chapters"`
}

type Chapter struct {
	Name     string    `json:"name"`
	Articles []Article `json:"articles"`
}

type Article struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type Titles []Title

func (law *Law) addTitle(title Title) []Title {
	law.Titles = append(law.Titles, title)
	return law.Titles
}

func (title *Title) addChapter(chapter Chapter) []Chapter {
	title.Chapters = append(title.Chapters, chapter)
	return title.Chapters
}

func (chapter *Chapter) addArticle(article Article) []Article {
	chapter.Articles = append(chapter.Articles, article)
	return chapter.Articles
}
