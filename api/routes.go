package api

import (
	"bitbucket.com/reneval/lawparser/core"
	"net/http"
)

type Route struct {
	Name    string
	Method  string
	Pattern string
	Func    func(e *core.Env, w http.ResponseWriter, r *http.Request) error
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/data",
		Index,
	},
	// Route{
	// 	"TodoIndex",
	// 	"GET",
	// 	"/todos",
	// 	TodoIndex,
	// },
	Route{
		"ParseShow",
		"GET",
		"/parse",
		ParseShow,
	},
	Route{
		"TestDB",
		"GET",
		"/testdb",
		TestDBHandler,
	},
}
