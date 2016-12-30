package domain

//Env holds injected instances
type Env struct {
	User UserStore
	Auth Authorizer
}
