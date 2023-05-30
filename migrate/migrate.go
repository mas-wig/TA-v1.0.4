package main

import (
	"fmt"
	"log"

	"github.com/mas-wig/ta-v1.0.4/initializers"
	"github.com/mas-wig/ta-v1.0.4/models"
)

func init() {
	config, err := initializers.LoadConfig("../")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.Exec("SELECT UUID()")
	initializers.DB.AutoMigrate(&models.User{}, &models.Post{})
	fmt.Println("👍 Migration complete")
}
