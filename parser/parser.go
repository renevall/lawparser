package parser

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/goodsign/monday"

	"bitbucket.org/reneval/lawparser/models"
)

type foundTag struct {
	tagname string
	line    int
}

type preparedTag struct {
	name string
	exp  *regexp.Regexp
}

//stream lines to input channel
func StreamLines(uri string) <-chan string {
	fmt.Println("Streamlines", uri)

	out := make(chan string)

	go func() {
		f, err := os.Open(uri)
		if err != nil {
			fmt.Println("error opening file ", err)
			os.Exit(1)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			out <- scanner.Text()

		}
		close(out)
	}()
	return out
}

//prepareData loads file to memory from the streamer to form the law
func prepareData(uri string, wg *sync.WaitGroup) []string {

	var lines []string

	t := time.Now()
	file, err := ioutil.ReadFile(uri)
	if err != nil {
		log.Fatal(err)
	}
	lines = strings.Split(string(file), "\n")

	ts := time.Since(t)
	fmt.Println("Prepare data: ", ts)

	return lines

}

//FindCTags finds the key words looking in the file
func FindCTags(in <-chan string, wg *sync.WaitGroup) (chan foundTag, chan []int) {
	fmt.Println("FindCTags reached")
	wg.Add(1)
	var keys []int

	chm := make(chan foundTag)
	chi := make(chan []int)

	go func(chan foundTag, chan []int) {
		ti := time.Now()

		i := 0
		pTags := prepareTags(tags)
		var ft []foundTag

		for text := range in {
			for _, reg := range pTags {
				if reg.exp.MatchString(text) {
					//t[i+1] = tag.name
					f := foundTag{tagname: reg.name, line: i + 1}
					ft = append(ft, f)
					keys = append(keys, i+1)

					break
				}
			}
			i++
		}

		chi <- keys

		for _, v := range ft {
			chm <- v
		}
		sort.Ints(keys)

		ts := time.Since(ti)
		fmt.Println("FindCTags Time:", ts)
		close(chm)
		wg.Done()

	}(chm, chi)

	// fmt.Println(chkeys)
	return chm, chi

}

func prepareTags(tags []Tag) []preparedTag {
	t := time.Now()
	var ptags []preparedTag
	for _, v := range tags {
		r, _ := regexp.Compile(v.regex)
		ptag := preparedTag{v.name, r}
		ptags = append(ptags, ptag)

	}
	ts := time.Since(t)
	fmt.Println("prepareTags Time:", ts)
	return ptags
}

//FindBasicData process the data before first article
func FindBasicData(done chan<- struct{}, in <-chan string, wg *sync.WaitGroup) *models.Law {
	fmt.Println("find basic data")
	law := new(models.Law)
	t := time.Now()
	wg.Add(1)
	go func(*models.Law) {
		// defer wg.Done()
		for text := range in {
			for _, tag := range intro {
				r, _ := regexp.Compile(tag.regex)
				matches := make(map[int]*regexp.Regexp)
				matches[0], _ = regexp.Compile("(\\d{1,2}\\s(de|del))\\s\\w+\\s+(del|de)\\s\\d+")
				matches[1], _ = regexp.Compile("\\sdel|\\sde")
				if r.MatchString(text) {
					switch tag.name {
					case "Name":
						fillBasicData(tag.name, text, law, matches)
						break

					case "Number":
						fillBasicData(tag.name, text, law, matches)
						break

					case "Aproved":
						fillBasicData(tag.name, text, law, matches)
						break

					case "Diary":
						fillBasicData(tag.name, text, law, matches)
						break

					case "Arto":
						fmt.Println("End of Law Header Reached")
						done <- struct{}{}
						wg.Done()
						ts := time.Since(t)
						fmt.Println("FindBasicData", ts)
						return
					}
				}
			}
		}
	}(law)

	return law
}

func fillBasicData(tag string, value string, law *models.Law, matches map[int]*regexp.Regexp) {

	switch tag {
	case "Name":
		law.Name = value
		break
	case "Number":
		law.Name = value
		break
	case "Aproved":
		// TODO: parse date, using now for db test
		location, _ := time.LoadLocation("")
		data := matches[0].FindString(value)
		fmt.Println("Date match found", data)

		data = matches[1].ReplaceAllString(data, "")
		//law.ApprovalDate =
		date, _ := monday.ParseInLocation("02 January 2006", data, location, monday.LocaleEsES)
		law.ApprovalDate = date
		break

	case "Diary":
		law.Journal = value
		break
	}

}

