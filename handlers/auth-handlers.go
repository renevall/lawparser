package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func SignIn() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Write([]byte("Not Implemented"))

	}
}

func SignUp() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Write([]byte("Not Implemented"))

	}
}
