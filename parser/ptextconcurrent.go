package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sync"

	"bitbucket.org/reneval/lawparser/models"
)

type FoundTag struct {
	tagname string
	line    int
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

func FindCTags(in <-chan string, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for text := range in {
			for _, tag := range tags {
				r, _ := regexp.Compile(tag.regex)
				if r.MatchString(text) {
				}
			}
		}
	}()

}

func FindBasicData(in <-chan string, wg *sync.WaitGroup) *models.Law {
	fmt.Println("find basic data")
	law := new(models.Law)
	wg.Add(1)
	go func(*models.Law) {
		defer wg.Done()
		for text := range in {
			for _, tag := range intro {
				r, _ := regexp.Compile(tag.regex)
				if r.MatchString(text) {
					fmt.Println("match!")
					fmt.Println(tag.name)
					switch tag.name {
					case "Name":
						fmt.Println(text)
						fillBasicData(tag.name, text, law)
						break

					case "Number":
						fmt.Println(text)
						fillBasicData(tag.name, text, law)
						break

					case "Aproved":
						fmt.Println(text)
						fillBasicData(tag.name, text, law)
						break

					case "Diary":
						fmt.Println(text)
						fillBasicData(tag.name, text, law)
						break

					case "Arto":
						fmt.Println("End of Law Header Reached")
						wg.Done()
						return
					}
				}
			}
		}
	}(law)
	return law
}

func fillBasicData(tag string, value string, law *models.Law) {
	switch tag {
	case "Name":
		law.Name = value
		break
	case "Number":
		law.Name = value
		break
	case "Aproved":
		law.ApprovalDate = value
		break

	case "Diary":
		law.Journal = value
		break
	}

}

func ParseConcurrent(uri string) *models.Law {
	wg := new(sync.WaitGroup)

	fmt.Println("Parse concurrent")
	in := StreamLines(uri)
	chs := fanOut(in, wg)
	law := FindBasicData(chs[0], wg)
	FindCTags(chs[1], wg)

	wg.Wait()

	return law

}

func fanOut(ch <-chan string, wg *sync.WaitGroup) []chan string {

	cs := []chan string{
		make(chan string),
		make(chan string),
	}
	go func() {
		for i := range ch {
			for _, w := range cs {
				w <- i
			}
		}
		for _, c := range cs {
			close(c)
		}
	}()
	return cs
}

var intro = Tags{
	Tag{"Name", "LEY DE|CÓDIGO"},
	Tag{"Number", "No\\."},
	Tag{"Aproved", "Aprobada"},
	Tag{"Diary", "Publicada"},
	Tag{"Arto", "Art\\.\\s\\d+|Artículo\\s\\d+"},
}
