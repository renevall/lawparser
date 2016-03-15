package parser

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

func OpenTextFile(uri string) []string {
	file, err := ioutil.ReadFile(uri)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(file), "\n")
	return lines
}

func ParseText(uri string) *Law {
	data := OpenTextFile(uri)
	ref, order := FindTags(data)
	parsed_law := MakeLaw2(data, order, ref)
	json_law := jsonFormat(parsed_law)
	return json_law
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

	// fmt.Println("Ended Find Tags")
	// fmt.Println(keys)
	// fmt.Println("--------------------------------------------------------------")
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
			// fmt.Println("Chapter index: ", chapter_index)
		}

		if ref[k] == "Arto" {

			article_index = article_index + 1

			fmt.Println("procesando linea: ", k)
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

func MakeLaw2(data []string, index []int, ref map[int]string) *Stack {
	var article_index int
	article_txt := []string{}

	article_index = -1
	last := index[len(index)-1]

	mStack := NewStack(3)

	for r, k := range index {
		if ref[k] == "Titulo" {
			mStack.Push(Title{Name: data[k]})
		}

		if ref[k] == "Capitulo" {
			mStack.Push(Chapter{Name: data[k]})
		}

		if ref[k] == "Arto" {
			mArticle := Article{Name: data[k-1]}
			article_index = article_index + 1
			if k != last {
				for x := k; x < index[r+1]-1; x += 1 {
					article_txt = append(article_txt, data[x])
				}
			} else {
				for x := k; x <= len(data)-1; x += 1 {
					article_txt = append(article_txt, data[x])
				}
			}

			mArticle.Text = strings.Join(article_txt, " ")
			mStack.Push(mArticle)
			article_txt = nil
		}
	}

	//fmt.Printf("%+v", titles)

	// fmt.Println(index)
	return mStack
}

func jsonFormat(stack *Stack) *Law {
	mLaw := new(Law)
	current_title, current_chapter := -1, -1

	for _, element := range stack.data {
		switch element := element.(type) {
		case Title:
			mLaw.addTitle(element)
			current_title += 1
			current_chapter = -1

		case Chapter:
			mLaw.Titles[current_title].Chapters = mLaw.Titles[current_title].addChapter(element)
			current_chapter += 1

		case Article:

			if len(mLaw.Titles[current_title].Chapters) > 0 {
				mLaw.Titles[current_title].Chapters[current_chapter].Articles =
					mLaw.Titles[current_title].Chapters[current_chapter].addArticle(element)
			}
		}
	}
	return mLaw
}

var tags = Tags{
	Tag{"Titulo", "TÍTULO"},
	Tag{"Capitulo", "Capí?tulo\\s|Capitulo\\s"},
	Tag{"Artículo", "Artículo\\s\\d+"},
	Tag{"Arto", "Art.\\s\\d+"},
}
