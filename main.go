package main

import (
	"time"

	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var DB *mongo.Database


func ConnectDatabase(){
	clientOptions  := options.Client().ApplyURI("mongodb://localhost:27017")

	ctx, cancel:= context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client,err:=mongo.Connect(ctx, clientOptions)

	if err!=nil{
		panic(err)
	}

	DB= client.Database("testdb")
} 


type User struct {
	ID primitive.ObjectID   `bson:"_id, omitempty" json:"id"`
	Name string          `bson:"name" json:"name"`
	Age int             `bson:"age" json:"age"`

}


