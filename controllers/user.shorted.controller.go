package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"backend/config"
	"backend/models"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

const (
	cacheExpiry   = 24 * time.Hour
)

// get all links user
func GetMyLinks(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	// convert userID
	var userID int
	switch v := userIDRaw.(type) {
	case int:
		userID = v
	case int64:
		userID = int(v)
	case float64:
		userID = int(v)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit := 10

	links, totalCount, err := models.GetUserLinksPaginated(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch links",
		})
		return
	}

	totalPage := int(math.Ceil(float64(totalCount) / float64(limit)))

	// hateoas
	extraQuery := url.Values{}
	baseURL := c.Request.Host
	path := "/api/v1/links"

	hateoas := utils.Hateoas(
		baseURL,
		path,
		page,
		limit,
		totalPage,
		extraQuery,
	)

	// respond
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    links,
		"link":    hateoas,
	})
}


// get detail link user
func GetMyLinkBySlug(c *gin.Context) {
	slug := c.Param("slug")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	link, err := models.GetUserLinkBySlug(int(userID.(int64)), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Link not found or you don't have permission",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    link,
	})
}

type UpdateLinkRequest struct {
	URL string `json:"url" binding:"required,url"`
}


// update link user
func UpdateMyLink(c *gin.Context) {
	slug := c.Param("slug")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	var req UpdateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid URL format",
		})
		return
	}

	err := models.UpdateUserLink(int(userID.(int64)), slug, req.URL)
	if err != nil {
		if err == models.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Link not found or you don't have permission",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update link",
		})
		return
	}

	// Update cache Redis jika ada
	ctx := context.Background()
	config.Rdb.Set(ctx, slug, req.URL, cacheExpiry)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Link updated successfully",
	})
}


// delete link user
func DeleteMyLink(c *gin.Context) {
	slug := c.Param("slug")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	err := models.DeleteUserLink(int(userID.(int64)), slug)
	if err != nil {
		if err == models.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Link not found or you don't have permission",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete link",
		})
		return
	}

	// delete from redis
	ctx := context.Background()
	config.Rdb.Del(ctx, slug)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Link deleted successfully",
	})
}


// dashboard
func GetDashboardStats(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	userID := int(userId.(int64))
	ctx := context.Background()

	// redis cache
	redisKey := fmt.Sprintf("user:%d:stats", userID)

	// ambil redis jika ada
	cachedData, err := config.Rdb.Get(ctx, redisKey).Result()
	if err == nil {
		var stats models.DashboardStats
		if err := json.Unmarshal([]byte(cachedData), &stats); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    stats,
				"cached":  true,
			})
			return
		}
	}

	stats, err := models.GetUserDashboardStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get dashboard stats",
			"error":   err.Error(),
		})
		return
	}

	// save redis 15 mnt
	statsJSON, err := json.Marshal(stats)
	if err == nil {
		config.Rdb.Set(ctx, redisKey, statsJSON, 15*time.Minute)
	}

	// save SaveDailyAnalytics
	today := time.Now().Format("2006-01-02")
	go models.SaveDailyAnalytics(userID, today, stats.TotalLinks, stats.TotalVisits)

	analyticsKey := fmt.Sprintf("analytics:%d:7d", userID)
	config.Rdb.Set(ctx, analyticsKey, statsJSON, 24*time.Hour)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
		"cached":  false,
	})
}