package main

import (
	// "bitbucket.com/reneval/lawparser/api"

	"net/http"

	"bitbucket.org/reneval/lawparser/auth"
	"bitbucket.org/reneval/lawparser/database"
	"bitbucket.org/reneval/lawparser/handlers"
	// "bitbucket.org/reneval/lawparser/parser"
	"log"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {

	db, err := database.NewDB()
	if err != nil {
		log.Fatalln("Could not connect to database")
	}
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200"},
		AllowedHeaders: []string{"Accept", "Content-Type",
			"Content-Length", "Accept-Encoding", "X-CSRF-Token",
			"Authorization", "Origin"},
	})
	r := httprouter.New()
	// r.POST("/api/upload", handlers.FileUpload(db))

	// TODO: CREATE A ROUTER MODULE FOR CLEANER MAIN FILE

	//Law
	r.GET("/api/laws", handlers.GetLawsJSON(db))            //Get all Laws
	r.GET("/api/law/:id", handlers.GetFullLawJSON(db))      //Get a law header
	r.GET("/api/law/:id/full", handlers.GetFullLawJSON(db)) //Get a full law
	r.GET("/api/tmp/laws", handlers.GetLawsTMP())           //Get tmp law (before veryfing it)
	r.GET("/api/tmp/law/:name", handlers.ReadTMPLaw())      //Get tmp law (before veryfing it)

	r.POST("/api/laws/parse", handlers.ParseLawFile(db)) //parse law
	r.POST("/api/laws", handlers.ParseLawFile(db))       //create new law

	r.GET("/api/laws/save/:name", handlers.SaveReviewedToDB(db)) //create new law

	r.PUT("/api/laws/:id", handlers.GetFullLawJSON(db))    //Update Law
	r.PATCH("/api/laws/:id", handlers.GetFullLawJSON(db))  //Partially Update Law
	r.DELETE("/api/laws/:id", handlers.GetFullLawJSON(db)) //Delete Law

	//Article
	r.GET("/api/articles", handlers.GetAllArticles(db))

	r.POST("/api/login", auth.LogIn(db))

	//r.POST("/api/setToken", auth.setToken())

	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.NewStatic(http.Dir("app")),
	)
	n.Use(c)
	n.UseHandler(r)
	n.Run(":8080")
}
