package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id         primitive.ObjectID `json:"userId" bson:"_id"`
	FirstName  string             `json:"firstName" bson:"firstName"`
	LastName   string             `json:"lastName" bson:"lastName"`
	UserName   string             `json:"userName" bson:"userName"`
	Email      string             `json:"email" bson:"email"`
	Password   string             `json:"password" bson:"password"`
	IsActive   int                `json:"isActive" bson:"isActive"`
	CreatedAt  time.Time          `json:"createdOn" bson:"createdOn"`
	ModifiedAt time.Time          `json:"modifiedOn" bson:"modifiedOn"`
}

type UserRead struct {
	Id         primitive.ObjectID `json:"userId" bson:"_id"`
	FirstName  string             `json:"firstName" bson:"firstName"`
	LastName   string             `json:"lastName" bson:"lastName"`
	Email      string             `json:"email" bson:"email"`
	IsActive   int                `json:"isActive" bson:"isActive"`
	CreatedAt  time.Time          `json:"createdOn" bson:"createdOn"`
	ModifiedAt time.Time          `json:"modifiedOn" bson:"modifiedOn"`
}
