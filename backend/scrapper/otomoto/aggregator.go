package otomoto

import (
	"log"
	"strings"
	"sync"

	db "github.com/thegroobi/web-listing-scrapper/database"
	"github.com/thegroobi/web-listing-scrapper/models"
)

var (
	aggregator sync.Map
	aggMutex   = sync.Mutex{}
)

func getOrCreateListing(l string) *models.CarListing {
	value, _ := aggregator.LoadOrStore(l, &models.CarListing{Link: l})
	return value.(*models.CarListing)
}

func saveOrUpdateListing(l *models.CarListing) error {
	err := db.SaveOrUpdateListing(l)

	if err != nil && !strings.Contains(err.Error(), "Error 1062") {
		log.Printf("\nError saving or updating listing %s: %s", listing.Link, err.Error())
		return err
	}
	return nil
}

func saveAllListings() {
	aggregator.Range(func(_, v interface{}) bool {
		l := v.(*models.CarListing)
		err := saveOrUpdateListing(l)
		if err != nil {
			log.Printf("Failed to save or update listing: %s", err.Error())
		}
		return true
	})
	aggregator = sync.Map{}
}
