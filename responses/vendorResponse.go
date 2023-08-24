package responses

import "github.com/rishi-chandwani/go-rest-api/models"

type VendorsResponse struct {
	Status  int64           `json:"status"`
	Message string          `json:"message"`
	Data    []models.Vendor `json:"vendorDetails"`
}

type VendorResponse struct {
	Status  int64         `json:"status"`
	Message string        `json:"message"`
	Data    models.Vendor `json:"vendorDetails"`
}
