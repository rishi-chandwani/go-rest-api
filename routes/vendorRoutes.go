package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rishi-chandwani/go-rest-api/controllers"
)

func VendorRoutes(router *gin.Engine) {
	router.GET("/vendors", controllers.GetAllVendors())
	router.POST("/vendor", controllers.AddNewVendor())
	router.GET("/vendor/:id", controllers.GetVendorById())
}
