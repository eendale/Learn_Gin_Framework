package routes

import (
	"github.com/Endale2/Learn_Gin_Framework/controllers"
	"github.com/gin-gonic/gin"
)


func  UserRoutes(r *gin.Engine){
  userRoute:= r.Group("/users")

  userRoute.GET("/", controllers.GetUsers)
  userRoute.GET("/:id", controllers.GetUser)
  userRoute.POST("/", controllers.CreateUser)
  userRoute.PUT("/:id", controllers.UpdateUser)
  userRoute.DELETE("/:id", controllers.DeleteUser)
}