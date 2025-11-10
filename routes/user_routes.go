package routes

import (
	"github.com/Endale2/Learn_Gin_Framework/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
  usersRoute:=r.Group("/users")
  usersRoute.GET("/", controllers.GetUsers)
  usersRoute.POST("/", controllers.CreateUser)
  usersRoute.DELETE("/:id", controllers.DeleteUser)
  usersRoute.PUT("/:id", controllers.UpdateUser)
}