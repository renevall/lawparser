package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"bitbucket.org/reneval/lawparser/files"
	"bitbucket.org/reneval/lawparser/parser"
	"github.com/gin-gonic/gin"
)

type File struct {
}

//LawParseHandler handles the parse request coming from the web
func (f *File) LawParseHandler(c *gin.Context) {
	//TODO: use path from config file
	dir := "./parsed_laws/"

	file, header, err := files.ReadFromHTTP(c.Request)
	defer file.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	fname := strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))
	path, err := files.TempFile(dir, fname)
	defer os.Remove(path.Name())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	err = files.SaveUploadedFile(file, header.Filename, "./tmp/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	//TODO: Make the parser return an error, check if using request input
	law := parser.ParseConcurrent("./tmp/" + header.Filename)
	b, err := json.Marshal(law)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	if err := ioutil.WriteFile(path.Name()+".json", b, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, "Save succeded")

}

//ParseLawFile recieves a Law in txt format from a http request, parses it,
//then proceeds to save it to a JSON file
// func ParseLawFile(db *sqlx.DB) httprouter.Handle {
// 	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
// 		// w.Header().Set("Access-Control-Allow-Origin", "*")

// 		//TODO: use path from config file
// 		dir := "./parsed_laws/"

// 		log.Println("METHOD IS " + r.Method + " AND CONTENT-TYPE IS " + r.Header.Get("Content-Type"))

// 		file, header, err := files.ReadFromHTTP(r)
// 		defer file.Close()

// 		if err != nil {
// 			log.Println(err)
// 			res := response.Error{}
// 			res.Wrap(response.StatusError, err.Error())
// 			if err := json.NewEncoder(w).Encode(res); err != nil {
// 				panic(err)

// 			}
// 		}
// 		fname := strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))

// 		path, err := files.TempFile(dir, fname)
// 		defer os.Remove(path.Name())
// 		fmt.Println(path)

// 		if err != nil {
// 			log.Println(err)
// 			res := response.Error{}
// 			res.Wrap(response.StatusError, err.Error())
// 			if err := json.NewEncoder(w).Encode(res); err != nil {
// 				panic(err)

// 			}
// 		}

// 		err = files.SaveUploadedFile(file, header.Filename, "./tmp/")
// 		if err != nil {
// 			log.Println(err)
// 			res := response.Error{}
// 			res.Wrap(response.StatusError, err.Error())
// 			if err := json.NewEncoder(w).Encode(res); err != nil {
// 				panic(err)

// 			}
// 		}

// 		law := parser.ParseConcurrent("testlaws/" + header.Filename)

// 		res := response.Response{}
// 		res.Wrap(response.StatusSuccess, "Saved Succed")

// 		b, err := json.Marshal(law)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		fmt.Println("path: ", path.Name())
// 		if err := ioutil.WriteFile(path.Name()+".json", b, 0644); err != nil {
// 			log.Println(err)
// 			res := response.Error{}
// 			res.Wrap(response.StatusError, err.Error())
// 			if err := json.NewEncoder(w).Encode(res); err != nil {
// 				panic(err)

// 			}
// 		}

// 		//TODO: SavetoJson instead of priting to screen
// 		if err := json.NewEncoder(w).Encode(res); err != nil {
// 			panic(err)

// 		}

// 	}
// }
