package controllers

import (
	"context"
	"time"

	"github.com/Endale2/Learn_Gin_Framework/db"
	"github.com/Endale2/Learn_Gin_Framework/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)



func  getCollection()*mongo.Collection{
	return db.DB.Collection("users")
}

// Create a new user
func CreateUser(c *gin.Context) {
	userCollection:=getCollection()
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Error binding JSON"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(201, user)
}

// Get all users
func GetUsers(c *gin.Context) {
	userCollection:=getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		c.JSON(500, gin.H{"error": "Failed to parse users"})
		return
	}
    
	c.JSON(200, users)
}

// Delete user
func DeleteUser(c *gin.Context) {
	userCollection:=getCollection()
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete user"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully"})
}

// Update an existing user
func UpdateUser(c *gin.Context) {
	userCollection := getCollection()
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Error binding JSON"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
		},
	}

	result, err := userCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	user.ID = objectID
	c.JSON(200, user)
}
