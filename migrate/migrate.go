package main

import (
	"fmt"
	"log"

	"github.com/mas-wig/ta-v1.0.4/initializers"
	"github.com/mas-wig/ta-v1.0.4/models"
)

func init() {
	config, err := initializers.LoadConfig("..")
	if err != nil {
		log.Fatal("file .env tidak ditemukan", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.AutoMigrate(
		&models.User{},
		&models.EncodePresensi{},
		&models.DecodePresensi{},
		&models.DecodeProgressLatihan{},
		&models.EncodeProgressLatihan{},
	)
	fmt.Println("!! Migration complete")
}
