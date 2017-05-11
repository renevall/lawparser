package router

import (
	"fmt"
	"net/http"

	"bitbucket.org/reneval/lawparser/domain"

	"github.com/gin-gonic/gin"
)

type Parser interface {
	Parse(uri string) (*domain.Tesauro, error)
}

type Tesauro struct {
	Parser
}

//GetLaw process a GET request of a single Law
func (t *Tesauro) ParseTesauro(c *gin.Context) {

	tesauro, err := t.Parser.Parse("Hola")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Record not Found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tesauro})
}
