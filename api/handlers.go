package api

import (
	"bitbucket.com/reneval/lawparser/core"
	"bitbucket.com/reneval/lawparser/parser"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//Index renders the index page
func Index(env *core.Env, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "Welcome!\n")
	return nil
}

//ParseShow shows the json of a parsed law
func ParseShow(env *core.Env, w http.ResponseWriter, r *http.Request) error {

	law := parser.ParseText("testlaws/test3.txt")

	w.Header().Set("Content-Type", "application/json; charset= UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(law); err != nil {
		panic(err)
	}
	return nil
}

//NotFoundHandler shows a custom 404
func NotFoundHandler(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("Not Found Handler", r.URL)
	return nil

}

//TestDBHandler helps to text db cnx
func TestDBHandler(env *core.Env, w http.ResponseWriter, r *http.Request) error {
	var law core.Law
	laws, err := law.GetAll(env)
	if err != nil {
		return core.StatusError{500, err}
	}

	for laws.Next() {
		err := laws.Scan(&law.Name, &law.Number, &law.Date, &law.Gaceta, &law.PublishDate, &law.ID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%+v", law)
	}

	return nil
}
