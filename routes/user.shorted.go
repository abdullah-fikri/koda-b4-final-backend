package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func UserShortLinkRouter(r *gin.Engine) {
	short := r.Group("/api/v1/links")

	short.Use(middleware.Auth())
	
	short.GET("/", controllers.GetMyLinks)              
	short.GET("/:slug", controllers.GetMyLinkBySlug)   
	short.PUT("/:slug", controllers.UpdateMyLink)     
	short.DELETE("/:slug", controllers.DeleteMyLink) 
	 
	short.GET("/dashboard/stats", controllers.GetDashboardStats)
}