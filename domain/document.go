package domain

type Document struct {
	ID     int
	Titles []DocTitle
}

type DocTitle struct {
	ID        int
	Name      string
	Parent    int
	Titles    []DocTitle
	Paragraph []DocParagraph
	Level     int
}

type DocParagraph struct {
	ID     int
	Parent int
	Text   string
}

func (d *Document) AddTitle(title DocTitle) *DocTitle {
	d.Titles = append(d.Titles, title)
	return &d.Titles[len(d.Titles)-1]
}

func (dt *DocTitle) AddChild(c DocTitle) *DocTitle {
	dt.Titles = append(dt.Titles, c)
	return &dt.Titles[len(dt.Titles)-1]
}

func (dt *DocTitle) AddParagraph(p DocParagraph) []DocParagraph {
	dt.Paragraph = append(dt.Paragraph, p)
	return dt.Paragraph
}
