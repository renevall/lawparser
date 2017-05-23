package main

import (
	"log"

	"bitbucket.org/reneval/lawparser/auth"
	"bitbucket.org/reneval/lawparser/domain"
	"bitbucket.org/reneval/lawparser/files"
	"bitbucket.org/reneval/lawparser/parser"
	db "bitbucket.org/reneval/lawparser/postgres"
	"bitbucket.org/reneval/lawparser/router"
	_ "github.com/lib/pq"
)

func main() {

	dataB, err := db.InitDB()
	if err != nil {
		log.Panic(err)
	}
	defer dataB.Close()

	//profile
	user := &db.User{dataB}
	env := &domain.Env{User: user}

	//auth
	reader := &auth.AuthService{UserReader: user}
	env.LoginReader = reader
	env.Authorizer = reader

	//law
	newLaw := &domain.Law{}
	law := &db.Law{dataB, newLaw}
	env.Law = law

	//Parser
	parser := &parser.Publication{}
	env.Parser = parser

	fileReader := &files.FileReader{}
	env.JSONFileReader = fileReader
	env.FileUploader = fileReader
	env.DirReader = fileReader
	env.FileRemover = fileReader

	router := router.InitRouter(env)
	router.Run(":8080")
}
