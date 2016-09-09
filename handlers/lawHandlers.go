package handlers

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/reneval/lawparser/database"
	"bitbucket.org/reneval/lawparser/files"
	"bitbucket.org/reneval/lawparser/response"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

//SaveReviewedToDB saves json file(reviewed) given a http request name param
func SaveReviewedToDB(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		name := p.ByName("name")

		law, err := files.LoadJSONLaw(name)

		if err != nil {
			res := response.Error{}
			res.Wrap(response.StatusError, err.Error())
			if err := json.NewEncoder(w).Encode(res); err != nil {
				panic(err)
			}
		}

		database.InsertLawToDB(db, law)

	}
}
