package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	// "github.com/gorilla/mux"
	"net/http"

	"bitbucket.org/reneval/lawparser/models"
	"bitbucket.org/reneval/lawparser/parser"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

type Response struct {
	err  string
	code bool
}

type Response2 struct {
	status        bool
	originalName  string
	generatedName string
}

func Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func ParseShow(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	law := parser.ParseText("testlaws/test3.txt")

	w.Header().Set("Content-Type", "application/json; charset= UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(law); err != nil {
		panic(err)
	}
}

func GetAllArticles(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var article models.Article
		articles, err := article.GetArticles(db)
		if err != nil {
			log.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset= UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(articles); err != nil {
			panic(err)
		}

	}

}

func GetFullLawJSON(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var law models.Law
		err := law.GetFullLaw(db, 1)
		if err != nil {
			log.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset= UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(law); err != nil {
			panic(err)
		}

	}
}
func FileUpload(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// w.Header().Set("Access-Control-Allow-Origin", "*")

		log.Println("METHOD IS " + r.Method + " AND CONTENT-TYPE IS " + r.Header.Get("Content-Type"))
		r.ParseMultipartForm(32 << 20)
		fmt.Println(r.MultipartForm.File)

		file, handler, err := r.FormFile("uploads[]")
		if err != nil {
			fmt.Println(Response{err.Error(), true})
			fmt.Println("open file")
			json.NewEncoder(w).Encode(Response{err.Error(), true})

			return
		}
		defer file.Close()
		// fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./tmp/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			json.NewEncoder(w).Encode(Response{err.Error(), true})
			fmt.Println(Response{err.Error(), true})
			fmt.Println("create file")

			return
		}
		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			fmt.Println("copy file")
			json.NewEncoder(w).Encode(Response{err.Error(), true})
			return
		}
		fmt.Println(Response{"File '" + handler.Filename + "' submited successfully", false})
		// json.NewEncoder(w).Encode(Response2{true, handler.Filename + "Server", handler.Filename})

		law := parser.ParseText("testlaws/" + handler.Filename)
		inserted := parser.InsertLawToDB(db, law)
		if inserted != nil {
			panic(inserted)
		}
		w.Header().Set("Content-Type", "application/json; charset= UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(law); err != nil {
			panic(err)
		}
	}
}
