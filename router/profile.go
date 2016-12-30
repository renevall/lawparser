package router

import (
	"net/http"
	"strconv"

	"bitbucket.org/reneval/lawparser/domain"
	"github.com/gin-gonic/gin"
)

// //Loader retrieves an User from DB
// type Loader interface {
// 	Load(id string) (*domain.User, error)
// }

//ProfileRetriever reads a Profile via db package
type ProfileRetriever interface {
	FindByID(id uint64) (*domain.User, error)
}

//Profile Holds a User Profile Data
type Profile struct {
	service ProfileRetriever
}

//New Returns a User Profile
//TODO: Check if this is needed
func New(p ProfileRetriever) *Profile {
	return &Profile{
		service: p,
	}
}

//ProfileHandler deals with a Profile Request
func (p *Profile) ProfileHandler(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	user, err := p.service.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, "Profile Not Found")
	}
	c.JSON(http.StatusOK, user)

}
