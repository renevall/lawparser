package parser

import (
	// "fmt"
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

func OpenTextFile(uri string) []string {
	file, err := ioutil.ReadFile(uri)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(file), "\n")
	return lines
}

func ParseText(uri string) Titles {
	data := OpenTextFile(uri)
	ref, order := FindTags(data)
	parsed_law := MakeLaw(data, order, ref)
	return parsed_law
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

func MakeLaw(data []string, index []int, ref map[int]string) Titles {
	var title_index, chapter_index, article_index int
	article_txt := []string{}
	var titles = Titles{}

	title_index, chapter_index, article_index = -1, -1, -1
	last := index[len(index)-1]

	for r, k := range index {
		if ref[k] == "Titulo" {
			title_index = title_index + 1
			chapter_index = -1
			article_index = -1
			titles = append(titles, Title{Name: data[k]})
		}

		if ref[k] == "Capitulo" {
			chapter_index = chapter_index + 1
			article_index = -1

			titles[title_index].Chapters =
				append(titles[title_index].Chapters, Chapter{Name: data[k]})
			//fmt.Println("Chapter index: ", chapter_index)
		}

		if ref[k] == "Arto" {

			article_index = article_index + 1

			// fmt.Println("procesando linea: ", k)
			titles[title_index].Chapters[chapter_index].Articles =
				append(titles[title_index].Chapters[chapter_index].Articles, Article{Name: data[k-1]})
			if k != last {
				for x := k; x < index[r+1]-1; x += 1 {
					article_txt = append(article_txt, data[x])
				}
			} else {
				for x := k; x <= len(data)-1; x += 1 {
					article_txt = append(article_txt, data[x])
				}
			}

			titles[title_index].Chapters[chapter_index].Articles[article_index].Text =
				strings.Join(article_txt, " ")
			article_txt = nil

		}

	}

	//fmt.Printf("%+v", titles)

	// fmt.Println(index)
	return titles
}

var tags = Tags{
	Tag{"Titulo", "TÍTULO"},
	Tag{"Capitulo", "Capí?tulo\\s|Capitulo\\s"},
	Tag{"Artículo", "Artículo\\s\\d+"},
	Tag{"Arto", "Art.\\s\\d+"},
}
