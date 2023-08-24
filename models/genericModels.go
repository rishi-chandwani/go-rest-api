package models

import "time"

type CommonFilters struct {
	CreatedOn  time.Time `json:"createdOn" bson:"createdOn"`
	ModifiedOn time.Time `json:"modifiedOn" bson:"modifiedOn"`

	PageNo int64          `json:"pageNo"`
	Limit  int64          `json:"limit"`
	Sort   map[string]int `json:"sort"`
}

type UserFilters struct {
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	Email     string `json:"email" bson:"email"`
	IsActive  int    `json:"isActive" bson:"isActive"`

	GenericFilters CommonFilters `json:"filters"`
}

type VendorFilters struct {
	Name   string `json:"vendorName" bson:"vendorName"`
	Mobile string `json:"mobile" bson:"mobile"`
	Email  string `json:"email" bson:"email"`
	Gstn   string `json:"gstNumber" bson:"gstNumber"`

	GenericFilters CommonFilters `json:"filters"`
}

type ArticleFilters struct {
	Title string `json:"articleTitle" bson:"articleTitle"`

	GenericFilters CommonFilters `json:"filters"`
}
