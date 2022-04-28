package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/yongtenglei/newThing/web/logic"
)

func SetUp() *gin.Engine {
	router := gin.Default()

	v1Group := router.Group("/v1/api/user")
	{
		// register
		v1Group.POST("/register", logic.RegisterHandler)
		// login
		v1Group.POST("/login", logic.LoginHandler)
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
