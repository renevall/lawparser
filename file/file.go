//Package file provides the method for file handling.
package file

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

//ReadFromHTTP reads the input files from an http request object
func ReadFromHTTP(r *http.Request) (multipart.File, *multipart.FileHeader, error) {
	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("uploads[]")
	if err != nil {
		fmt.Println("file could not be opened")
		return nil, nil, errors.Wrap(err, "Could not parse file from request")
	}

	defer file.Close()

	return file, header, nil

}

//SaveFile picks a file a moves it to the requested location
func SaveFile(file multipart.File, name string, path string) error {
	f, err := os.OpenFile("./tmp/"+name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		//json.NewEncoder(w).Encode(Response{err.Error(), true})
		//fmt.Println(Response{err.Error(), true})
		fmt.Println("Could not open tmp folder")

		return errors.Wrap(err, "Could not open tmp folder")
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Println("copy file")
		return errors.Wrap(err, "Could not copy file")
	}

	return nil

}
