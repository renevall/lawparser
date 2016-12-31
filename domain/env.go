package domain

import "net/http"

//Env holds injected instances
type Env struct {
	User        UserStore
	LoginReader LoginReader
	Authorizer  JWTAuthorizer
	Law         LawStore
}

type ClaimerVerifier interface {
	Valid() error
}

type JWTAuthorizer interface {
	JWTAuthorize(r *http.Request) (ClaimerVerifier, error)
}
