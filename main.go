package main

import (
	"log"
	"os"

	"github.com/Endale2/Learn_Gin_Framework/db"
	"github.com/Endale2/Learn_Gin_Framework/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


func main(){
	err:=godotenv.Load()
	if err!=nil{
		log.Fatal("Failed to load  .env")
		return
	}

	db.ConnectDatabase()

	r:=gin.Default()

	routes.UserRoutes(r)

	r.Run(os.Getenv("PORT"))

}

