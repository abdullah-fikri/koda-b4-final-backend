package controllers

import (
	"context"
	"net/http"
	"strings"

	"backend/config"
	"backend/models"
	"backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type ShortenRequest struct {
	URL string `json:"url" binding:"required,url"`
}

type ShortenResponse struct {
	Shortcode     string `json:"shortcode"`
	ShortURL string `json:"short_url"`
	OriginalUrl  string `json:"long_url"`
}

// create short links 
func CreateShortLink(c *gin.Context) {
    var body ShortenRequest
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "Invalid URL format",
        })
        return
    }

    var userID *int = nil
    authHeader := c.GetHeader("Authorization")
    
    if after, ok :=strings.CutPrefix(authHeader, "Bearer "); ok  {
        token := after
        
        payload, err := utils.VerifyAccessToken(token)
        if err == nil {
            id := int(payload.Id)
            userID = &id
        }
    }

	// generate unik slug
    slug, err := utils.GenerateUniqueSlug()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
            "message": "Failed to generate slug",
        })
        return
    }

    // insert short link db
    if err := models.CreateShortLinkModel(slug, body.URL, userID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
            "error": "Failed to save short link",
        })
        return
    }

    // cache redis
    ctx := context.Background()
    config.Rdb.Set(ctx, "link:"+slug+":destination", body.URL, cacheExpiry)
    
    // click redis
    config.Rdb.Set(ctx, "link:"+slug+":clicks", 0, 0) 

    c.JSON(http.StatusCreated, ShortenResponse{
        Shortcode:     slug,
        ShortURL: "http://localhost:8008/" + slug,
        OriginalUrl:  body.URL,
    })
}


// redirect short links
func RedirectShortLink(c *gin.Context) {
	slug := c.Param("shortcode")
	
	// skip jika bukan slug valid
	if len(slug) < 5 || len(slug) > 20 {
		c.JSON(http.StatusNotFound, gin.H{"success": false,"error": "Invalid slug format"})
		return
	}
	
	
	ctx := context.Background()

	// redis
	url, err := config.Rdb.Get(ctx, "link:"+slug+":destination").Result()
	if err == nil {
		go utils.IncrementClicksRedis(slug)
		c.Redirect(http.StatusMovedPermanently, url)
		return
	}

	// cek db jikda redis gada
	if err != redis.Nil {
		c.Error(err)
	}

	url, err = models.GetURLBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false,"message": "Short link not found"})
		return
	}

	// update cache redis
	config.Rdb.Set(ctx, "link:"+slug+":destination", url, cacheExpiry)
	
	// inisialisasi click counter di redis dari DB
	currentClicks, _ := models.GetCurrentClicks(slug)
	config.Rdb.Set(ctx, "link:"+slug+":clicks", currentClicks, 0)

	go utils.IncrementClicksRedis(slug)

	c.Redirect(http.StatusMovedPermanently, url)
}
