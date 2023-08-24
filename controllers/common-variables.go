package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/rishi-chandwani/go-rest-api/configs"
)

var appDB = configs.GetDatabase(configs.MongoDb, "enterprise")
var usersCollection = configs.GetCollection(appDB, "users")
var vendorCollection = configs.GetCollection(appDB, "vendors")
var articleCollection = configs.GetCollection(appDB, "articles")

var validate = validator.New()
