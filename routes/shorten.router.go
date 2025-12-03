package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func ShortLinkRouter(r *gin.Engine) {
	short := r.Group("/api/v1/links")

	short.POST("/", controllers.CreateShortLink)

	r.GET("/:shortcode", controllers.RedirectShortLink)

}