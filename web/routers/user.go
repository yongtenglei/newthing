package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/yongtenglei/newThing/web/logic"
	"github.com/yongtenglei/newThing/web/middlewares"
	"net/http"
)

func SetUp() *gin.Engine {
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// register
	router.POST("/register", logic.RegisterHandler)

	// login
	router.POST("/login", logic.LoginHandler)

	router.POST("/refreshToken", logic.RefreshTokenHandler)

	v1Group := router.Group("/v1/api/user", middlewares.JWTAuth())
	{
		// info
		v1Group.GET("/info/:mobile", logic.InfoHandler)
		// update
		v1Group.PUT("/update", logic.UpdateHandler)
		// delete
		v1Group.DELETE("/delete/:mobile", logic.DeleteHandler)
		// re password
		v1Group.PUT("/repassword", logic.RePasswordHandler)

	}
	return router

}
