package controllers

import (
	"backend/models"
	"backend/utils"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)


func GetProfileController(ctx *gin.Context){
	userData,_ := ctx.Get("user")
	user := userData.(utils.UserPayload)

	profile, err := models.GetUserProfile(int64(user.Id))
	if err != nil{
		ctx.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"success": true,
		"message": "success get user profile",
		"data": profile,
	})
}



func UpdateProfileController(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(int64)

	var req models.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	updated, err := models.UpdateUserModel(userID, req)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"message": "profile updated successfully",
		"data":    updated,
	})
}



func UploadUserPicture(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(int64)

	file, err := ctx.FormFile("picture")
	if err != nil {
		ctx.JSON(400, gin.H{"success": false, "message": "file not provided"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := []string{".jpg", ".jpeg", ".png"}
	if !slices.Contains(allowed, ext) {
		ctx.JSON(400, gin.H{"success": false, "message": "invalid file format"})
		return
	}

	src, err := file.Open()
	if err != nil {
		ctx.JSON(500, gin.H{"success": false, "message": "cannot open file: " + err.Error()})
		return
	}
	defer src.Close()

	uploadedURL, err := utils.UploadImage(src)
	if err != nil {
		ctx.JSON(500, gin.H{"success": false, "message": "failed upload to cloudinary: " + err.Error()})
		return
	}

	if err := models.UpdateUserProfilePicture(userID, uploadedURL); err != nil {
		ctx.JSON(500, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"message": "upload success",
		"data":    gin.H{"profile_picture": uploadedURL},
	})
}