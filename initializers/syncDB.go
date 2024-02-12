package initializers

import (
	"github.com/trusthemind/go-auth/models"
)

func SyncDB() {
	DB.AutoMigrate(&models.User{})
}
