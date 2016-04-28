package main

import (
	// "bitbucket.com/reneval/lawparser/api"
	"bitbucket.com/reneval/lawparser/db"
	"bitbucket.com/reneval/lawparser/handlers"
	// "bitbucket.com/reneval/lawparser/parser"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func main() {

	db, err := db.NewDB()
	if err != nil {
		log.Fatalln("Could not connect to database")
	}
	
	// router := api.NewRouter()
	
	r:= httprouter.New()
	r.GET("/", handlers.Index)
	r.POST("/upload", handlers.FileUpload(db))

	log.Fatal(http.ListenAndServe(":8080", r))
}
