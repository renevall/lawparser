package router

import "github.com/gin-gonic/gin"

//TODO: DELETE, USING DEPENDY INJECTION
// Database is a middleware that add a db cnx to the context
// func Database(db *sqlx.DB) gin.HandlerFunc {
// 	fmt.Println("adding db to context")
// 	return func(c *gin.Context) {
// 		c.Set("DB", db)
// 		c.Next()
// 	}
// }

// type Authorizer interface {
// 	JWTAuthorize(r *http.Request) (interface {
// 		Valid() error
// 	}, error)
// }

// type Claimer interface {
// 	New() interface {
// 		Valid() error
// 	}
// }

// TokenAuthMiddleware exists to protect /profile and /logout
func (auth *RequestAuth) TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, err := auth.Authorizer.JWTAuthorize(c.Request)
		if err != nil {
			c.JSON(401, gin.H{"message": "Token Required", "status": 401})
			c.Abort()
		}
		c.Set("Claims", claim)
		c.Next()
	}
}
