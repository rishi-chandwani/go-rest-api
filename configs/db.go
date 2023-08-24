package configs

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectToDb function will be used to connect to MongoDB.
// If there is some error than app will panic and shutdown.
// If there is no error than this function will return MongoDB Client Object
func ConnectToDb() *mongo.Client {
	dbClient, connectErr := mongo.NewClient(options.Client().ApplyURI(GetEnvMongoLink()))
	if connectErr != nil {
		panic(connectErr)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	connectErr = dbClient.Connect(ctx)
	if connectErr != nil {
		panic(connectErr)
	}

	pingErr := dbClient.Ping(ctx, nil)
	if pingErr != nil {
		panic(pingErr)
	}

	log.Println("MongoDB Connection was successful")
	defer cancel()
	return dbClient
}

var MongoDb = ConnectToDb()

// GetDatabase function accepts pointer to MongoDB Client and Database Name as string to which app should connect
// It will return pointer to MongoDB Database object
// TODO - Check what happens if Database itself doesn't exists?
func GetDatabase(dbClient *mongo.Client, databaseName string) *mongo.Database {
	return dbClient.Database(databaseName)
}

// GetCollection function accepts pointer to MongoDB Database object and Collection Name as string
// It will return pointer to MongoDB Collection object
// TODO - Check what happens if Collection doesn't exists?
func GetCollection(database *mongo.Database, collectionName string) *mongo.Collection {
	return database.Collection(collectionName)
}
