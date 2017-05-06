package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"bitbucket.org/reneval/lawparser/domain"
)

type Tesauro struct {
}

//Parse parses a Tesauro Book
func (t *Tesauro) Parse(uri string) error {
	lines := Stream("./testlaws/index.txt")
	findTitles(lines)
	return nil
}

func findTitles(lines <-chan string) ([]string, error) {

	var titles []string
	end := false
	listmap := NewStack(10)
	var para []string

	fTitle, err := regexp.Compile("^[a-zA-Z\u00C0-\u017F\\s]+[\\.]{2,}")
	indexEnd, err := regexp.Compile("^I\\.[a-zA-Z\u00C0-\u017F\\ ]+$")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	go func() {
		for text := range lines {

			//finding titles
			if !end {
				end = indexEnd.MatchString(text)
				if end {
					fmt.Println("End of Index reached at:", text)
					sort.Strings(titles)
				}
				found := fTitle.MatchString(text)
				if found {
					match := strings.Split(text, "..")
					titles = append(titles, strings.TrimSpace((match[0])))
				}
			} else { //keep processing tesauro
				//Find titles
				row := strings.TrimSpace(text)
				i := sort.SearchStrings(titles, row)
				if i < len(titles) && titles[i] == row {
					// fmt.Println("Found Title: ", text)
					listmap.Push(&domain.TTitle{ID: 0, Title: text})
					if len(para) > 0 {
						paragraph := &domain.TParagraph{
							ID:   0,
							Text: strings.Join(para, " "),
						}
						listmap.Push(paragraph)
						para = nil
					}
				} else {
					para = append(para, text)
				}

			}

		}
		fmt.Printf("%q", listmap)

	}()
	return titles, nil
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
