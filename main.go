package main

import (
	// "bitbucket.com/reneval/lawparser/api"

	"net/http"

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

	r := httprouter.New()
	r.POST("/api/upload", handlers.FileUpload(db))
	r.GET("/api/articles", handlers.GetAllArticles(db))
	r.GET("/api/law", handlers.GetFullLawJSON(db))
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.NewStatic(http.Dir("app")),
	)
	n.UseHandler(r)
	n.Run(":8080")
}
