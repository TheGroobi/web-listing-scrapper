package main

import (
	"github.com/thegroobi/web-listing-scrapper/scrapper/otomoto"
)

var link = "https://www.otomoto.pl/osobowe/honda/civic/"

func main() {
	otomoto.ScrapArticles(link)
}
