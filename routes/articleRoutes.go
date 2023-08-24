package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rishi-chandwani/go-rest-api/controllers"
)

func ArticleRoutes(router *gin.Engine) {
	router.GET("/articles", controllers.GetAllArticles())
	router.POST("/article", controllers.CreateNewArticle())
	router.GET("/article/:articleId", controllers.GetArticleById())
}
