package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	// "io"
	// "os"
	// "github.com/gorilla/mux"
	"bitbucket.com/reneval/lawparser/parser"
	"net/http"
)

type appHandler func(http.ResponseWriter, *http.Request) (int, error)

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if status, err := fn(w, r); err != nil {
		switch status {
		case http.StatusNotFound:
			NotFound(w, r)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

func Index(w http.ResponseWriter, r *http.Request) (int, error) {
	fmt.Fprint(w, "Welcome!\n")

	return http.StatusOK, nil
}

func ParseShow(w http.ResponseWriter, r *http.Request) (int, error) {

	law := parser.ParseText("testlaws/test3.txt")

	w.Header().Set("Content-Type", "application/json; charset= UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(law); err != nil {
		panic(err)
	}

	// if err != nil {
	// 	// Much better!
	// 	return http.StatusInternalServerError, err
	// }

	return http.StatusOK, nil
}

func FileUpload(w http.ResponseWriter, r *http.Request) (int, error) {
	fmt.Println(r.Header)
	r.ParseForm()
	file, header, err := r.FormFile("file")

	if err != nil {
		// Much better!
		return http.StatusInternalServerError, err
	}

	name := r.PostFormValue("lawName")
	fmt.Fprintf(w, "Hello, %s!", name)
	fmt.Println(file)
	fmt.Println(header)

	return http.StatusOK, nil

}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	http.ServeFile(w, r, "app/index.html")
	// fmt.Println("What?")
}

// FILE SERVER CUSTOM HANDLER//

type myFileHandler struct {
	http.ResponseWriter
	ignore bool
}

func (mfh *myFileHandler) WriteHeader(status int) {
	//mfh.ResponseWriter.WriteHeader(200)
	if status == 404 {
		mfh.ignore = true
		t, _ := template.ParseFiles("app/index.html")
		t.Execute(mfh.ResponseWriter, nil)

	}

}

func (mfh *myFileHandler) Write(p []byte) (int, error) {
	if mfh.ignore {
		return len(p), nil
	}
	return mfh.ResponseWriter.Write(p)
}

type NotFoundHook struct {
	h http.Handler
}

func (nfh NotFoundHook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	nfh.h.ServeHTTP(&myFileHandler{ResponseWriter: w}, r)
}

// END CUSTOM HANDLER //
