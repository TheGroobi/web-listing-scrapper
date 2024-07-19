package models

import "gorm.io/gorm/logger"

type CarListing struct {
	Title    string
	Link     string `gorm:"uniqueIndex"`
	Gearbox  string
	BodyType string
	FuelType string
	Color    string
	Version  string
	Year     uint16
	Power    uint16
	Mileage  uint32
	ID       uint64 `gorm:"primaryKey"`
	Price    float64
}
type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}
type CustomLogger struct {
	logger.Interface
}
type AggregatedCarListing struct {
	Summary   *CarListing
	Completed bool
}
