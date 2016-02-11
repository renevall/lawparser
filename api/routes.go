package api

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/data",
		Index,
	},
	Route{
		"ParseShow",
		"GET",
		"/parse",
		ParseShow,
	},
	Route{
		"FileUpload",
		"POST",
		"/upload",
		FileUpload,
	},
}
