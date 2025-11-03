package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Endale2/Learn_Gin_Framework/database"
	"github.com/Endale2/Learn_Gin_Framework/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GET /todos
func GetTodos(c *gin.Context){
	collection:=database.DB.Collection("todos")
	ctx, cancel:=context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor,err:=collection.Find(ctx, bson.M{})

	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"unable to fetch todos!"})
		return
	}

	defer cursor.Close(ctx)

	var  todos []models.Todo

	err=cursor.All(ctx, &todos)

	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to Parse todos"})
		return
	}

	c.JSON(http.StatusOK, todos)


}

// GET /todos/:id
func GetTodoById(c*gin.Context){
	idParam :=c.Param("id")

	objectID, err:=primitive.ObjectIDFromHex(idParam)

	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid ID!"})
		return
	}

	collection:=database.DB.Collection("todos")
   ctx , cancel :=context.WithTimeout(context.Background(), 5*time.Second)
   defer cancel()
   
   var  todo models.Todo
	err = collection.FindOne(ctx, bson.M{"_id":objectID}).Decode(&todo)

	if err!=nil{
		c.JSON(http.StatusNotFound, gin.H{"error":"Todo is Not  Found!"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// POST /todos
func AddTodo(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := database.DB.Collection("todos")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	todo.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, todo)
}

// PUT /todos/:id
func UpdateTodo(c *gin.Context) {
	idParam := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := database.DB.Collection("todos")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title": todo.Title,
			"done":  todo.Done,
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	todo.ID = objectID
	c.JSON(http.StatusOK, todo)
}





// DELETE /todos/:id
func DeleteTodo(c *gin.Context){
	idParam :=c.Param("id")
	objectID , err:=primitive.ObjectIDFromHex(idParam)

	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"message":"Invalid ID"})
		return
	}

	collection:=database.DB.Collection("todos")
	ctx , cancel :=context.WithTimeout(context.Background(), 5*time.Second)
     defer cancel()

	 result, err := collection.DeleteOne(ctx, bson.M{"_id":objectID})

	 if err!=nil || result.DeletedCount==0{
		c.JSON(http.StatusNotFound, gin.H{"message":"Todo Not  Found!"})
		return
	 }
  c.JSON(http.StatusOK, gin.H{"message":"Todo is  Successfuly  deleted!"})

}
