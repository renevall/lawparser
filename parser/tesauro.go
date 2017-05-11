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

type Tesauro struct {
}

//Parse parses a Tesauro Book
func (t *Tesauro) Parse(uri string) (*domain.Tesauro, error) {
	wg := new(sync.WaitGroup)

	lines := Stream("./testlaws/index.txt")
	titles, err := findTitles(wg, lines)
	if err != nil {
		return nil, err
	}
	stack := <-titles
	tesauro := formTesauro(stack)

	wg.Wait()
	return tesauro, nil
}

func findTitles(wg *sync.WaitGroup, lines <-chan string) (<-chan *Stack, error) {
	wg.Add(1)
	var titles []string
	end := false
	listmap := NewStack(10)
	out := make(chan *Stack)
	var para []string

	fTitle, err := regexp.Compile("^[a-zA-Z\u00C0-\u017F\\.]+[[a-zA-Z\u00C0-\u017F\\ ]+[\\.]{2,}")
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
					titles = append(titles, strings.TrimSpace((match[0])))
				}
				end = indexEnd.MatchString(text)
				if end {
					fmt.Println("End of Index reached at:", text)
					fmt.Println("Index contains: ")
					fmt.Printf("%q", titles)
					sort.Strings(titles)

				}
			}

			//finding titles

			if end { //keep processing tesauro
				//Find titles
				row := strings.TrimSpace(text)
				i := sort.SearchStrings(titles, row)
				if i < len(titles) && titles[i] == row {
					listmap.Push(domain.TTitle{ID: 0, Name: text})
					//end of paragraph if a new title is found.
					if len(para) > 0 {
						paragraph := domain.TParagraph{
							ID:   0,
							Text: strings.Join(para, " "),
						}
						listmap.Push(paragraph)
						para = nil
					}
				} else {
					if text == "" {
						paragraph := domain.TParagraph{
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

func formTesauro(stack *Stack) *domain.Tesauro {
	tesauro := &domain.Tesauro{}
	var currentChild, currentMain *domain.TTitle
	lasWasTitle := false
	isMainT, err := regexp.Compile("^[IVXL]+\\.[a-zA-Z\u00C0-\u017F\\ ]+$")
	if err != nil {
		return nil
	}

	for _, element := range stack.data {
		switch element := element.(type) {
		case domain.TTitle:

			main := isMainT.MatchString(element.Name)
			fmt.Println("Title found")
			if currentChild == nil { //first time
				currentChild = tesauro.AddTitle(element)
				if main {
					currentMain = currentChild
				}
			} else if lasWasTitle { //subtitle
				currentChild = currentChild.AddChild(element)
			} else { // last one was not first time, and current isn't nil, (last was )
				if main {
					currentChild = tesauro.AddTitle(element)
					currentMain = currentChild
				} else {
					currentChild = currentMain.AddChild(element)
				}
			}

			lasWasTitle = true
		case domain.TParagraph:
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
