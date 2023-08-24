package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rishi-chandwani/go-rest-api/models"
	"github.com/rishi-chandwani/go-rest-api/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllVendors() gin.HandlerFunc {
	return func(cntx *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		var allVendors []models.Vendor
		var vendorFilters models.VendorFilters
		defer cancel()

		// Validate request data with User Model
		if bindErr := cntx.BindJSON(&vendorFilters); bindErr != nil {
			cntx.JSON(http.StatusBadRequest, responses.VendorsResponse{Status: http.StatusBadRequest, Message: "Error - Invalid Request Data - " + bindErr.Error(), Data: allVendors})
			return
		}

		findOptions := options.Find()
		if vendorFilters.GenericFilters.Limit != 0 && vendorFilters.GenericFilters.PageNo != 0 {
			startFrom := (vendorFilters.GenericFilters.Limit * vendorFilters.GenericFilters.PageNo) - vendorFilters.GenericFilters.Limit
			findOptions.SetLimit(int64(vendorFilters.GenericFilters.Limit))
			findOptions.SetSkip(int64(startFrom))
		}

		if vendorFilters.GenericFilters.Sort != nil {
			findOptions.SetSort(vendorFilters.GenericFilters.Sort)
		}

		filterOptions := getVendorSearchFilters(vendorFilters)

		allVendorCur, fetchErr := vendorCollection.Find(ctx, filterOptions, findOptions)
		if fetchErr != nil {
			cntx.JSON(http.StatusInternalServerError, responses.VendorsResponse{Status: http.StatusInternalServerError, Message: "Error - Fetching all users - " + fetchErr.Error(), Data: allVendors})
			return
		}

		defer allVendorCur.Close(ctx)

		for allVendorCur.Next(ctx) {
			var oneVendor models.Vendor

			if oneFetchErr := allVendorCur.Decode(&oneVendor); oneFetchErr != nil {
				cntx.JSON(http.StatusInternalServerError, responses.VendorsResponse{Status: http.StatusInternalServerError, Message: "Error - Generating return data - " + oneFetchErr.Error(), Data: allVendors})
				return
			}

			allVendors = append(allVendors, oneVendor)
		}

		cntx.JSON(http.StatusOK, responses.VendorsResponse{Status: http.StatusOK, Message: "Details fetched successfully", Data: allVendors})
	}
}

func AddNewVendor() gin.HandlerFunc {
	return func(cntx *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		var vendor models.Vendor
		var readVendor models.Vendor
		defer cancel()

		if bindErr := cntx.BindJSON(&vendor); bindErr != nil {
			cntx.JSON(http.StatusBadRequest, responses.VendorResponse{Status: http.StatusBadRequest, Message: "Error - Invalid request data - " + bindErr.Error(), Data: vendor})
			return
		}

		if validateErr := validate.Struct(&vendor); validateErr != nil {
			cntx.JSON(http.StatusBadRequest, responses.VendorResponse{Status: http.StatusBadRequest, Message: "Error - Missing required data - " + validateErr.Error(), Data: vendor})
			return
		}

		existFilter := bson.M{"vendorName": vendor.Name, "mobile": vendor.Mobile, "email": vendor.Email}
		chkErr := vendorCollection.FindOne(ctx, existFilter).Decode(&readVendor)
		if chkErr == nil {
			cntx.JSON(http.StatusBadRequest, responses.VendorResponse{Status: http.StatusBadRequest, Message: "Error - Vendor with same Name OR Mobile OR Email already exists", Data: readVendor})
			return
		}

		vendor.ID = primitive.NewObjectID()
		vendor.CreatedOn = time.Now()
		vendor.IsActive = true

		insData, insErr := vendorCollection.InsertOne(ctx, vendor)
		if insErr != nil {
			cntx.JSON(http.StatusInternalServerError, responses.VendorResponse{Status: http.StatusInternalServerError, Message: "Error - Insert failed - " + insErr.Error(), Data: vendor})
			return
		}

		log.Println(insData)

		fetchInsVendor := bson.M{"_id": vendor.ID}
		fetchErr := vendorCollection.FindOne(ctx, fetchInsVendor).Decode(&readVendor)
		if fetchErr != nil {
			cntx.JSON(http.StatusInternalServerError, responses.VendorResponse{Status: http.StatusInternalServerError, Message: "Error - Fetching newly inserted vendor details - " + fetchErr.Error(), Data: vendor})
			return
		}

		cntx.JSON(http.StatusOK, responses.VendorResponse{Status: http.StatusOK, Message: "Vendor added successfully", Data: readVendor})
	}
}

func GetVendorById() gin.HandlerFunc {
	return func(cntx *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		var vendor models.Vendor
		var vendorId = cntx.Param("id")
		defer cancel()

		vendorIdObj, convErr := primitive.ObjectIDFromHex(vendorId)
		if convErr != nil {
			cntx.JSON(http.StatusInternalServerError, responses.VendorResponse{Status: http.StatusInternalServerError, Message: "Error - while converting to Object - " + convErr.Error(), Data: vendor})
			return
		}

		vendorFilter := bson.M{"_id": vendorIdObj}
		fetchErr := vendorCollection.FindOne(ctx, vendorFilter).Decode(&vendor)
		if fetchErr != nil {
			cntx.JSON(http.StatusInternalServerError, responses.VendorResponse{Status: http.StatusInternalServerError, Message: "Error - Fetching details failed - " + fetchErr.Error(), Data: vendor})
			return
		}

		cntx.JSON(http.StatusOK, responses.VendorResponse{Status: http.StatusOK, Message: "Details fetched successfully", Data: vendor})
	}
}

func getVendorSearchFilters(vendorFilters models.VendorFilters) map[string]string {
	var filters = map[string]string{}

	if vendorFilters.Email != "" {
		filters["email"] = vendorFilters.Email
	}

	if vendorFilters.Gstn != "" {
		filters["gstNumber"] = vendorFilters.Gstn
	}

	if vendorFilters.Mobile != "" {
		filters["mobile"] = vendorFilters.Mobile
	}

	if vendorFilters.Name != "" {
		filters["vendorName"] = vendorFilters.Name
	}

	return filters
}
