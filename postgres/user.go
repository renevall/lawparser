package postgres

import (
	"fmt"

	"bitbucket.org/reneval/lawparser/domain"
	"github.com/pkg/errors"
)

//User implements the UserStore interface
type User struct {
	*DB
}

// //New returns a new db interface
// func New(s domain.UserStore) *User {
// 	return &User{
// 		store: s,
// 	}
// }

//FindByID return one user with the query ID
func (u *User) FindByID(id uint64) (*domain.User, error) {
	var user = domain.User{}

	query := `SELECT user_id, first_name, last_name, email 
              FROM "user" WHERE user_id = $1;`
	// err := db.sqler.DB().Get(&user, query, id)
	err := u.DB.Get(&user, query, id)
	if err != nil {
		err = errors.Wrap(err, "couldn't find user by id")
		return nil, err
	}
	return &user, nil

}

//FindByEmail finds an User by his email
func (u *User) FindByEmail(email string) (*domain.User, error) {
	fmt.Println("llego a postgres")
	var user domain.User
	query := `SELECT user_id, password FROM "user" WHERE "email" = $1`
	//err := db.sqler.DB().Get(&user, query, email)
	err := u.DB.Get(&user, query, email)
	if err != nil {
		err = errors.Wrap(err, "couldn't find user by email")
		return nil, err
	}

	return &user, nil

}

//FindAll returns all users in DB.
func (u *User) FindAll() ([]*domain.User, error) {
	fmt.Println("llego a postgres")
	var users []*domain.User
	query := `SELECT * FROM "user"`
	err := u.DB.Select(&users, query)
	if err != nil {
		err = errors.Wrap(err, "couldn't find any user")
		return nil, err
	}

	return users, nil

}
