package router

import (
	"net/http"

	"bitbucket.org/reneval/lawparser/files"
	"github.com/gin-gonic/gin"
)

type Law struct {
}

//GetLawsTMP responds with all tmp laws (flat file)
func (l *Law) GetLawsTmp(c *gin.Context) {
	laws, err := files.ListDirFiles("./parsed_laws")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, laws)
}
