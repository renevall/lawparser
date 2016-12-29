package router

import (
	"net/http"
	"time"

	"bitbucket.org/reneval/lawparser/domain"

	//remove dependecy

	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
)

// InitRouter initializes the router
func InitRouter(env *domain.Env) *gin.Engine {
	router := gin.Default()

	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	profile := &Profile{service: env.User}
	router.GET("/profile/:id", profile.ProfileHandler)

	return router
}

//NotImplemented is a test dummy for handler
func NotImplemented(c *gin.Context) {
	content := gin.H{"Response": "Not Implemented"}
	c.JSON(http.StatusOK, content)
}
