package main

import (
	"bitbucket.com/reneval/lawparser/api"
	"bitbucket.com/reneval/lawparser/core"

	// "bitbucket.com/reneval/lawparser/parser"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {

	// ctx := context.Background()
	cs := "postgres://penshiru:dream015@localhost/penshiru?sslmode=disable"
	db, err := sql.Open("postgres", cs)
	if err != nil {
		log.Fatal(err)
	}

	env := &core.Env{
		DB: db,
		// Port: os.Getenv("PORT"),
		// Host: os.Getenv("HOST"),
	}

	// ctx = db.OpenSQL(ctx, "main", "postgres", ds)
	//
	// defer db.Close(ctx) // closes all DB connections

	router := api.NewRouter(env)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", router))

}
