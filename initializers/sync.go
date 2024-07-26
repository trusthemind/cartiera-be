package initializers

import (
	"fmt"

	"github.com/trusthemind/go-cars-app/models"
)

func SyncDB() error {
	var err error

	migrations := []struct {
		model interface{}
		name  string
	}{
		{&models.User{}, "User"},
		{&models.Engine{}, "Engine"},
		{&models.Detail{}, "Detail"},
		{&models.Car{}, "Car"},
		{&models.CarIntent{}, "CarIntent"},
		{&models.PaymentIntent{}, "PaymentIntent"},
	}

	for _, m := range migrations {
		if err = DB.AutoMigrate(m.model); err != nil {
			fmt.Printf("[ERROR] %v", err)
			return err
		}
	}

	return nil
}
