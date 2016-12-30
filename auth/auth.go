package auth

import (
	"errors"

	"bitbucket.org/reneval/lawparser/domain"
)

//AuthService helps with dependency injection and decoupling
type AuthService struct {
	UserReader domain.UserStore
}

//Login logs the user
func (auth *AuthService) Login(email, pass string) (*domain.User, error) {
	var authUser *domain.User
	user, err := auth.UserReader.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if CompareHash(pass, user.Password) {
		authUser, err = auth.UserReader.FindByID(user.ID)
		if err != nil {
			return nil, err
		}
		authUser.Token = GenerateToken(*authUser)
	} else {
		err = errors.New("hash not equal")
		return nil, err
	}

	return authUser, nil

}

// NOTES

// auth package depends on a database layer(postgres), so that when DI,
// the instance that is passed should contain the db method implementation
// that it is set on the interface
