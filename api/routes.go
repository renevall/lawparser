package api

// import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc appHandler
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/data",
		appHandler(Index),
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
		appHandler(ParseShow),
	},
	Route{
		"FileUpload",
		"POST",
		"/api/law/uploadfile",
		appHandler(FileUpload),
	},
}
