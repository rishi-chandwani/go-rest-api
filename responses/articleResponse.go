package responses

import "github.com/rishi-chandwani/go-rest-api/models"

type ArticlesResponse struct {
	Status  int64             `json:"status"`
	Message string            `json:"msg"`
	Data    []models.Articles `json:"details"`
}

type ArticleResponse struct {
	Status  int64           `json:"status"`
	Message string          `json:"msg"`
	Data    models.Articles `json:"details"`
}
