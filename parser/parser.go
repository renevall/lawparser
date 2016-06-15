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
	fmt.Println(ts)

	return lines

}

//FindCTags finds the key words looking in the file
func FindCTags(in <-chan string, wg *sync.WaitGroup) (chan map[int]string, chan []int) {
	wg.Add(1)

	chm := make(chan map[int]string)
	chkeys := make(chan []int)

	go func(chan map[int]string, chan []int) {

		var keys []int
		t := make(map[int]string)

		i := 0

		for text := range in {

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

//FindBasicData process the data before first article
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

//ParseConcurrent parses a law using goroutines
func ParseConcurrent(uri string) *models.Law {
	wg := new(sync.WaitGroup)

	done := make(chan struct{})
	defer close(done)

	fmt.Println("Parse concurrent")
	in := StreamLines(uri)
	chs := fanOut(done, in, wg)

	law := FindBasicData(done, chs[0], wg)
	FindCTags(chs[1], wg)
	lines := prepareData(uri, wg)

	wg.Wait()
	fmt.Println(len(lines))

	return law

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
	// fmt.Println(data)
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
	Tag{"Arto", "Art\\.\\s\\d+|Artículo\\s\\d+"},
}
