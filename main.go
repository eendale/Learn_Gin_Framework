package main

import (

	"os"

	"github.com/Endale2/Learn_Gin_Framework/routes"
	"github.com/Endale2/Learn_Gin_Framework/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

)


func main(){
	err:=godotenv.Load()
	if err!=nil{
		panic("Error loading .env file")
	}
	db.ConnectDatabase()
	r := gin.Default()
	routes.UserRoutes(r)
	r.Run(os.Getenv("PORT"))
}

