package domain

type Tesauro struct {
	ID     int
	Titles []Title
}

type TTitle struct {
	ID     int
	Title  string
	Parent int
	Titles []Title
}

type TParagraph struct {
	ID    int
	Title int
	Text  string
}
