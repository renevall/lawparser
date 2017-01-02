package router

import (
	"net/http"

	"bitbucket.org/reneval/lawparser/domain"
	"bitbucket.org/reneval/lawparser/files"
	"github.com/gin-gonic/gin"
)

//LawReader interface reads Law via db package
type LawReader interface {
	GetLaws() ([]domain.Law, error)
}

//LawReader interface reads Law via file package
type LawJSONReader interface {
	LoadJSONLaw(name string) (*domain.Law, error)
}

type Law struct {
	Reader     LawReader
	JSONReader LawJSONReader
}

//GetLawsTMP responds with all tmp laws (flat file)
func (l *Law) GetLawsTMP(c *gin.Context) {
	laws, err := files.ListDirFiles("./parsed_laws")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": laws})
}

//GetLawsJSON responds with all laws from db
func (l *Law) GetLawsJSON(c *gin.Context) {
	var laws []domain.Law

	laws, err := l.Reader.GetLaws()
	if err != nil {
		c.JSON(http.StatusOK, err)
		return
	}

	if len(laws) < 1 {
		c.JSON(http.StatusOK, "")
		return
	}

	c.JSON(http.StatusOK, laws)
}

// ReadTMPLaw Reads a TMP Law (Flat file)  and renders it as JSON to be consumed
func (l *Law) ReadTMPLaw(c *gin.Context) {
	name := c.Param("name")
	law, err := l.JSONReader.LoadJSONLaw(name)

	if err != nil {
		c.JSON(http.StatusOK, err)
	}

	c.JSON(200, gin.H{"code": 200, "data": law})

}
