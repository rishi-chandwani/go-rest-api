package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rishi-chandwani/go-rest-api/models"
	"github.com/rishi-chandwani/go-rest-api/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllArticles() gin.HandlerFunc {
	return func(cntx *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		var allArticles []models.Articles
		var articleFilters models.ArticleFilters
		defer cancel()

		// Validate request data with Articles Model
		if bindErr := cntx.BindJSON(&articleFilters); bindErr != nil {
			cntx.JSON(http.StatusBadRequest, responses.ArticlesResponse{Status: http.StatusBadRequest, Message: "Error - Invalid Request Data - " + bindErr.Error(), Data: allArticles})
			return
		}

		findOptions := options.Find()
		if articleFilters.GenericFilters.Limit != 0 && articleFilters.GenericFilters.PageNo != 0 {
			startFrom := (articleFilters.GenericFilters.Limit * articleFilters.GenericFilters.PageNo) - articleFilters.GenericFilters.Limit
			findOptions.SetLimit(int64(articleFilters.GenericFilters.Limit))
			findOptions.SetSkip(int64(startFrom))
		}

		if articleFilters.GenericFilters.Sort != nil {
			findOptions.SetSort(articleFilters.GenericFilters.Sort)
		}

		filterOptions := getArticlesSearchFilters(articleFilters)

		allArticlesCur, fetchErr := articleCollection.Find(ctx, filterOptions, findOptions)
		if fetchErr != nil {
			cntx.JSON(http.StatusInternalServerError, responses.ArticlesResponse{Status: http.StatusInternalServerError, Message: "Error - Fetching all articles - " + fetchErr.Error(), Data: allArticles})
			return
		}

		defer allArticlesCur.Close(ctx)

		for allArticlesCur.Next(ctx) {
			var oneArticle models.Articles

			if oneFetchErr := allArticlesCur.Decode(&oneArticle); oneFetchErr != nil {
				cntx.JSON(http.StatusInternalServerError, responses.ArticlesResponse{Status: http.StatusInternalServerError, Message: "Error - Generating return data - " + oneFetchErr.Error(), Data: allArticles})
				return
			}

			allArticles = append(allArticles, oneArticle)
		}

		cntx.JSON(http.StatusOK, responses.ArticlesResponse{Status: http.StatusOK, Message: "Details fetched successfully", Data: allArticles})
	}
}

func CreateNewArticle() gin.HandlerFunc {
	return func(cntx *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		var article models.Articles
		var articleRead models.Articles
		defer cancel()

		if bindErr := cntx.BindJSON(&article); bindErr != nil {
			cntx.JSON(http.StatusBadRequest, responses.ArticleResponse{Status: http.StatusBadRequest, Message: "Error - Invalid Data - " + bindErr.Error(), Data: articleRead})
			return
		}

		if validateErr := validate.Struct(&article); validateErr != nil {
			cntx.JSON(http.StatusBadRequest, responses.ArticleResponse{Status: http.StatusBadRequest, Message: "Error - Mandatory fields not present - " + validateErr.Error(), Data: articleRead})
			return
		}

		article.ID = primitive.NewObjectID()
		article.CreatedOn = time.Now()
		article.IsActive = true
		article.ExpireOn = time.Now().Add(time.Hour * 10)

		_, insErr := articleCollection.InsertOne(ctx, article)
		if insErr != nil {
			cntx.JSON(http.StatusInternalServerError, responses.ArticleResponse{Status: http.StatusInternalServerError, Message: "Error - Insert failed - " + insErr.Error(), Data: article})
			return
		}

		fetchFilter := bson.M{"_id": article.ID}
		fetchErr := articleCollection.FindOne(ctx, fetchFilter).Decode(&articleRead)
		if fetchErr != nil {
			cntx.JSON(http.StatusInternalServerError, responses.ArticleResponse{Status: http.StatusInternalServerError, Message: "Error - Fetching inserted details - " + fetchErr.Error(), Data: article})
			return
		}

		cntx.JSON(http.StatusOK, responses.ArticleResponse{Status: http.StatusOK, Message: "Article added successfully", Data: articleRead})
	}
}

func GetArticleById() gin.HandlerFunc {
	return func(cntx *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		var articleDetails models.Articles
		var articleId = cntx.Param("articleId")
		defer cancel()

		idObj, convErr := primitive.ObjectIDFromHex(articleId)
		if convErr != nil {
			cntx.JSON(http.StatusInternalServerError, responses.ArticleResponse{Status: http.StatusInternalServerError, Message: "Error - converting ID to Object - " + convErr.Error(), Data: articleDetails})
			return
		}

		fetchFilter := bson.M{"_id": idObj}
		fetchErr := articleCollection.FindOne(ctx, fetchFilter).Decode(&articleDetails)
		if fetchErr != nil {
			cntx.JSON(http.StatusOK, responses.ArticleResponse{Status: http.StatusNoContent, Message: "No Data Found", Data: articleDetails})
			return
		}

		cntx.JSON(http.StatusOK, responses.ArticleResponse{Status: http.StatusOK, Message: "Details fetched successfully", Data: articleDetails})
	}
}

func getArticlesSearchFilters(articleFilters models.ArticleFilters) bson.D {
	var filterData = make(bson.D, 1)

	if articleFilters.Title != "" {
		filterData = append(filterData, bson.E{Key: "articleTitle", Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: articleFilters.Title, Options: "i"}}}})
	}

	return filterData
}
