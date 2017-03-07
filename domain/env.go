package domain

import "net/http"

//Env holds injected instances
type Env struct {
	User           UserStore
	LoginReader    LoginReader
	Authorizer     JWTAuthorizer
	Law            LawStore
	JSONFileReader JSONFileReader
}

type ClaimerVerifier interface {
	Valid() error
}

type JWTAuthorizer interface {
	JWTAuthorize(r *http.Request) (ClaimerVerifier, error)
}

type JSONFileReader interface {
	LoadJSONLaw(name string) (*Law, error)
}
