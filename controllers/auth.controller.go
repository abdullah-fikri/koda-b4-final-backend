package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)



func RegisterUser(ctx *gin.Context) {
	var req models.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	roleVal, exists := ctx.Get("role")

	if !exists {
		req.Role = "user"
	} else {
		role := roleVal.(string)
		if role == "admin" {
			req.Role = "user"
		} else {
			req.Role = "user"
		}
	}

	user, err := models.Register(req)
	if err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	config.Rdb.Del(context.Background(), "/users")
	ctx.JSON(200, gin.H{
		"success": true,
		"message": "Register success",
		"data":    user,
	})
}


// login
func LoginController(ctx *gin.Context) {
	var req models.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	user, err := models.Login(req.Email)
	if err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "wrong email or password",
		})
		return
	}

	if !utils.VerifyPassword(req.Password, user.Password) {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "wrong email or password",
		})
		return
	}
	intId := int(user.ID)
	tokens := utils.GenerateTokens(intId, user.Role)
	err = models.SaveRefreshToken(intId, tokens["refresh_token"])
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "failed to save refresh token",
		})
		return
	}

	profile, err := models.GetUserProfile(user.ID)
	if err == nil {
		jsonData, _ := json.Marshal(profile)
		config.Rdb.Set(context.Background(),
			fmt.Sprintf("user:%d:profile", user.ID),
			jsonData,
			24*time.Hour,
		)
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"message": "Login success",
		"data": map[string]any{
			"user":  user,
			"token": tokens,
		},
	})

}


// forgot password
func ForgotPassword(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"message": "invalid request"})
		return
	}

	user, _ := models.Forgot(body.Email)

	otp := fmt.Sprintf("%06d", rand.Intn(999999))

	config.Rdb.Set(context.Background(), "otp:"+user.Email, otp, 10*time.Minute)

	c.JSON(200, gin.H{
		"message": "OTP created (dev mode)",
		"otp":     otp,
	})
}

func ResetPassword(c *gin.Context) {
	var body struct {
		OTP     string `json:"otp"`
		NewPass string `json:"new_password"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"message": "invalid request"})
		return
	}

	ctx := context.Background()

	keys, _ := config.Rdb.Keys(ctx, "otp:*").Result()

	var email string
	for _, key := range keys {
		val, _ := config.Rdb.Get(ctx, key).Result()
		if val == body.OTP {
			email = strings.TrimPrefix(key, "otp:")
			break
		}
	}

	if email == "" {
		c.JSON(400, gin.H{"message": "invalid or expired OTP"})
		return
	}
	hash := utils.HashPassword(body.NewPass)
	config.Db.Exec(ctx, `UPDATE users SET password=$1 WHERE email=$2`, hash, email)
	config.Rdb.Del(ctx, "otp:"+email)

	c.JSON(200, gin.H{"message": "password updated"})
}


