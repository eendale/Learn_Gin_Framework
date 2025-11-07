package db

import (
	"context"
	"time"

	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var DB  *mongo.Database 

func  ConnectDatabase(){
	clientOptions:=options.Client().ApplyURI(os.Getenv("MONGO_URI"))

	ctx, cancel:=context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err:=mongo.Connect(ctx, clientOptions)

	if err!=nil{
		panic("Error  Connecting to  DB")

	}

	DB = client.Database(os.Getenv("DATABASE_NAME"))

}