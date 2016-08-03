package auth

import (
	"fmt"
	"net/http"
	"strings"

	"bitbucket.org/reneval/lawparser/config"
	"bitbucket.org/reneval/lawparser/models"
	"bitbucket.org/reneval/lawparser/response"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

//LogIn User
func LogIn(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		res := response.Error{}

		email := r.FormValue("email")
		password := r.FormValue("password")
		fmt.Println(email)

		if strings.TrimSpace(email) == "" {
			res.Wrap(response.StatusError, "Email can't be empty")
			res.Respond(w)
		} else if strings.TrimSpace(password) == "" {
			res.Wrap(response.StatusError, "Password can't be empty")
			res.Respond(w)
		} else {
			user := new(models.User)
			err := user.FindByEmail(db, email)
			if err != nil {
				res.Wrap(response.StatusError, "Could not contact DB")
				res.Respond(w)
			}

			if user.Email == "" {
				res.Wrap(response.StatusError, "Account not found")

			} else {
				decryptedPass := decrypt([]byte(config.GetString("keys.crypt")), user.Password)
				if decryptedPass != password {
					res.Wrap(response.StatusError, "Account not found")
					res.Respond(w)
				} else {
					res.Wrap(response.StatusSuccess, "Account found!")
					res.Respond(w)
				}
			}
			//TODO Complete code to check for account
		}

	}
}

//create new user
func SignUp(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	}
}
