package main

import (
	"fmt"
	"log"

	"github.com/thegroobi/web-listing-scrapper/models"

	"github.com/gocolly/colly"
)

var link = "https://www.otomoto.pl/osobowe/cenntro"

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("otomoto.pl"),
	)

	c.OnHTML("div[data-test-id='search-results' a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		err := e.Request.Visit(link)
		if err != nil {
			log.Fatal("Error visiting the link: ", err)
		}
	})

	c.OnHTML("article[data-id]", func(e *colly.HTMLElement) {
		listing := models.CarListing{
			Title: e.ChldText("h1 a"),
			Link:  e.ChildAttr("h1 a", "href"),
		}
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})
}
