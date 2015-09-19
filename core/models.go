package core

import (
	"database/sql"
)

//Law is the representation of the Law Table in the DB
type Law struct {
	ID          sql.NullInt64
	Name        string
	Number      int
	Date        sql.NullString
	Gaceta      string
	PublishDate sql.NullString
	// Titles      []Title
}

//GetAll returns all laws in db
func (model *Law) GetAll(env *Env) (*sql.Rows, error) {
	query := "SELECT * FROM \"Law\""
	rows, err := env.DB.Query(query)
	if err != nil {
		return rows, err
	}

	return rows, nil
}

//GetOne return one Law
func (model *Law) GetOne(env *Env, id int) {

}
