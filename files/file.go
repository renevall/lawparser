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
	"path/filepath"
	"strings"

	"bitbucket.org/reneval/lawparser/domain"

	"github.com/pkg/errors"
)

type FileReader struct {
}

func (f *FileReader) UploadFromHTTP(r *http.Request, dir string) (string, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}

	file, header, err := ReadFromHTTP(r)
	if err != nil {
		fmt.Println("Error on ReadFromHTTP", err)
		return "", err
	}
	defer file.Close()

	fname := strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))

	path, err := TempFile(dir, fname)
	if err != nil {
		fmt.Println("Error on TempFile", err)
		return "", err
	}
	defer os.Remove(path.Name())

	name, err := SaveUploadedFile(file, "./"+path.Name(), dir)
	if err != nil {
		fmt.Println("Error on SaveUploadedFile", err)
		return "", err
	}

	return name, nil
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
func SaveUploadedFile(file multipart.File, name string, dir string) (string, error) {
	fmt.Println("name is:", name)
	fmt.Println("dir is:", dir)
	path := dir + name
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return "", errors.Wrap(err, "Could not open tmp folder")
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Println("copy file")
		return "", errors.Wrap(err, "Could not copy file")
	}

	return path, nil

}

//OpenFile returns the file considering the uri it recieves
func OpenFile(uri string) ([]byte, error) {
	file, err := ioutil.ReadFile(uri)
	if err != nil {
		return nil, errors.Wrap(err, "Could not open file")
	}
	return file, nil
}

//DeleteFile removes a file from filesystem
func (f *FileReader) DeleteFile(uri string) error {
	err := os.Remove(uri)
	if err != nil {
		return err
	}
	return nil
}

//ListDirFiles list all files but dirs
func (f *FileReader) ListDirFiles(uri string) ([]domain.File, error) {
	files, err := ioutil.ReadDir(uri)
	var filelist []domain.File
	if err != nil {
		return nil, errors.Wrap(err, "Could not open Folder")
	}

	for _, f := range files {
		if !f.IsDir() {
			file := domain.File{Name: f.Name(), Path: uri}
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

//LoadJSONPub parses json file into a law object given a name.
func (f *FileReader) LoadJSONPub(name string) (*domain.Publication, error) {
	if name == "" {
		return nil, errors.Wrap(errors.New("Expected Param"), "Param name not set")
	}

	//TODO: use config file
	path := path.Join("./tmp_publication", name)

	file, err := OpenFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "file open failed")
	}
	pub := new(domain.Publication)

	//TODO: Review if it is posible to not unmarshall and send json from file
	//10ms diference so far
	err = json.Unmarshal(file, pub)
	if err != nil {
		return nil, errors.Wrap(err, "parsing to law failed")
	}
	return pub, nil
}
