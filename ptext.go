package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"sort"
	"strings"
)

type Tag struct {
	name  string
	regex string
}

type Tags []Tag

var title int = 0

var titles Titles

func OpenTextFile(uri string) []string {
	file, err := ioutil.ReadFile(uri)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(file), "\n")
	return lines
}

func ParseText(uri string) {
	data := OpenTextFile(uri)
	ref, order := FindTags(data)
	MakeLaw(data, order, ref)
}

func FindTags(data []string) (map[int]string, []int) {
	var keys []int
	m := make(map[int]string)
	for i, text := range data {
		for _, tag := range tags {
			r, _ := regexp.Compile(tag.regex)
			if r.MatchString(text) {
				m[i+1] = tag.name
				break
			}
		}
	}

	for k, v := range m {
		keys = append(keys, k)
		if v == "Titulo" {
			title = title + 1
		}
	}

	sort.Ints(keys)

	return m, keys

}

func MakeLaw(data []string, index []int, ref map[int]string) {
	var title_index, chapter_index int
	title_index, chapter_index = -1, -1
	// fmt.Println("title_index at first", title_index)
	// fmt.Println("chapter_index at first", chapter_index)
	for _, k := range index {
		if ref[k] == "Titulo" {
			title_index = title_index + 1
			chapter_index = -1
			titles = append(titles, Title{name: data[k]})
		}

		if ref[k] == "Capitulo" {
			chapter_index = chapter_index + 1

			titles[title_index].chapters =
				append(titles[title_index].chapters, Chapter{name: data[k]})

		}

		if ref[k] == "Arto" {
			// fmt.Println("title_index =", title_index)
			// fmt.Println("chapter_index =", chapter_index, " ", data[k-1])
			fmt.Println("procesando linea: ", k)
			titles[title_index].chapters[chapter_index].articles =
				append(titles[title_index].chapters[chapter_index].articles, Article{name: data[k-1]})

		}

	}

	fmt.Printf("%+v", titles)
}

var tags = Tags{
	Tag{"Titulo", "TÍTULO"},
	Tag{"Capitulo", "Capí?tulo\\s|Capitulo\\s"},
	Tag{"Artículo", "Artículo\\s\\d+"},
	Tag{"Arto", "Art.\\s\\d+"},
}
