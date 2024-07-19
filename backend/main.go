package main

import (
	"github.com/thegroobi/web-listing-scrapper/config"
	db "github.com/thegroobi/web-listing-scrapper/database"
	"github.com/thegroobi/web-listing-scrapper/scrapper/otomoto"
)

var link = "https://www.otomoto.pl/osobowe/ds-automobiles"

func main() {
	cfg := config.LoadConfig()
	db.InitDB(cfg)
	db.DBStatus()

	otomoto.ScrapArticles(link)
}
