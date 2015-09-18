package main

import (
	"bitbucket.com/reneval/lawparser/api"
	"bitbucket.com/reneval/lawparser/db"
	// "bitbucket.com/reneval/lawparser/parser"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/context"
	"log"
	"net/http"
)

func main() {

	ctx := context.Background()
	ctx = db.OpenSQL(ctx, "main", "sqlite3", "root:hunter2@unix(/tmp/mysql.sock)/myCoolDB")
	defer db.Close(ctx) // closes all DB connections

	router := api.NewRouter()
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", router))

}
