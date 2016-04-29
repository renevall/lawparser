package main

import (
	// "bitbucket.com/reneval/lawparser/api"
	"bitbucket.org/reneval/lawparser/db"
	"bitbucket.org/reneval/lawparser/handlers"
	// "bitbucket.org/reneval/lawparser/parser"
	"log"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := db.NewDB()
	if err != nil {
		log.Fatalln("Could not connect to database")
	}

	// router := api.NewRouter()

	r := httprouter.New()
	r.GET("/", handlers.Index)
	r.POST("/upload", handlers.FileUpload(db))
	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(":8080")
}
