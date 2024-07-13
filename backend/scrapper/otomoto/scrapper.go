package otomoto

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

func ScrapArticles(link string) {
	listingCount := 0
	c := colly.NewCollector(
		colly.AllowedDomains("otomoto.pl", "www.otomoto.pl"),
		colly.CacheDir("./cache/otomoto_cache"),
	)

	detailCollector := c.Clone()

	fmt.Println("Starting the collector!")

	c.OnHTML("article[data-id]", func(e *colly.HTMLElement) {
		href := e.ChildAttr("h1 a", "href")
		fmt.Printf("Found link: %s\n", href)

		if !strings.Contains(href, "/oferta/") {
			return
		}

		err := detailCollector.Visit(href)
		if err != nil {
			log.Printf("Failed to visit link %s: %v", href, err)
		}

		listingCount++
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting URL: ", r.URL.String())
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("\nFinished browsing URL: ", r.Request.URL)
		fmt.Println("Scraped articles in total: ", listingCount)
	})

	err := c.Visit(link)
	if err != nil {
		log.Panic("Error visiting the website: ", err)
	}
}
