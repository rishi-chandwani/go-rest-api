package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rishi-chandwani/go-rest-api/controllers"
)

func UserRouters(router *gin.Engine) {
	router.GET("/users", controllers.GetAllUsers())
	router.POST("/user", controllers.CreateUser())
	router.GET("/user/:id", controllers.GetUserById())
}
