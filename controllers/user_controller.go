package controllers

import (
	"context"
	"errors"
	"net/http"
	"strings"
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
	userColl := getCollection()
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(400, gin.H{"error": "Error Binding JSON"})
		return
	}
	if user.FirstName == "" || user.Email == "" {
		c.JSON(400, gin.H{"error": "First Name or Email is required!"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := userColl.InsertOne(ctx, user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create User!"})
		return
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	c.JSON(201, user)
}

func GetUsers(c *gin.Context) {
	userColl := getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := userColl.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode users"})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No users found"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	userColl := getCollection()
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid Id!"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	var user models.User
	err = userColl.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(404, gin.H{"error": "User Not Found!"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to Fetch User"})
		return
	}

	c.JSON(200, user)
}

func UpdateUser(c *gin.Context) {

	userColl := getCollection()
	var user models.User

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Invalid ID"})
		return
	}

	err = c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(404, gin.H{"error": "Error Binding User"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
   updateFields:=bson.M{}

   if user.FirstName!=""{
	if len(user.FirstName)<2{
		c.JSON(400, gin.H{"error":"First Name must be at least 2 characters!"})
		return
	}
	updateFields["first_name"] =user.FirstName
   }

   if user.LastName!=""{
	if len(user.LastName)<2{
		c.JSON(400, gin.H{"error":"Last Name must be at least 2 characters"})
		return
	}

	updateFields["last_name"]=user.LastName
   }

   if user.Email!=""{
	if !strings.Contains(user.Email, "@"){
		c.JSON(400, gin.H{"error":"Invalid Email!"})
		return
	}

	updateFields["email"]=user.Email
   }
	update := bson.M{"$set":updateFields}
	result, err := userColl.UpdateOne(ctx, bson.M{"_id": objID}, update)

	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to Fetch User"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(404, gin.H{"error": "User Not Found!"})
		return
	}

	c.JSON(200, gin.H{"message": "User Updated Successfully!"})
}

func DeleteUser(c *gin.Context) {
	userColl := getCollection()
	objID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID!"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := userColl.DeleteOne(ctx, bson.M{"_id": objID})

	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to Delete User!"})
		return
	}
	if result.DeletedCount == 0 {
		c.JSON(404, gin.H{"error": "User Not Found!"})
		return
	}

	c.JSON(200, gin.H{"message": "Deleted Successfully!"})
}

func SearchUser(c *gin.Context) {
	userColl := getCollection()

	email := strings.TrimSpace(c.Query("email"))

	if email == "" {
		c.JSON(400, gin.H{"error": "Email Required!"})
		return
	}
	var users []models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := userColl.Find(ctx, bson.M{
		"email": bson.M{
			"$regex": ".*" + email, "$options": "i",
		},
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Error Fetching users"})
		return
	}

	defer cursor.Close(ctx)
	err = cursor.All(ctx, &users)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error Decoding  Users!"})
		return
	}
	if len(users) == 0 {
		c.JSON(404, gin.H{"error": "No Users Found!"})
		return
	}

	c.JSON(200, users)
}
