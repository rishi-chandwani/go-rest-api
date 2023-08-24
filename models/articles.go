package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Articles struct {
	ID          primitive.ObjectID `json:"articleId" bson:"_id"`
	Title       string             `json:"articleTitle" bson:"articleTitle"`
	Description string             `json:"articleDesc" bson:"articleDesc"`
	Photos      map[string]string  `json:"photos" bson:"photos"`
	IsActive    bool               `json:"isActive" bson:"isActive"`
	ExpireOn    time.Time          `json:"expireOn" bson:"expireOn"`
	CreatedOn   time.Time          `json:"createdOn" bson:"createdOn"`
	ModifiedOn  time.Time          `json:"modifiedOn" bson:"modifiedOn"`
}
