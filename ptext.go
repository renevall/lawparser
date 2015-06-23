package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"sort"
	"strings"
)

type Tag struct {
	name  string
	regex string
}

type Tags []Tag

var title int = 0

func OpenTextFile(uri string) []string {
	file, err := ioutil.ReadFile(uri)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(file), "\n")
	return lines
}

func ParseText(uri string) {
	data := OpenTextFile(uri)
	ref, order := FindTags(data)
	MakeLaw(data, order, ref)
}

func FindTags(data []string) (map[int]string, []int) {
	var keys []int
	m := make(map[int]string)
	for i, text := range data {
		for _, tag := range tags {
			r, _ := regexp.Compile(tag.regex)
			if r.MatchString(text) {
				m[i+1] = tag.name
				break
			}
		}
	}

	for k, v := range m {
		keys = append(keys, k)
		if v == "Titulo" {
			title = title + 1
		}
	}

	sort.Ints(keys)

	return m, keys

}

func MakeLaw(data []string, index []int, ref map[int]string) {
	for _, k := range index {
		if ref[k] == "Titulo" {
			fmt.Println("Key:", k, "Value:", ref[k])

			// 	titles = append(titles, Title{name: data[k]})
		}
	}

	fmt.Println("Title total =", title)
	fmt.Println(data[5287])
}

var tags = Tags{
	Tag{"Titulo", "TÍTULO"},
	Tag{"Capitulo", "Capítulo\\s+[IV]+"},
	Tag{"Artículo", "Artículo\\s\\d+"},
	Tag{"Arto", "Art.\\s\\d+"},
}
