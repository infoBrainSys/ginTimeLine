package route

import (
	"github.com/gin-gonic/gin"
	"timeLineGin/internal/controller"
	"timeLineGin/internal/logic/middleware"
)

func RegisterRoute(e *gin.Engine) {
	e.Any("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
	e.Any("/auth", middleware.Authorize(), func(c *gin.Context) { c.JSON(200, gin.H{"message": "pass"}) })

	e.POST("/signUp", controller.SignUp)
	e.POST("/signIn", middleware.SignTokenString(), controller.SignIn)
	e.DELETE("/signOut", middleware.AddInBlackLists())

}
