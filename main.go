package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rishi-chandwani/go-rest-api/configs"
	"github.com/rishi-chandwani/go-rest-api/routes"
)

func main() {
	router := gin.Default()

	configs.ConnectToDb()

	routes.UserRouters(router)
	routes.VendorRoutes(router)
	routes.ArticleRoutes(router)
	router.Run(":" + configs.GetApplicationPort())
}
