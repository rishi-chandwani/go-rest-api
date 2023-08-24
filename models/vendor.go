package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Vendor struct {
	ID         primitive.ObjectID `json:"vendorId" bson:"_id"`
	Name       string             `json:"vendorName" bson:"vendorName"`
	Address    string             `json:"address" bson:"address"`
	Mobile     string             `json:"mobile" bson:"mobile"`
	Email      string             `json:"email" bson:"email"`
	Gstn       string             `json:"gstNumber" bson:"gstNumber"`
	IsActive   bool               `json:"isActive" bson:"isActive"`
	CreatedOn  time.Time          `json:"createdOn" bson:"createdOn"`
	ModifiedOn time.Time          `json:"modifiedOn" bson:"modifiedOn"`
}
