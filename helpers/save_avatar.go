package helpers

import (
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