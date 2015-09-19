package api

import (
	"flag"
	"net/http"

	"bitbucket.com/reneval/lawparser/core"
	"github.com/gorilla/mux"
)

//NewRouter creates a mux router
func NewRouter(env *core.Env) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		// var handler http.Handler
		// handler = route.Func
		// handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(core.Handler{Env: env, H: route.Func})

	}

	dir := flag.String("directory", "app/", "directory of web files")
	flag.Parse()

	// handle all requests by serving a file of the same name
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)

	//handle custom 404

	//router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	router.Handle("/", http.RedirectHandler("/app/", 302))
	router.PathPrefix("/app/").Handler(
		NotFoundHook{http.StripPrefix("/app/", fileHandler)})
	return router
}
