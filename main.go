package main

import (
	"time"

	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database


func ConnectDatabase(){
	clientOptions  := options.Client().ApplyURI("mongodb://localhost:27017")

	ctx, cancel:= context.WithTimeout(context.Background(), 10*time.Second)
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


func  CreateUser(c *gin.Context){
	collection:=DB.Collection("users")
	var  user User
   err:=c.ShouldBindJSON(&user)
   if err!=nil{
	c.JSON(400, gin.H{"error":err.Error()})
	return 
   }
   ctx , cancel :=context.WithTimeout(context.Background(), 5*time.Second)
   defer cancel()

   result, err:=collection.InsertOne(ctx, user)

   if err!=nil{
	c.JSON(500, gin.H{"error":err.Error()})
	return 
   }

   user.ID = result.InsertedID.(primitive.ObjectID)

   c.JSON(201, user)
}


func  GetUsers(c *gin.Context){
	collection:=DB.Collection("users")

	ctx, cancel:=context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err:=collection.Find(ctx, bson.M{})
	if err!=nil{
		c.JSON(500, gin.H{"error":err.Error()})
		return
	}

	var  users []User

	err=cursor.All(ctx, &users)

	if err!=nil{
		c.JSON(500, gin.H{"error":err.Error()})
		return
	}

	c.JSON(200, users)

}

//routes

func Routes(r *gin.Engine){
	r.POST("/users", CreateUser)
	r.GET("/users", GetUsers)
}

func  main(){

	r:=gin.Default()
   ConnectDatabase()
   Routes(r)

   r.Run(":8080")
}