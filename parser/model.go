package parser

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
