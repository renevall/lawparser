package domain

type Publication struct {
	ID     int
	Titles []PubTitle
}

type PubTitle struct {
	ID        int
	Name      string
	Parent    int
	Titles    []PubTitle
	Paragraph []PubParagraph
	Level     int
}

type PubParagraph struct {
	ID     int
	Parent int
	Text   string
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
