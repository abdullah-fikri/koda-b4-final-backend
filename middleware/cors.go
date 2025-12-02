package middleware

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)


func CorsMiddleware() gin.HandlerFunc {
	frontEnd := os.Getenv("FRONTEND")
    return cors.New(cors.Config{
        AllowOrigins: []string{
            frontEnd,
            "http://localhost:5173",
        },
        AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders: []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
    })
}
