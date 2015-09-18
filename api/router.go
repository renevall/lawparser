package api

import (
	"flag"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	dir := flag.String("directory", "app/", "directory of web files")
	flag.Parse()

	// handle all requests by serving a file of the same name
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)
	
	//handle custom 404
	
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	router.Handle("/", http.RedirectHandler("/app/", 302))
	router.PathPrefix("/app/").Handler(
		NotFoundHook{http.StripPrefix("/app/",fileHandler)})
	// router.PathPrefix("/app").Handler(http.StripPrefix("/app/",fileHandler))
	return router
}
