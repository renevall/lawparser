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

	//TODO: Use Groups
	profile := &Profile{service: env.User}
	router.GET("/profile/:id", profile.ProfileHandler)

	auth := &RequestAuth{LoginReader: env.LoginReader, Authorizer: env.Authorizer}
	router.POST("/login", auth.AuthHandler)

	test := router.Group("/test", auth.TokenAuthMiddleware())
	{
		test.GET("/logout", NotImplemented)
	}

	file := &File{}
	router.POST("/api/tmp/laws", file.LawParseHandler)

	law := &Law{ReaderWriter: env.Law, JSONReader: env.JSONFileReader}
	router.GET("/api/tmp/laws", law.GetLawsTMP)
	router.GET("/api/laws", law.GetLawsJSON)
	router.GET("/api/laws/:id", law.GetLaw)
	router.GET("/api/tmp/laws/:name", law.ReadTMPLaw)
	router.POST("/api/laws", law.SaveLawDB)
	router.PUT("/api/tmp/laws/:name", law.UpdateTmpLaw)

	router.GET("/api/law/autocomplete", law.AutoComplete)

	router.GET("/api/index/law/:id", law.IndexLaw)

	return router
}

//NotImplemented is a test dummy for handler
func NotImplemented(c *gin.Context) {
	content := gin.H{"Response": "Not Implemented"}
	c.JSON(http.StatusOK, content)
}

// NOTES

// Originally we would wire up the dependencies on this file,
// but it felt out of place, and we were going against the decoupling
// flow. Ideally we would wire up things on the main. We came up
// with and Env structure defined at the domain level, where we would
// set the dependencies as member, the we pass it to the router package
// where we prepare our handlers with our requires dependencies. This way
// so far makes it so router has no coupling with the data layer.
