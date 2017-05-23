package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"

	"bitbucket.org/reneval/lawparser/domain"
)

type Publication struct {
}

type Index struct {
	Text  string
	Level int
}

type Indexes []Index

func (slice Indexes) Len() int {
	return len(slice)
}

func (slice Indexes) Less(i, j int) bool {
	return slice[i].Text < slice[j].Text
}

func (slice Indexes) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

//ParsePub parses a Publication Book
func (t *Publication) ParsePub(uri string) (*domain.Publication, error) {
	wg := new(sync.WaitGroup)

	lines := Stream(uri)
	titles, err := findTitles(wg, lines)
	if err != nil {
		return nil, err
	}
	stack := <-titles
	document := formPublication(stack)

	wg.Wait()
	return document, nil
}

func findTitles(wg *sync.WaitGroup, lines <-chan string) (<-chan *Stack, error) {
	wg.Add(1)
	var titles Indexes
	end := false
	listmap := NewStack(10)
	out := make(chan *Stack)
	var para []string

	fTitle, err := regexp.Compile("^[\\ ]*[a-zA-Z\u00C0-\u017F\\.-]+[[a-zA-Z\u00C0-\u017F\\ -]+[\\.]{2,}")
	indexEnd, err := regexp.Compile("^I\\.[a-zA-Z\u00C0-\u017F\\ ]+$")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	go func() {
		for text := range lines {
			if !end {
				found := fTitle.MatchString(text)
				if found {
					match := strings.Split(text, "..")
					titles = append(titles, Index{Text: strings.TrimSpace(match[0]), Level: findIndexLevel(match[0])})
				}
				end = indexEnd.MatchString(text)
				if end {
					fmt.Println("End of Index reached at:", text)
					fmt.Println("Index contains: ")
					fmt.Printf("%q", titles)
					sort.Sort(titles)
				}
			}

			//finding titles
			if end { //keep processing tesauro
				//Find titles
				row := strings.TrimSpace(text)
				i := sort.Search(len(titles), func(i int) bool { return titles[i].Text >= row })
				if i < len(titles) && titles[i].Text == row {
					listmap.Push(domain.PubTitle{ID: 0, Name: text, Level: titles[i].Level})
					//end of paragraph if a new title is found.
					if len(para) > 0 {
						paragraph := domain.PubParagraph{
							ID:   0,
							Text: strings.Join(para, " "),
						}
						listmap.Push(paragraph)
						para = nil
					}
				} else {
					if text == "" {
						paragraph := domain.PubParagraph{
							ID:   0,
							Text: strings.Join(para, " "),
						}
						listmap.Push(paragraph)
						para = nil
					} else {
						para = append(para, text)
					}
				}
			}
		}
		out <- listmap
		wg.Done()
	}()
	return out, nil
}

// Stream opens file and send it line by line
func Stream(uri string) <-chan string {
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

func formPublication(stack *Stack) *domain.Publication {
	tesauro := &domain.Publication{}
	var currentChild *domain.PubTitle
	lasWasTitle := false
	// isMainT, err := regexp.Compile("^[IVXL]+\\.[a-zA-Z\u00C0-\u017F\\ ]+$")
	// if err != nil {
	// 	return nil
	// }
	var currents []*domain.PubTitle

	for _, element := range stack.data {
		switch element := element.(type) {
		case domain.PubTitle:
			next := element.Level
			fmt.Println("Current level is: ", next, " and text is: ", element.Name)
			if len(currents) == 0 {
				//first time
				currentChild = tesauro.AddTitle(element)
				currents = append(currents, currentChild)
			} else if lasWasTitle {
				//title - subtitles
				if next == len(currents) {
					currentChild = currents[next-1].AddChild(element)
					currents = append(currents, currentChild)
				} else {
					currentChild = currents[next-1].AddChild(element)
					currents[next] = currentChild

				}
			} else {
				//after paragraph
				if next == 0 {
					currentChild = tesauro.AddTitle(element)
					currents[0] = currentChild
				} else {
					currentChild = currents[next-1].AddChild(element)
					currents[next] = currentChild

				}
			}

			lasWasTitle = true
		case domain.PubParagraph:
			fmt.Println("Paragraph found")
			if currentChild != nil {
				currentChild.AddParagraph(element)

			} else {
				fmt.Println("current child must always carry a value")
			}
			lasWasTitle = false
		}
	}
	return tesauro
}

func findIndexLevel(text string) int {
	i := 0
	for _, step := range text {
		if step == ' ' {
			i++
		} else {
			break
		}
	}
	return i
}
