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
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s",
			os.Getenv("PGHOST"), os.Getenv("PGUSER"), os.Getenv("PGDATABASE"), os.Getenv("PGSSLMODE"), os.Getenv("PGPASSWORD")),
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db

	var nodeNames []string
	if err := DB.Raw("SELECT node_name FROM spock.node").Scan(&nodeNames).Error; err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
}

func CloseDB() {
	db, err := DB.DB()
	if err != nil {
		panic(err)
	}
	db.Close()
}
