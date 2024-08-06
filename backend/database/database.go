package db

import (
	"fmt"
	"log"

	"github.com/thegroobi/web-listing-scrapper/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *models.Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Panicf("Failed to connect to db: %v", err.Error())
	}

	if err := db.AutoMigrate(&models.CarListing{}); err != nil {
		log.Panicf("Failed to migrate the db : %v", err)
	}
	log.Println("Database migrated correctly")

	DB = db
}

func SaveListing(l *models.CarListing) error {
	res := DB.Save(l)
	return res.Error
}

func SaveOrUpdateListing(l *models.CarListing) error {
	var listing models.CarListing
	res := DB.Where("link = ?", listing.Link).First(&listing)

	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return res.Error
	}

	if res.RowsAffected > 0 {
		res = DB.Model(&listing).Updates(l)
	} else {
		res = DB.Save(l)
	}

	return res.Error
}

func UpdateListing(l *models.CarListing) error {
	listing := models.CarListing{}
	q := DB.Where("link = ?", l.Link).First(&listing)
	if q.Error != nil {
		return q.Error
	}

	res := q.Updates(l)
	return res.Error
}

func ListingExists(link string) (bool, error) {
	listing := models.CarListing{}
	q := DB.Where("link = ?", link).First(&listing)
	if q.Error != nil && q.Error != gorm.ErrRecordNotFound {
		return false, q.Error
	}
	return q.RowsAffected > 0, nil
}

func DBStatus() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Panic("Failed to get database instance: ", err.Error())
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Println("Database connection is not active: ", err.Error())
	} else {
		log.Println("Database connection is active")
	}
}
