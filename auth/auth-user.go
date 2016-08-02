package auth

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

//AuthUser hold the data once a user is signed
type AuthUser struct {
	ID        int       `json:"id"`
	UserName  string    `json:"user_name"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Status    string    `json:"status"`
	UserLevel string    `json:"user_level"`
	Email     string    `json:"email"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Gender    string    `json:"gender"`
	PicURL    string    `json:"pic_url"`

	Claim jwt.StandardClaims
}
