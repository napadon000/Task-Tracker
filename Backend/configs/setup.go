package configs

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	dotenv "github.com/dsh2dsh/expx-dotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var MongoClient *mongo.Client
var UserCollection *mongo.Collection
var TaskCollection *mongo.Collection

func init() {
	if err := dotenv.New().Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func ConnectDatabase() {

	clientOption := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))

	client, err := mongo.Connect(clientOption)

	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// and returns an error if the ping fails
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		// log.Printf("Failed to ping MongoDB: %v", err)
		// return
		panic(err)
	}

	fmt.Println("Connected to MongoDB!")
	MongoClient = client

	UserCollection = client.Database("tasktracker").Collection("user")
	TaskCollection = client.Database("tasktracker").Collection("task")
}
