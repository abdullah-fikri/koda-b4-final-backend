package routes

import (
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine){
	r.Use(middleware.CorsMiddleware())
	r.GET("/", func(ctx *gin.Context){
		ctx.JSON(200, gin.H{
			"success": true,
			"message": "backend is running well",
		})
	})

	AuthRoutes(r)
	ShortLinkRouter(r)
	UserShortLinkRouter(r)
	UserRoutes(r)
}