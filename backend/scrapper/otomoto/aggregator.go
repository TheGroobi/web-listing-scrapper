package otomoto

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

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
	for {
		err := db.SaveOrUpdateListing(l)
		if err != nil {
			if strings.Contains(err.Error(), "Error 1213") { // Deadlock error code
				log.Println("Deadlock detected. Retrying transaction")
				time.Sleep(2 * time.Second)
				continue
			}
			if !strings.Contains(err.Error(), "Error 1062") { // Duplicate entry error code
				log.Printf("Error saving or updating listing %s: %s", l.Link, err.Error())
				return err
			}
		}
		return nil
	}
}

func saveAllListings() {
	aggregator.Range(func(_, v interface{}) bool {
		l := v.(*models.CarListing)
		err := saveOrUpdateListing(l)
		if err != nil {
			log.Printf("Failed to save or update listing: %s", err.Error())
		}
		fmt.Printf("saving listing: %v\n", l)
		return true
	})
	aggregator = sync.Map{}
}
