package parser

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/goodsign/monday"

	"bitbucket.org/reneval/lawparser/domain"
)

type foundTag struct {
	tagname string
	line    int
	text    string
}

type preparedTag struct {
	name string
	exp  *regexp.Regexp
}

//ParseConcurrent parses a law using goroutines
//TODO: Use domain, not model
func ParseConcurrent(uri string) *domain.Law {
	fmt.Println("Parse concurrent called")
	wg := new(sync.WaitGroup)

	done := make(chan struct{})
	defer close(done)

	in := StreamLines(uri)
	chs := fanOut(done, in, wg)
	law := NewLaw()

	FindBasicData(law, done, chs[0], wg)
	tags, chi := FindCTags(chs[1], wg)
	lines := prepareData(uri, wg)
	lawLIFO := makeLaw(lines, law, tags, chi, wg)
	jsonFormat2(lawLIFO, law)

	wg.Wait()
	// fmt.Println(lawLIFO)

	return law

}

func NewLaw() *domain.Law {
	return &domain.Law{}
}

// StreamLines opens file and send it line by line
func StreamLines(uri string) <-chan string {
	fmt.Println("Streaming Lines started")

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

//fanOut distributes the readed files lines to different channels to process in parallel
func fanOut(done <-chan struct{}, ch <-chan string, wg *sync.WaitGroup) []chan string {
	fmt.Println("Sending lines to different channels")
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

//FindBasicData process the data before first article
func FindBasicData(law *domain.Law, done chan<- struct{}, in <-chan string, wg *sync.WaitGroup) {
	// t := time.Now()
	var intro []string
	wg.Add(1)
	matches := make(map[int]*regexp.Regexp)
	matches[0], _ = regexp.Compile("(\\d{1,2}\\s(de|del))\\s\\w+\\s+(del|de)\\s\\d+")
	matches[1], _ = regexp.Compile("\\sdel|\\sde")
	matches[2], _ = regexp.Compile("\\.")
	go func(*domain.Law) {
		defer wg.Done()
		for text := range in {
			intro = append(intro, text)
			for _, tag := range introTags {
				r, _ := regexp.Compile(tag.regex)
				if r.MatchString(text) {
					switch tag.name {
					case "Name":
						fillBasicData(tag.name, text, law, matches)
						break

					case "Number":
						// fmt.Println("Law number found in:", text)
						n, err := regexp.Compile("[0-9]+")
						if err != nil {
							fmt.Println(err)
						}
						number := n.FindString(text)
						// fmt.Println("found: ", number)
						fillBasicData(tag.name, number, law, matches)
						break

					case "Aproved":
						a, err := regexp.Compile("[0-9]{1,2}\\s\\w+\\s\\w+\\s\\w+\\s[0-9]+")
						if err != nil {
							fmt.Println(err)
						}
						date := a.FindString(text)
						fillBasicData(tag.name, date, law, matches)
						break

					case "Diary":
						d, err := regexp.Compile("[0-9]+")
						if err != nil {
							fmt.Println(err)
						}
						journal := d.FindString(text)
						fillBasicData(tag.name, journal, law, matches)
						break

					case "Arto":
						law.Intro = strings.Join(intro, " ")
						done <- struct{}{}
						// ts := time.Since(t)
						// fmt.Println("FindBasicData", ts)
						return
					}
				}
			}
		}
	}(law)
}

// FindCTags look for keywords in the file
func FindCTags(in <-chan string, wg *sync.WaitGroup) (chan foundTag, chan []int) {
	// fmt.Println("FindCTags reached")
	wg.Add(1)
	var keys []int

	chm := make(chan foundTag)
	chi := make(chan []int)

	go func() {
		// ti := time.Now()

		i := 0
		pTags := prepareTags(tags)
		var ft []foundTag

		for text := range in {
			for _, reg := range pTags {
				if reg.exp.MatchString(text) {
					f := foundTag{tagname: reg.name, line: i + 1, text: text}
					ft = append(ft, f)
					keys = append(keys, i+1)
					break
				}
			}
			i++

		}
		chi <- keys

		// fmt.Println("Second stage reached")
		for _, v := range ft {
			chm <- v
		}

		// sort.Ints(keys)

		// ts := time.Since(ti)
		// fmt.Println("FindCTags Time:", ts)
		close(chm)
		wg.Done()

	}()

	return chm, chi

}

// prepareData loads file to memory from the streamer to form the law
func prepareData(uri string, wg *sync.WaitGroup) []string {

	var lines []string

	// t := time.Now()
	file, err := ioutil.ReadFile(uri)
	if err != nil {
		log.Fatal(err)
	}
	lines = strings.Split(string(file), "\n")

	// ts := time.Since(t)
	// fmt.Println("Prepare data: ", ts)

	return lines

}

func makeLaw(lines []string, law *domain.Law, tag <-chan foundTag,
	index <-chan []int, wg *sync.WaitGroup) *Stack {

	mStack := NewStack(3) //estimated stack size
	tags := <-index
	last := tags[len(tags)-1]
	r := 0

	for t := range tag {
		switch t.tagname {
		case "Libro":
			mStack.Push(domain.Book{Name: lines[t.line], Reviewed: false})
		case "Titulo":
			// fmt.Println("Titulo ", lines[t.line])
			mStack.Push(domain.Title{Name: lines[t.line], Reviewed: false})

		case "Capitulo":
			// fmt.Println("Capitulo ", lines[t.line])
			mStack.Push(domain.Chapter{Name: lines[t.line], Reviewed: false})

		case "Arto":
			article := feedArticle(lines, last, t, tags, r)
			mStack.Push(article)

		}
		r++
	}

	return mStack
}

func jsonFormat2(stack *Stack, mLaw *domain.Law) *domain.Law {

	currentBook, currentTitle, currentChapter := -1, -1, -1
	hasBook, hasChapter, hasTitle := false, false, false

	for _, element := range stack.data {
		switch element := element.(type) {

		case domain.Book:
			fmt.Printf("Book", element.Name)
			mLaw.AddBook(element)
			currentBook++
			currentChapter = -1
			currentTitle = -1
			hasBook = true
		case domain.Title:
			fmt.Println("Title", element.Name)
			currentTitle++
			currentChapter = -1
			hasTitle = true
			if len(mLaw.Books) > 0 {
				fmt.Println("Adding Title under Book: ", currentBook)
				mLaw.Books[currentBook].Titles = mLaw.Books[currentBook].AddTitle(element)
			} else {
				fmt.Println("Adding Title under Law")
				mLaw.AddTitle(element)
			}
		case domain.Chapter:
			fmt.Println("Chapter", element.Name)

			if len(mLaw.Books) > 0 {
				fmt.Println("Adding Chapter under Title: ", currentTitle, " from book: ", currentBook)
				mLaw.Books[currentBook].Titles[currentTitle].Chapters =
					mLaw.Books[currentBook].Titles[currentTitle].AddChapter(element)
			} else {
				if len(mLaw.Titles) > 0 {
					fmt.Println("Adding Chapter under Title: ", currentTitle)
					mLaw.Titles[currentTitle].Chapters = mLaw.Titles[currentTitle].AddChapter(element)
				} else {
					fmt.Println("Adding Chapter under Law: ")
					mLaw.AddChapter(element)
				}
			}

			currentChapter++
			hasChapter = true
		case domain.Article:
			fmt.Println("Article", element.Name)

			if hasBook {
				mLaw.Books[currentBook].Titles[currentTitle].Chapters[currentChapter].AddArticle(element)
			} else {
				if hasTitle {
					//case it has title but no chapter
					if currentChapter == -1 {
						element := domain.Chapter{Name: "No Title"}
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
			//if law has titles
		}
	}
	return mLaw
}

func prepareTags(tags []Tag) []preparedTag {
	// t := time.Now()
	var ptags []preparedTag
	for _, v := range tags {
		r, _ := regexp.Compile(v.regex)
		ptag := preparedTag{v.name, r}
		ptags = append(ptags, ptag)

	}
	// ts := time.Since(t)
	// fmt.Println("prepareTags Time:", ts)
	return ptags
}

func fillBasicData(tag string, value string, law *domain.Law, matches map[int]*regexp.Regexp) {

	switch tag {
	case "Name":
		law.Name = value
		break
	case "Number":
		law.Number, _ = strconv.Atoi(value)
		break
	case "Aproved":
		location, _ := time.LoadLocation("")
		data := matches[0].FindString(value)

		data = matches[1].ReplaceAllString(data, "")
		data = matches[2].ReplaceAllString(data, "")
		date, err := monday.ParseInLocation("2 January 2006", data, location, monday.LocaleEsES)
		if err != nil {
			fmt.Println(err)
		}
		law.ApprovalDate = date
		break

	case "Diary":
		law.Journal = value
		break
	}

}

func feedArticle(lines []string, last int, t foundTag, tags []int, r int) domain.Article {
	article_txt := []string{}
	mArticle := domain.Article{Name: lines[t.line-1], Reviewed: false}

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

var introTags = Tags{
	Tag{"Name", "^LEY DE|^C(Ó?|O?)DIGO"},
	Tag{"Number", "No\\.|N°."},
	Tag{"Aproved", "Aprobad(a|o)"},
	Tag{"Diary", "Publicad(a|o)"},
	Tag{"Arto", `^\f?(?:Art.\s\d+|Arto.\s\d+|Artículo\s\d+|Articulo\s\d+)`},
}

//TODO: Make tags consider words like "Único" and weird ass symbol: "\f"
var tags = Tags{
	Tag{"Titulo", "^\f?(TÍTULO\\s?([IVX\u00C0-\u00FF]|$)|^TITULO\\s?([IVX\u00C0-\u00FF]|$)|TITULO\\s\\w+$|TÍTULO\\s\\w+$)"},
	Tag{"Capitulo", "^\f?(?:Capítulo\\s[\u00C0-\u00FF]?\\w+$|Capí?tulo\\s?\\w{0,3}$|Capitulo\\s?\\w{0,3}$|CAP(Í?|I?)TULO\\s?)"},
	Tag{"Arto", `^\f?(?:Art.\s\d+|Arto.\s\d+|Artículo\s\d+|Articulo\s\d+)`},
	Tag{"Libro", "^\f?LIBRO\\s[IVXLCDM]+$"},
}