//ParseConcurrent parses a law using goroutines
func ParseConcurrent(uri string) *models.Law {
	wg := new(sync.WaitGroup)

	done := make(chan struct{})
	defer close(done)

	fmt.Println("Parse concurrent")
	in := StreamLines(uri)
	chs := fanOut(done, in, wg)

	law := FindBasicData(done, chs[0], wg)
	tags, chi := FindCTags(chs[1], wg)
	lines := prepareData(uri, wg)
	lawLIFO := makeLaw(lines, law, tags, chi, wg)
	jsonFormat2(lawLIFO, law)

	wg.Wait()
	fmt.Println(len(lines))

	return law

}

func makeLaw(lines []string, law *models.Law, tag <-chan foundTag,
	index <-chan []int, wg *sync.WaitGroup) *Stack {

	mStack := NewStack(3)
	tags := <-index
	last := tags[len(tags)-1]
	r := 0

	for t := range tag {
		switch t.tagname {
		case "Titulo":
			// fmt.Println("Titulo ", lines[t.line])
			mStack.Push(models.Title{Name: lines[t.line], Reviewed: false})

		case "Capitulo":
			// fmt.Println("Capitulo ", lines[t.line])
			mStack.Push(models.Chapter{Name: lines[t.line], Reviewed: false})

		case "Arto":
			article := feedArticle(lines, last, t, tags, r)
			mStack.Push(article)

		}
		r++
	}

	return mStack
}

func feedArticle(lines []string, last int, t foundTag, tags []int, r int) models.Article {
	article_txt := []string{}
	mArticle := models.Article{Name: lines[t.line-1], Reviewed: false}

	if t.line != last {
		for x := t.line; x < tags[r+1]-1; x += 1 {
			article_txt = append(article_txt, lines[x])
		}
	} else {
		for x := t.line; x <= len(lines)-1; x += 1 {
			article_txt = append(article_txt, lines[x])
		}
	}

	mArticle.Text = strings.Join(article_txt, " ")

	return mArticle

}

func jsonFormat2(stack *Stack, mLaw *models.Law) *models.Law {

	currentTitle, currentChapter := -1, -1
	hasChapter, hasTitle := false, false

	for _, element := range stack.data {
		switch element := element.(type) {
		case models.Title:
			fmt.Println("Title", element.Name)
			mLaw.AddTitle(element)
			currentTitle++
			currentChapter = -1
			hasTitle = true

		case models.Chapter:
			fmt.Println("Chapter", element.Name)

			if len(mLaw.Titles) > 0 {
				fmt.Println("Adding Chapter under Title: ", currentTitle)
				mLaw.Titles[currentTitle].Chapters = mLaw.Titles[currentTitle].AddChapter(element)
			} else {
				fmt.Println("Adding Chapter under Law: ")
				mLaw.AddChapter(element)
			}
			currentChapter++
			hasChapter = true
		case models.Article:
			fmt.Println("Article", element.Name)

			//if law has titles
			if hasTitle {
				//case it has title but no chapter
				if currentChapter == -1 {
					element := models.Chapter{Name: "No Title"}
					mLaw.Titles[currentTitle].Chapters = mLaw.Titles[currentTitle].AddChapter(element)
					currentChapter++
				}
				//if it has titles and Chapters
				if hasChapter {
					fmt.Println("Adding Article under Title: ", currentTitle,
						" and Chapter: ", currentChapter)

					mLaw.Titles[currentTitle].Chapters[currentChapter].Articles =
						mLaw.Titles[currentTitle].Chapters[currentChapter].AddArticle(element)
				}
				//if Law doesnt have titles but does chapters
			} else if hasChapter {
				mLaw.Chapters[currentChapter].Articles =
					mLaw.Chapters[currentChapter].AddArticle(element)
				//if law only has articles
			} else {
				mLaw.AddArticle(element)
			}
		}
	}
	return mLaw
}

//fanOut distributes the readed files lines to different channels to process in parallel
func fanOut(done <-chan struct{}, ch <-chan string, wg *sync.WaitGroup) []chan string {

	cs := []chan string{
		make(chan string),
		make(chan string),
	}

	go func() {
		stoped := false
		for mainch := range ch {
			if !stoped {
				stoped = broadcastCancel(done, cs[0], mainch)

			}
			broadcast(cs[1], mainch)

		}

		for _, c := range cs {
			close(c)
		}
	}()
	return cs
}

func broadcast(ch chan<- string, data string) {
	ch <- data
}

func broadcastCancel(done <-chan struct{}, ch chan<- string, data string) bool {
	select {
	case ch <- data:
		return false
	case <-done:
		return true
	}

	return false
}

var intro = Tags{
	Tag{"Name", "LEY DE|CÓDIGO"},
	Tag{"Number", "No\\."},
	Tag{"Aproved", "Aprobada"},
	Tag{"Diary", "Publicada"},
	Tag{"Arto", "Art\\.\\s\\d+"},
}
