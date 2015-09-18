package api

import (
	"encoding/json"
	"fmt"
	// "github.com/gorilla/mux"
	"bitbucket.com/reneval/lawparser/parser"
	"net/http"
)

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

func NotFoundHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("Not Found Handler" ,r.URL )
	
}
