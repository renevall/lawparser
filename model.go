package main

type Title struct {
	name     string
	chapters []Chapter
}

type Chapter struct {
	name     string
	articles []Article
}

type Article struct {
	name string
	text string
}

type Titles []Title
