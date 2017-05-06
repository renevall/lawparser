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
	parser := &parser.Tesauro{}
	env.Parser = parser

	fileReader := &files.FileReader{}
	env.JSONFileReader = fileReader

	router := router.InitRouter(env)
	router.Run(":8080")

	// db, err := database.NewDB()
	// if err != nil {
	// 	log.Fatalln("Could not connect to database")
	// }
	// c := cors.New(cors.Options{
	// 	AllowedOrigins: []string{"http://localhost:4200"},
	// 	AllowedHeaders: []string{"Accept", "Content-Type",
	// 		"Content-Length", "Accept-Encoding", "X-CSRF-Token",
	// 		"Authorization", "Origin"},
	// })
	// r := httprouter.New()
	// // r.POST("/api/upload", handlers.FileUpload(db))

	// // TODO: CREATE A ROUTER MODULE FOR CLEANER MAIN FILE

	// //Law
	// r.GET("/api/laws", handlers.GetLawsJSON(db))            //Get all Laws
	// r.GET("/api/law/:id", handlers.GetFullLawJSON(db))      //Get a law header
	// r.GET("/api/law/:id/full", handlers.GetFullLawJSON(db)) //Get a full law
	// r.GET("/api/tmp/laws", handlers.GetLawsTMP())           //Get tmp law (before veryfing it)
	// r.GET("/api/tmp/law/:name", handlers.ReadTMPLaw())      //Get tmp law (before veryfing it)

	// r.POST("/api/laws/parse", handlers.ParseLawFile(db)) //parse law
	// r.POST("/api/laws", handlers.ParseLawFile(db))       //create new law

	// r.GET("/api/laws/save/:name", handlers.SaveReviewedToDB(db)) //create new law

	// r.PUT("/api/laws/:id", handlers.GetFullLawJSON(db))    //Update Law
	// r.PATCH("/api/laws/:id", handlers.GetFullLawJSON(db))  //Partially Update Law
	// r.DELETE("/api/laws/:id", handlers.GetFullLawJSON(db)) //Delete Law

	// //Article
	// r.GET("/api/articles", handlers.GetAllArticles(db))

	// r.POST("/api/login", auth.LogIn(db))

	// //r.POST("/api/setToken", auth.setToken())

	// n := negroni.New(
	// 	negroni.NewRecovery(),
	// 	negroni.NewLogger(),
	// 	negroni.NewStatic(http.Dir("app")),
	// )
	// n.Use(c)
	// n.UseHandler(r)
	// n.Run(":8080")
}
