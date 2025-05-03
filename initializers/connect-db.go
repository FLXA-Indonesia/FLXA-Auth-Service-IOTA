package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	dsn := os.Getenv("DB_CONN_STRING")

	// Ensure sslmode=require is in the connection string
	dsn += " sslmode=require"

	DB, err = gorm.Open(postgres.Open(dsn))

	if err != nil {
		log.Fatal("Connection to db failed ", err.Error())
	}

	log.Print("Connected to DB")
}
