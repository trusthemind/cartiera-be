package helpers

import (
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/trusthemind/go-cars-app/initializers"
	"github.com/trusthemind/go-cars-app/models"
)

func UpdateAvatar(userID uint, avatarPath string) error {
	var user models.User
	if err := initializers.DB.First(&user, userID).Error; err != nil {
		return err
	}

	user.Avatar = avatarPath
	if err := initializers.DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func SavePhotoToTable(ctx *gin.Context, files []*multipart.FileHeader) (result string, ok bool, error error) {
	var urlArr []string

	for _, file := range files {
		fileName := uuid.New().String() + filepath.Ext(file.Filename)

		if err := ctx.SaveUploadedFile(file, "uploads/"+fileName); err != nil {
			return "", false, err
		} else {
			urlArr = append(urlArr, "photos/"+fileName)
		}
	}

	return strings.Join(urlArr, ","), true, nil
}
