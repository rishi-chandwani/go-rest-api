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

func GetAllUsers() gin.HandlerFunc {
	return func(cntx *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		var allUsers []models.UserRead
		var userFilters models.UserFilters
		defer cancel()

		// Validate request data with User Model
		if bindErr := cntx.BindJSON(&userFilters); bindErr != nil {
			cntx.JSON(http.StatusBadRequest, responses.UsersResponse{Status: http.StatusBadRequest, Message: "Error - Invalid Request Data - " + bindErr.Error(), Data: allUsers})
			return
		}

		findOptions := options.Find()
		if userFilters.GenericFilters.Limit != 0 && userFilters.GenericFilters.PageNo != 0 {
			startFrom := (userFilters.GenericFilters.Limit * userFilters.GenericFilters.PageNo) - userFilters.GenericFilters.Limit
			findOptions.SetLimit(int64(userFilters.GenericFilters.Limit))
			findOptions.SetSkip(int64(startFrom))
		}

		if userFilters.GenericFilters.Sort != nil {
			findOptions.SetSort(userFilters.GenericFilters.Sort)
		}

		filterOptions := getUserSearchFilters(userFilters)

		allUsersCur, fetchErr := usersCollection.Find(ctx, filterOptions, findOptions)
		if fetchErr != nil {
			cntx.JSON(http.StatusInternalServerError, responses.UsersResponse{Status: http.StatusInternalServerError, Message: "Error - Fetching all users - " + fetchErr.Error(), Data: allUsers})
			return
		}

		defer allUsersCur.Close(ctx)

		for allUsersCur.Next(ctx) {
			var oneUser models.UserRead

			if oneFetchErr := allUsersCur.Decode(&oneUser); oneFetchErr != nil {
				cntx.JSON(http.StatusInternalServerError, responses.UsersResponse{Status: http.StatusInternalServerError, Message: "Error - Generating return data - " + oneFetchErr.Error(), Data: allUsers})
				return
			}

			allUsers = append(allUsers, oneUser)
		}

		cntx.JSON(http.StatusOK, responses.UsersResponse{Status: http.StatusOK, Message: "Details fetched successfully", Data: allUsers})
	}
}

func CreateUser() gin.HandlerFunc {
	return func(cntx *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		var user models.User
		var userRead models.UserRead
		defer cancel()

		// Validate request data with User Model
		if bindErr := cntx.BindJSON(&user); bindErr != nil {
			cntx.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Error - Invalid Request Data - " + bindErr.Error(), Data: userRead})
			return
		}

		// Validate all the required fields as mentioned in Model are present in the request
		if missErr := validate.Struct(&user); missErr != nil {
			cntx.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Error - Missing Mandatory Details", Data: userRead})
			return
		}

		existFilter := bson.M{"firstName": user.FirstName, "lastName": user.LastName, "email": user.Email}
		chkErr := usersCollection.FindOne(ctx, existFilter).Decode(&userRead)
		if chkErr == nil {
			cntx.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Error - User with same First Name OR Last Name OR Email already exists", Data: userRead})
			return
		}

		user.Id = primitive.NewObjectID()
		user.CreatedAt = time.Now()
		user.IsActive = 1

		_, insErr := usersCollection.InsertOne(ctx, user)
		if insErr != nil {
			cntx.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error - Something went wrong while adding new User", Data: userRead})
			return
		}

		fetchInsUser := bson.M{"_id": user.Id}
		fetchErr := usersCollection.FindOne(ctx, fetchInsUser).Decode(&userRead)
		if fetchErr != nil {
			cntx.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error - While fetching user details after Insertion", Data: userRead})
			return
		}

		cntx.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "User added successfully", Data: userRead})
	}
}

func GetUserById() gin.HandlerFunc {
	return func(cntx *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		userId := cntx.Param("id")
		var userRead models.UserRead
		defer cancel()

		userObjId, convErr := primitive.ObjectIDFromHex(userId)
		if convErr != nil {
			cntx.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error - While converting to Object - " + convErr.Error(), Data: userRead})
			return
		}

		fetchErr := usersCollection.FindOne(ctx, bson.M{"_id": userObjId}).Decode(&userRead)
		if fetchErr != nil {
			cntx.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error - While fetching user details", Data: userRead})
			return
		}

		cntx.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: " User details fetched successfully", Data: userRead})
	}
}

func getUserSearchFilters(userFilters models.UserFilters) bson.D {
	var filterData = make(bson.D, 3)

	if userFilters.FirstName != "" {
		filterData = append(filterData, bson.E{Key: "firstName", Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: userFilters.FirstName, Options: "i"}}}})
	}

	if userFilters.LastName != "" {
		filterData = append(filterData, bson.E{Key: "lastName", Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: userFilters.LastName, Options: "i"}}}})
	}

	if userFilters.Email != "" {
		filterData = append(filterData, bson.E{Key: "email", Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: userFilters.Email, Options: "i"}}}})
	}

	return filterData
}
