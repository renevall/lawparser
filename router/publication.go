package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"os"

	"bitbucket.org/reneval/lawparser/domain"
	"github.com/gin-gonic/gin"
)

type fileparser interface {
	ParsePub(uri string) (*domain.Publication, error)
}
type filereader interface {
	LoadJSONPub(uri string) (*domain.Publication, error)
}
type fileuploader interface {
	UploadFromHTTP(*http.Request, string) (string, error)
}
type dirReader interface {
	ListDirFiles(string) ([]domain.File, error)
}

type Publication struct {
	fileparser
	filereader
	fileuploader
	dirReader
}

func NewPublication(
	parser fileparser,
	reader filereader,
	uploader fileuploader,
	dirReader dirReader,
) *Publication {
	return &Publication{
		fileparser:   parser,
		filereader:   reader,
		fileuploader: uploader,
		dirReader:    dirReader,
	}
}

//ParsePublication process a GET request of a single Law
func (p *Publication) ParsePublication(c *gin.Context) {

	publication, err := p.fileparser.ParsePub("Hola")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Record not Found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": publication})
}

//ListPublications returns a list of all Publications
func (p *Publication) ListPublications(c *gin.Context) {
	publication, err := p.filereader.LoadJSONPub(c.Param("id"))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Record not Found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": publication})
}

//UploadPublication uploads and parses a file coming from an http request
func (p *Publication) UploadPublication(c *gin.Context) {

	path, err := p.UploadFromHTTP(c.Request, "./tmp_publication/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not upload File"})
		return
	}

	pub, err := p.fileparser.ParsePub(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not parse File"})
		return
	}

	b, err := json.Marshal(pub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not marshall File"})
		return
	}

	name := strings.TrimSuffix(path, filepath.Ext(path))

	if err := ioutil.WriteFile(name+".json", b, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not write parsed File"})
	}
	os.Remove(path)

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": pub})
	// c.JSON(http.StatusOK, gin.H{"status": "success", "data": "Hola"})

}

func (p *Publication) GetTMPPub(c *gin.Context) {
	files, err := p.dirReader.ListDirFiles("./tmp_publication")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": files})
}
