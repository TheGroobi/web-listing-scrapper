package main

import (
	"github.com/thegroobi/web-listing-scrapper/api/server"
	"github.com/thegroobi/web-listing-scrapper/config"
	db "github.com/thegroobi/web-listing-scrapper/database"
)

var link = "https://www.otomoto.pl/osobowe/honda/accord"

func main() {
	cfg := config.LoadConfig()
	db.InitDB(cfg)
	db.DBStatus()

	// otomoto.ScrapArticles(link)
	server.StartServer(cfg.ServerPort)
}
