package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
)

//Auth logs the user in using JWT
func Auth(db *sqlx.DB, email string, pass string) {

}

func generateJWT(username string) string {
	mySigningKey := []byte(config.GetString("TOKEN_KEY"))
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Issuer:    username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(mySigningKey)
	return tokenString
}
