package controllers

import (
	"context"
	"net/http"
	"regexp"
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
		c.JSON(500, gin.H{"error": "Failed to fetch users"})
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



func GetUser(c *gin.Context) {
    userCollection := getCollection()
    id := c.Param("id")

    // Convert string ID to ObjectID
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(400, gin.H{"error": "Invalid user ID"})
        return
    }

    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var user models.User
    err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(404, gin.H{"error": "User not found"})
        } else {
            c.JSON(500, gin.H{"error": "Failed to fetch user"})
        }
        return
    }

    // Return the user
    c.JSON(200, user)
}


//update  user 



func  UpdateUser(c *gin.Context){
	userCollection:=getCollection()
	var user models.User 
	id:=c.Param("id")
	objID, err:=primitive.ObjectIDFromHex(id)

	if err!=nil{
		c.JSON(400,gin.H{"error":"Invalid ID"})
		return
	}

	err=c.ShouldBindJSON(&user)
	if err!=nil{
		c.JSON(400, gin.H{"error":"Binding  Fails"})
		return
	}

	ctx,cancel:=context.WithTimeout(context.Background(), 5*time.Second)
	defer  cancel()

	update:=bson.M{
		"$set":bson.M{
			"first_name":user.FirstName,
			"last_name":user.LastName,
			"email":user.Email,
		},
	}

	result, err:= userCollection.UpdateOne(ctx, bson.M{"_id":objID}, update)
	if err!=nil{
         c.JSON(500, gin.H{"error":"Failed to update user!"})
		return
	}

	if result.MatchedCount==0{
		c.JSON(404 , gin.H{"error":"User  Not  Found"})
		return 
	}

	user.ID =objID

	c.JSON(200,user)

}



//delete  user 


func  DeleteUser(c  *gin.Context){
	userCollection:=getCollection()
	id :=c.Param("id")
	objID,err:=primitive.ObjectIDFromHex(id)

	if err!=nil{
		c.JSON(400, gin.H{"error":"Invalid  Id!"})
		return
	}

	ctx, cancel:=context.WithTimeout(context.Background(), 5*time.Second)
	defer  cancel()

	result , err:= userCollection.DeleteOne(ctx , bson.M{"_id":objID})

	if err!=nil{
		c.JSON(500, gin.H{"error":"Failed to Delete  User!"})
		return
	}

	if  result.DeletedCount==0{
		c.JSON(404, gin.H{"error":"User  Not Found!"})
		return
	}

	c.JSON(200, gin.H{"message":"Deleted  Successfully!"})
}


func SearchUser(c *gin.Context) {
    userCollection := getCollection()
    email := c.Query("email")

    if email == "" {
        c.JSON(400, gin.H{"error": "email query is required"})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var users []models.User

	safeQuery:=regexp.QuoteMeta(email)

    cursor, err := userCollection.Find(ctx, bson.M{"email":bson.M{"$regex":"^" + safeQuery , "$options":"i"}})
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to fetch users"})
        return
    }
    defer cursor.Close(ctx)

    err = cursor.All(ctx, &users)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to  Decode users!"})
        return
    }

    c.JSON(200, users)
}
