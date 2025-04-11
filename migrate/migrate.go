package main

import (
	"github.com/FLXA-Auth-Service-IOTA/initializers"
	"github.com/FLXA-Auth-Service-IOTA/models"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.User{})
	initializers.DB.AutoMigrate(&models.Card{})
}
