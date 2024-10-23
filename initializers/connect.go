package initializers

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

)

var DB *gorm.DB

func ConnectToDB() {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"), 
		os.Getenv("PGDATABASE"),
		// os.Getenv("PGSSLMODE"),
	)


	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db
}

func CloseDB() {
	db, err := DB.DB()
	if err != nil {
		panic(err)
	}
	db.Close()
}
