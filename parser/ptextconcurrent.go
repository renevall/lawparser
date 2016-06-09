package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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

func FindCTags(in <-chan string) {
	// out := make(chan FoundTag)
	l := 1

	go func() {
		for text := range in {
			for _, tag := range tags {
				r, _ := regexp.Compile(tag.regex)
				if r.MatchString(text) {
					// fmt.Println("Tag found")
				}
			}
			l++
		}
		// close(out)
	}()

}

func FindBasicData(in <-chan string) {
	fmt.Println("find basic data")
	go func() {
		for text := range in {
			for _, tag := range intro {
				r, _ := regexp.Compile(tag.regex)
				if r.MatchString(text) {
					fmt.Println("match!")
					fmt.Println(tag.name)
					switch tag.name {
					case "Name":
						fmt.Println(text)
						break

					case "Number":
						fmt.Println(text)
						break

					case "Aproved":
						fmt.Println(text)
						break

					case "Diary":
						fmt.Println(text)
						break

					}
				}
			}
		}
	}()
}

func ParseConcurrent(uri string) {
	fmt.Println("Parse concurrent")
	in := StreamLines(uri)
	chs := fanOut(in)
	FindBasicData(chs[0])
	FindCTags(chs[1])

}

// func printLanes(lines <-chan string) {
// 	for l := range lines {
// 		fmt.Println("channel 2 :", l)
// 	}
// }

func fanOut(ch <-chan string) []chan string {

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
}
