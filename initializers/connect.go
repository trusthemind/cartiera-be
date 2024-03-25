package initializers

import (
	// "fmt"
	// "log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	// os.Setenv("PGHOST", "preferably-pretty-crayfish-iad.a1.pgedge.io")
	// os.Setenv("PGUSER", "admin")
	// os.Setenv("PGDATABASE", "cars_app_db")
	// os.Setenv("PGSSLMODE", "require")
	// os.Setenv("PGPASSWORD", "yvQRA68w032H2k4V07l4vPET")

	// db, err := gorm.Open(postgres.New(postgres.Config{
	//     DSN: fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s",
	//         os.Getenv("PGHOST"), os.Getenv("PGUSER"), os.Getenv("PGDATABASE"), os.Getenv("PGSSLMODE"), os.Getenv("PGPASSWORD")),
	// }), &gorm.Config{})

	// if err != nil {
	//     log.Fatalf("Failed to connect to database: %v", err)
	// }

	// DB = db

	// var nodeNames []string
	// if err := DB.Raw("SELECT node_name FROM spock.node").Scan(&nodeNames).Error; err != nil {
	//     log.Fatalf("Failed to execute query: %v", err)
	// }

	// For local connections
	var err error
	dsn := os.Getenv("DB")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
}

func CloseDB() {
	slqDB, err := DB.DB()
	if err != nil {
		panic(err)
	}
	slqDB.Close();
}
