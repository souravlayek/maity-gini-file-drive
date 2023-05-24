package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MyClient *mongo.Client
var DB *mongo.Database

func ConnectDB() {
	dbPath := os.Getenv("MONGO_URI")
	if dbPath == "" {
		log.Fatal("MONGO_URI not found")
	}
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	connectionURI := dbPath
	clientOptions := options.Client().
		ApplyURI(connectionURI).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	MyClient = client
	DB = client.Database("storage")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database connected successfully")
}
