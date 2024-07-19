package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/thegroobi/web-listing-scrapper/models"
)

func LoadConfig() *models.Config {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Panic("Error loading the environmental variables")
	}
	return &models.Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
	}
}
