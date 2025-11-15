package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Endale2/Learn_Gin_Framework/db"
	"github.com/Endale2/Learn_Gin_Framework/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getCollection() *mongo.Collection {
	return db.DB.Collection("users")
}

func CreateUser(c *gin.Context) {
	userCollection := getCollection()

	user := models.User{}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := userCollection.InsertOne(ctx, user)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to  Create  user!"})
		return
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	c.JSON(201, user)

}

//get all users

func GetUsers(c *gin.Context) {
	userCollection := getCollection()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := userCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(500, gin.H{"errors": err})
		return
	}

	defer cursor.Close(ctx)

	var users []models.User
	err = cursor.All(ctx, &users)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	c.JSON(200, users)
}
