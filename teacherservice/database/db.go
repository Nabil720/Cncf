package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

func Connect() error {
	// Environment variables থেকে MongoDB URI নিবে
	// (Vault থেকে main.go-তে already set করা হয়েছে)
	connectionString := os.Getenv("MONGODB_URI")
	if connectionString == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

	databaseName := os.Getenv("DATABASE_NAME")
	if databaseName == "" {
		databaseName = "kindergarten"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return err
	}

	Client = client
	Database = client.Database(databaseName)
	log.Println("Connected to MongoDB successfully!")
	return nil
}

func GetCollection(collectionName string) *mongo.Collection {
	return Database.Collection(collectionName)
}