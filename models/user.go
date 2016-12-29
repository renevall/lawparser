package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

//User struct mapped to DB
type User struct {
	UserName  string `json:"user_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	ContactNo string `json:"contact_no"`
	Status    string `json:"status"`
	UserLevel string `json:"user_level"`
	Password  string `json:"-"`
	Gender    string `json:"gender"`
	PicUrl    string `json:"pic_url"`
	BaseModel
}

//FindByEmail find the existence of an user by email
func (u *User) FindByEmail(db *sqlx.DB, email string) error {
	q := `SELECT * FROM User WHERE email = $1`
	err := db.Get(&u, q, email)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
