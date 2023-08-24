package responses

import "github.com/rishi-chandwani/go-rest-api/models"

type UsersResponse struct {
	Status  int64             `json:"status"`
	Message string            `json:"msg"`
	Data    []models.UserRead `json:"details"`
}

type UserResponse struct {
	Status  int64           `json:"status"`
	Message string          `json:"msg"`
	Data    models.UserRead `json:"details"`
}
