package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

func SignIn() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Write([]byte("Not Implemented"))

	}
}

func setToken() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Write([]byte("Not Implemented"))

		// Expires the token and cookie in 24 hours
		expireToken := time.Now().Add(time.Hour * 24).Unix()
		expireCookie := time.Now().Add(time.Hour * 24)

		//TODO authenticate to db here

		// We'll manually assign the claims but in production you'd insert values from a database
		claims := &jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "penshiru.io",
		}

		// Create the token using your claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Signs the token with a secret.
		//TODO read From config enviroment
		signedToken, _ := token.SignedString([]byte("secret"))

		// This cookie will store the token on the client side
		cookie := http.Cookie{Name: "Auth", Value: signedToken, Expires: expireCookie, HttpOnly: true}

		http.SetCookie(w, &cookie)

		// Redirect the user to his profile
		http.Redirect(w, r, "/profile", 301)

	}
}
