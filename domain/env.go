package domain

import "net/http"

//Env holds injected instances
type Env struct {
	User           UserStore
	LoginReader    LoginReader
	Authorizer     JWTAuthorizer
	Law            LawStore
	JSONFileReader JSONFileReader
	Parser         Parser
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

type Parser interface {
	Parse(uri string) error
}
