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

type Law struct {
	Reader LawReader
}

//GetLawsTMP responds with all tmp laws (flat file)
func (l *Law) GetLawsTMP(c *gin.Context) {
	laws, err := files.ListDirFiles("./parsed_laws")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, laws)
}

//GetLawsJSON responds with all laws from db
func (l *Law) GetLawsJSON(c *gin.Context) {
	var laws []domain.Law

	laws, err := l.Reader.GetLaws()
	if err != nil {
		c.JSON(http.StatusOK, err)
	}

	c.JSON(http.StatusOK, laws)
}
