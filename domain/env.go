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
	FileUploader
	DirReader
	FileRemover
}

type ClaimerVerifier interface {
	Valid() error
}

type JWTAuthorizer interface {
	JWTAuthorize(r *http.Request) (ClaimerVerifier, error)
}

type JSONFileReader interface {
	LoadJSONLaw(name string) (*Law, error)
	LoadJSONPub(name string) (*Publication, error)
}

type FileUploader interface {
	UploadFromHTTP(*http.Request, string) (string, error)
}

type Parser interface {
	ParsePub(uri string) (*Publication, error)
}

type DirReader interface {
	ListDirFiles(string) ([]File, error)
}

type FileRemover interface {
	DeleteFile(string) error
}
