package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDatabase() {
	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DATABASE_NAME")

	if uri == "" || dbName == "" {
		panic("Missing MONGO_URI or DATABASE_NAME in environment variables")
	}

	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic("Error connecting to DB: " + err.Error())
	}

	

	DB = client.Database(dbName)
	fmt.Println("✅ Connected to MongoDB:", dbName)
}
