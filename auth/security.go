package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/reneval/lawparser/domain"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username string `json:"username"`

	// recommended having
	jwt.StandardClaims
}

//HashPass takes a pass (string) and returns a base64 hash
func HashPass(pass string) (string, error) {

	// TODO: use config files
	var cost = 1
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), cost)
	if err != nil {
		return "", errors.Wrap(err, "Could not generate pass hash")
	}

	// Encode the hash as base64 and return
	hashBase64 := base64.StdEncoding.EncodeToString(hash)

	return hashBase64, nil

}

// CompareHash comares 2 hashes passwords
func CompareHash(reqPass, dbPass string) bool {
	//Decode
	hashBytes, err := base64.StdEncoding.DecodeString(dbPass)
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword(hashBytes, []byte(reqPass))
	return err == nil

}

// GenerateToken generates a jwt token
func GenerateToken(user domain.User) string {
	// Expires the token and cookie in 1 hour
	expireToken := time.Now().Add(time.Hour * 24).Unix()
	//expireCookie := time.Now().Add(time.Hour * 1)

	// We'll manually assign the claims but in production you'd insert values from a database
	claims := Claims{
		user.Customer.Email,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "localhost:9000", //TODO: use real info
		},
	}

	// Create the token using your claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signs the token with a secret.
	//TODO: USE ENV for secret
	signedToken, _ := token.SignedString([]byte("secret"))

	return signedToken
}

//JWTAuthorize checks if token in header is valid
func (auth *AuthService) JWTAuthorize(r *http.Request) (domain.ClaimerVerifier, error) {
	token, err := request.ParseFromRequestWithClaims(
		r,
		request.AuthorizationHeaderExtractor,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			//TODO: See why is this secret
			return []byte("secret"), nil
		})

	if err != nil {
		err := errors.Wrap(err, "Token Invalid")
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil

	}

	return nil, err
}
