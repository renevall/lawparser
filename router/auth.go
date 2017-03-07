package router

import (
	"fmt"
	"net/http"

	"bitbucket.org/reneval/lawparser/domain"
	"github.com/gin-gonic/gin"
)

//LoginRequest used to map request via gin
type LoginRequest struct {
	Email    string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

//LoginReader used to decouple code
type LoginReader interface {
	Login(email, pass string) (*domain.User, error)
}

//RequestAuth tries to auth the user via the postgres package
type RequestAuth struct {
	LoginReader LoginReader
	Authorizer  domain.JWTAuthorizer
}

//AuthHandler deals with the Auth Request
func (r *RequestAuth) AuthHandler(c *gin.Context) {
	var login LoginRequest
	err := c.BindJSON(&login)
	if err != nil {
		fmt.Println(err)
		return
	}
	//user, err := context.AuthService.Login(login.Email, login.Password)
	user, err := r.LoginReader.Login(login.Email, login.Password) //Dependency
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}

// NOTES

// The router uses a database layer on r.Service.Login, but using interfaces
// we decoupled the code. The main file should inject an object that implements
// the User Reader interface(with a Login Method). At the same times we have a
// dependency on the data layer, so when main wire up things,we should have the
// corresponding data dependency injected
