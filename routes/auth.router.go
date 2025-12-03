package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)



func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/api/v1/auth")
	auth.POST("/register", controllers.RegisterUser)
	auth.POST("/login", controllers.LoginController)
    auth.POST("/forgot-password", controllers.ForgotPassword)
    auth.POST("/reset-password", controllers.ResetPassword)
	auth.POST("/refresh", controllers.RefreshTokenController)
	auth.POST("/logout", controllers.LogoutController)
}