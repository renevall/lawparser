package domain

type Publication struct {
	ID     int        `json:"id"`
	Titles []PubTitle `json:"titles"`
}

type PubTitle struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	Parent    int            `json:"parent"`
	Titles    []PubTitle     `json:"titles"`
	Paragraph []PubParagraph `json:"paragraph"`
	Level     int            `json:"level"`
}

type PubParagraph struct {
	ID      int    `json:"id"`
	TitleID int    `json:"title_id"`
	Text    string `json:"text"`
}

func (p *Publication) AddTitle(title PubTitle) *PubTitle {
	p.Titles = append(p.Titles, title)
	return &p.Titles[len(p.Titles)-1]
}

func (pt *PubTitle) AddChild(c PubTitle) *PubTitle {
	pt.Titles = append(pt.Titles, c)
	return &pt.Titles[len(pt.Titles)-1]
}

func (pt *PubTitle) AddParagraph(p PubParagraph) []PubParagraph {
	pt.Paragraph = append(pt.Paragraph, p)
	return pt.Paragraph
}
