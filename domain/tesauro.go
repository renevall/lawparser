package domain

type Tesauro struct {
	ID     int
	Titles []TTitle
}

type TTitle struct {
	ID        int
	Name      string
	Parent    int
	Titles    []TTitle
	Paragraph []TParagraph
}

type TParagraph struct {
	ID     int
	Parent int
	Text   string
}

func (t *Tesauro) AddTitle(title TTitle) *TTitle {
	t.Titles = append(t.Titles, title)
	return &t.Titles[len(t.Titles)-1]
}

func (tt *TTitle) AddChild(c TTitle) *TTitle {
	tt.Titles = append(tt.Titles, c)
	return &tt.Titles[len(tt.Titles)-1]
}

func (tt *TTitle) AddParagraph(p TParagraph) []TParagraph {
	tt.Paragraph = append(tt.Paragraph, p)
	return tt.Paragraph
}
