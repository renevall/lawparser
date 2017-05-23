package router

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"fmt"

	"bitbucket.org/reneval/lawparser/domain"
	"github.com/gin-gonic/gin"
)

//LawReader interface reads Law via db package
type LawReaderWriter interface {
	GetLaws() ([]domain.Law, error)
	GetLaw(id string) (domain.Law, error)
	InsertLawDB(*domain.Law) error
	AutoComplete(query string) ([]string, error)
}

//LawReader interface reads Law via file package
type LawJSONReader interface {
	LoadJSONLaw(name string) (*domain.Law, error)
}

type DirReader interface {
	ListDirFiles(string) ([]domain.File, error)
}

type FileRemover interface {
	DeleteFile(string) error
}

type Law struct {
	ReaderWriter LawReaderWriter
	JSONReader   LawJSONReader
	DirReader
	FileRemover
}

//GetLawsTMP responds with all tmp laws (flat file)
func (l *Law) GetLawsTMP(c *gin.Context) {
	// TODO: Use Global Config
	files, err := l.DirReader.ListDirFiles("./parsed_laws")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": files})
}

//GetLawsJSON responds with all laws from db
func (l *Law) GetLawsJSON(c *gin.Context) {
	var laws []domain.Law

	laws, err := l.ReaderWriter.GetLaws()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error",
			"message": "Could not retrieve Laws"})
		return
	}

	if len(laws) < 1 {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": laws})
}

//GetLaw process a GET request of a single Law
func (l *Law) GetLaw(c *gin.Context) {
	id := c.Param("id")
	law, err := l.ReaderWriter.GetLaw(id)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Record not Found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": law})
}

func (l *Law) AutoComplete(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Not the expected params"})
		return
	}
	results, err := l.ReaderWriter.AutoComplete(query)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error",
			"message": "Could not reach DB"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": results})

}

//IndexLaw sends a Request to the indexer service to Index a Law
func (l *Law) IndexLaw(c *gin.Context) {
	url := "http://localhost:8585/law"
	id := c.Param("id")
	law, err := l.ReaderWriter.GetLaw(id)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Record not Found"})
		return
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(law)
	_, err = http.Post(url, "application/json; charset=utf-8", b)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Problem with indexer service"})
	}
	c.JSON(http.StatusOK, gin.H{})
}

// ReadTMPLaw Reads a TMP Law (Flat file)  and renders it as JSON to be consumed
func (l *Law) ReadTMPLaw(c *gin.Context) {
	name := c.Param("name")
	law, err := l.JSONReader.LoadJSONLaw(name)

	if err != nil {
		c.JSON(http.StatusOK, err)
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": law})

}

func (l *Law) SaveLawDB(c *gin.Context) {
	law := &domain.Law{}
	err := c.BindJSON(&law)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Ilegal JSON request"})
		return
	}
	err = l.ReaderWriter.InsertLawDB(law)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not save to DB"})
		return
	}

	l.FileRemover.DeleteFile("./parsed_laws/" + law.Init)

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": law})
}

func (l *Law) UpdateTmpLaw(c *gin.Context) {
	dir := "./parsed_laws/"
	name := c.Param("name")
	law := &domain.Law{}
	err := c.BindJSON(&law)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "error": "Ilegal JSON request"})
		return
	}

	b, err := json.Marshal(law)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	if err := ioutil.WriteFile(dir+name, b, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": law})

}
