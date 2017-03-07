//Package files provides the method for file handling.
package files

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	"bitbucket.org/reneval/lawparser/domain"
	"bitbucket.org/reneval/lawparser/models"

	"github.com/pkg/errors"
)

type FileReader struct {
}

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

//SaveUploadedFile picks a file a moves it to the requested location
func SaveUploadedFile(file multipart.File, name string, path string) error {
	f, err := os.OpenFile(path+name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
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

//OpenFile returns the file considering the uri it recieves
func OpenFile(uri string) ([]byte, error) {
	file, err := ioutil.ReadFile(uri)
	if err != nil {
		return nil, errors.Wrap(err, "Could not open file")
	}
	return file, nil
}

//ListDirFiles list all files but dirs
func ListDirFiles(uri string) ([]models.TmpLaw, error) {
	files, err := ioutil.ReadDir(uri)
	var filelist []models.TmpLaw
	if err != nil {
		return nil, errors.Wrap(err, "Could not open Folder")
	}

	for _, f := range files {
		if !f.IsDir() {
			file := models.TmpLaw{Name: f.Name(), Path: uri}
			filelist = append(filelist, file)
		}
	}

	return filelist, nil
}

//LoadJSONLaw parses json file into a law object given a name.
//Uses Config File for folder path
func (f *FileReader) LoadJSONLaw(name string) (*domain.Law, error) {

	if name == "" {
		return nil, errors.Wrap(errors.New("Expected Param"), "Param name not set")
	}

	//TODO: use config file
	path := path.Join("./parsed_laws", name)

	file, err := OpenFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "file open failed")
	}
	law := new(domain.Law)

	//TODO: Review if it is posible to not unmarshall and send json from file
	//10ms diference so far
	err = json.Unmarshal(file, law)
	if err != nil {
		return nil, errors.Wrap(err, "parsing to law failed")
	}
	return law, nil
}
