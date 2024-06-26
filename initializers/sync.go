package initializers

import (
	"github.com/trusthemind/go-cars-app/models"

)

func SyncDB() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Engine{})
	DB.AutoMigrate(&models.Detail{})
	DB.AutoMigrate(&models.Car{})
	DB.AutoMigrate(&models.PaymentIntent{})
}
