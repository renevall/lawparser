package parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// type Article struct {
// 	Name    string
// 	Content string
// }

// type Articles []Article

// func main() {
// 	// data, err := openPDF("test.html")
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// parseLaw(data)

// 	ParseText("test3.txt")
// }

func openPDF(uri string) (string, error) {
	file, err := ioutil.ReadFile(uri)
	if err != nil {
		log.Fatal(err)
	}
	return string(file), nil
}

func parseLaw(data string) {
	// var articles Articles

	doc, err := html.Parse(strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	r, _ := regexp.Compile("Artículo\\s\\d+|Art.\\s\\d+")

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			match := r.MatchString(n.Data)
			if match {
				m := &html.Node{}
				m = n.NextSibling

				if m != nil {
					l := &html.Node{}
					l = m.NextSibling
					fmt.Println(l.Data)

				}

				// articles = append(articles, Article{Name: n.Data})

				// 	if b.Data == "br" && b.NextSibling.Data == "br" {
				// 		break
				// 	}
				// }
			}

		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	//fmt.Printf("%v+", articles[:10])

}
