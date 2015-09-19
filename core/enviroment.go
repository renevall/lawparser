package core

import (
	"database/sql"
	"net/http"
)

//Env is A (simple) example of our application-wide configuration.
type Env struct {
	DB   *sql.DB
	Port string
	Host string
}

// The Handler struct that takes a configured Env and a function matching
// our useful signature.
type Handler struct {
	Env *Env
	H   func(e *Env, w http.ResponseWriter, r *http.Request) error
}
