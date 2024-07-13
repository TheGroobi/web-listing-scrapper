package main

import (
	"github.com/thegroobi/web-listing-scrapper/scrapper/otomoto"
)

var link = "https://www.otomoto.pl/osobowe/alfa-romeo"

func main() {
	otomoto.ScrapArticles(link)
}
