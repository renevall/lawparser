package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Parser interface {
	Parse(uri string) error
}

type Tesauro struct {
	Parser
}

//GetLaw process a GET request of a single Law
func (t *Tesauro) ParseTesauro(c *gin.Context) {

	err := t.Parser.Parse("Hola")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Record not Found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "worked"})
}
