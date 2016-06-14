package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
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

func FindCTags(in <-chan string, wg *sync.WaitGroup) (chan map[int]string, chan []int) {
	wg.Add(1)

	chm := make(chan map[int]string)
	chkeys := make(chan []int)

	go func(chan map[int]string, chan []int) {

		var keys []int
		t := make(map[int]string)

		i := 0

		for text := range in {
			fmt.Println(text)
			for _, tag := range tags {
				r, _ := regexp.Compile(tag.regex)
				if r.MatchString(text) {
					t[i+1] = tag.name
					break
				}
			}
			i++
		}

		chm <- t

		for k, v := range t {
			keys = append(keys, k)
			if v == "Titulo" {
				title = title + 1
			}
		}

		sort.Ints(keys)
		chkeys <- keys

		close(chm)
		close(chkeys)
		wg.Done()

	}(chm, chkeys)

	return chm, chkeys

}

func FindBasicData(done chan<- struct{}, in <-chan string, wg *sync.WaitGroup) *models.Law {
	fmt.Println("find basic data")
	law := new(models.Law)

	wg.Add(1)
	go func(*models.Law) {
		defer wg.Done()
		for text := range in {
			for _, tag := range intro {
				r, _ := regexp.Compile(tag.regex)
				if r.MatchString(text) {
					switch tag.name {
					case "Name":
						fillBasicData(tag.name, text, law)
						break

					case "Number":
						fillBasicData(tag.name, text, law)
						break

					case "Aproved":
						fillBasicData(tag.name, text, law)
						break

					case "Diary":
						fillBasicData(tag.name, text, law)
						break

					case "Arto":
						fmt.Println("End of Law Header Reached")
						done <- struct{}{}
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

	done := make(chan struct{})

	fmt.Println("Parse concurrent")
	in := StreamLines(uri)
	chs := fanOut(done, in, wg)

	law := FindBasicData(done, chs[0], wg)
	tags, _ := FindCTags(chs[1], wg)

	fmt.Println("Tags: ", <-tags)

	wg.Wait()

	return law

}

func fanOut(done <-chan struct{}, ch <-chan string, wg *sync.WaitGroup) []chan string {

	cs := []chan string{
		make(chan string),
		make(chan string),
	}

	var closed bool = false
	go func() {
		for i := range ch {
			for r, w := range cs {
				select {
				case <-done:
					closed = true
					break

				default:
					if r == 0 && !closed {
						w <- i
					} else if r != 0 {
						w <- i
					}
				}
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
