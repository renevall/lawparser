package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	// "github.com/gorilla/mux"
	"net/http"

	"bitbucket.com/reneval/lawparser/parser"
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

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func ParseShow(w http.ResponseWriter, r *http.Request) {

	law := parser.ParseText("testlaws/test3.txt")

	w.Header().Set("Content-Type", "application/json; charset= UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(law); err != nil {
		panic(err)
	}
}

func FileUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
	json.NewEncoder(w).Encode(Response2{true, handler.Filename + "Server", handler.Filename})
}
