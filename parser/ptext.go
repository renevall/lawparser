package parser

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"sort"
	"strings"

	"bitbucket.org/reneval/lawparser/models"
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

func ParseText(uri string) *models.Law {
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

func MakeLaw(data []string, index []int, ref map[int]string) []models.Title {
	var title_index, chapter_index, article_index int
	article_txt := []string{}
	var titles = []models.Title{}

	title_index, chapter_index, article_index = -1, -1, -1
	last := index[len(index)-1]

	for r, k := range index {
		if ref[k] == "Titulo" {
			title_index = title_index + 1
			chapter_index = -1
			article_index = -1
			titles = append(titles, models.Title{Name: data[k]})
		}

		if ref[k] == "Capitulo" {
			chapter_index = chapter_index + 1
			article_index = -1

			titles[title_index].Chapters =
				append(titles[title_index].Chapters, models.Chapter{Name: data[k]})
			// fmt.Println("Chapter index: ", chapter_index)
		}

		if ref[k] == "Arto" {

			article_index = article_index + 1

			fmt.Println("procesando linea: ", k)
			titles[title_index].Chapters[chapter_index].Articles =
				append(titles[title_index].Chapters[chapter_index].Articles, models.Article{Name: data[k-1]})
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
			mStack.Push(models.Title{Name: data[k]})
		}

		if ref[k] == "Capitulo" {
			mStack.Push(models.Chapter{Name: data[k]})
		}

		if ref[k] == "Arto" {
			mArticle := models.Article{Name: data[k-1]}
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

func jsonFormat(stack *Stack) *models.Law {
	mLaw := new(models.Law)
	current_title, current_chapter := -1, -1

	for _, element := range stack.data {
		switch element := element.(type) {
		case models.Title:
			mLaw.AddTitle(element)
			current_title += 1
			current_chapter = -1

		case models.Chapter:
			mLaw.Titles[current_title].Chapters = mLaw.Titles[current_title].AddChapter(element)
			current_chapter += 1

		case models.Article:

			if len(mLaw.Titles[current_title].Chapters) > 0 {
				mLaw.Titles[current_title].Chapters[current_chapter].Articles =
					mLaw.Titles[current_title].Chapters[current_chapter].AddArticle(element)
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

//InsertLawToDB inserts all parsed law to DB
func InsertLawToDB(db *sql.DB, law *models.Law) error {
	//TODO first insert, get id of inserted law
	for _, title := range law.Titles {
		//TODO second insert
		for _, chapter := range title.Chapters {
			//TODO insert of Chapter, get ID
			tx, err := db.Begin()
			if err != nil {
				log.Fatal(err)
			}
			for _, article := range chapter.Articles {

				//TODO INSERT ARTICLES, DONT FORGET TO SET FK
				article.ChapterID = 1
				err := article.CreateArticle(db, tx)
				if err != nil {
					log.Println(err)
					return nil
				}
			}
			tx.Commit()

		}
	}
	return nil
}
