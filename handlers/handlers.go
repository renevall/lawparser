package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	// "github.com/gorilla/mux"
	"net/http"

	"bitbucket.org/reneval/lawparser/files"
	"bitbucket.org/reneval/lawparser/models"
	"bitbucket.org/reneval/lawparser/parser"
	"bitbucket.org/reneval/lawparser/response"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
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

//GetLawsTMP responds with all tmp laws (flat file)
func GetLawsTMP() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json; charset= UTF-8")
		w.WriteHeader(http.StatusOK)
		laws, err := files.ListDirFiles("./parsed_laws")
		if err != nil {
			fmt.Println(Response{err.Error(), true})
			fmt.Println("open file")
			json.NewEncoder(w).Encode(Response{err.Error(), true})

			return
		}

		res := response.Response{}

		res.Wrap(response.StatusSuccess, laws)

		json.NewEncoder(w).Encode(res)
	}
}

//GetLawsJSON responds with all laws from db
func GetLawsJSON(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var law models.Law
		var laws []models.Law
		laws, err := law.GetLaws(db)
		if err != nil {
			log.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset= UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(laws); err != nil {
			panic(err)
		}

	}
}

//TODO: MAKE SAVE FROM JSON FILE
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

//ParseLawFile recieves a Law in txt format from a http request, parses it,
//then proceeds to save it to a JSON file
func ParseLawFile(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// w.Header().Set("Access-Control-Allow-Origin", "*")

		//TODO use path from config file
		dir := "./parsed_laws/"

		log.Println("METHOD IS " + r.Method + " AND CONTENT-TYPE IS " + r.Header.Get("Content-Type"))

		file, header, err := files.ReadFromHTTP(r)
		defer file.Close()

		if err != nil {
			log.Println(err)
			res := response.Error{}
			res.Wrap(response.StatusError, err.Error())
			if err := json.NewEncoder(w).Encode(res); err != nil {
				panic(err)

			}
		}
		fname := strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))

		path, err := files.TempFile(dir, fname)
		defer os.Remove(path.Name())
		fmt.Println(path)

		if err != nil {
			log.Println(err)
			res := response.Error{}
			res.Wrap(response.StatusError, err.Error())
			if err := json.NewEncoder(w).Encode(res); err != nil {
				panic(err)

			}
		}

		err = files.SaveUploadedFile(file, header.Filename, "./tmp/")
		if err != nil {
			log.Println(err)
			res := response.Error{}
			res.Wrap(response.StatusError, err.Error())
			if err := json.NewEncoder(w).Encode(res); err != nil {
				panic(err)

			}
		}

		law := parser.ParseConcurrent("testlaws/" + header.Filename)

		res := response.Response{}
		res.Wrap(response.StatusSuccess, "Saved Succed")

		b, err := json.Marshal(law)
		if err != nil {
			log.Println(err)
		}

		fmt.Println("path: ", path.Name())
		if err := ioutil.WriteFile(path.Name()+".json", b, 0644); err != nil {
			log.Println(err)
			res := response.Error{}
			res.Wrap(response.StatusError, err.Error())
			if err := json.NewEncoder(w).Encode(res); err != nil {
				panic(err)

			}
		}

		//TODO SavetoJson instead of priting to screen
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)

		}

	}
}

//TODO READ FROM TEMP LAW

func ReadTMPLaw() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		name := p.ByName("name")
		if name == "" {
			res := response.Error{}

			err := errors.Wrap(errors.New("Expected Param"), "Param name not present")
			res.Wrap(response.StatusError, err.Error())

			if err := json.NewEncoder(w).Encode(res); err != nil {
				panic(err)
			}
		}

		path := path.Join("./parsed_laws", name)
		file, err := files.OpenFile(path)
		if err != nil {
			res := response.Error{}
			res.Wrap(response.StatusError, err.Error())
			if err := json.NewEncoder(w).Encode(res); err != nil {
				panic(err)
			}
			return
		}
		law := new(models.Law)

		//TODO Review if it is posible to not unmarshall and send json from file
		//10ms diference so far
		json.Unmarshal(file, law)

		res := response.Response{}
		res.Wrap(response.StatusSuccess, law)

		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}

	}
}
