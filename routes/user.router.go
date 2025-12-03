package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine){
	user := r.Group("/api/v1/user")
	user.Use(middleware.Auth())

	user.GET("/profile", controllers.GetProfileController)
	user.PUT("/profile", controllers.UpdateProfileController)
	user.POST("/profile/upload", controllers.UploadUserPicture)
}