package main

import (
	"bitbucket.com/reneval/lawparser/api"
	"log"
	"net/http"
)

func main() {

	router := api.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
