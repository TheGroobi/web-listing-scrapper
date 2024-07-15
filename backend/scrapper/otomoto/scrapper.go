package otomoto

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/thegroobi/web-listing-scrapper/models"
)

var (
	listing           models.CarListing
	power             uint64
	mileage           uint64
	pageCount         int = 1
	listingCount      int
	mutex             sync.Mutex
	currentPage       int = 1
	anchorsFound      int = 0
	totalAnchorsFound int = 0
)

func ScrapArticles(link string) {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	c := colly.NewCollector(
		colly.AllowedDomains("otomoto.pl", "www.otomoto.pl"),
		colly.Async(true),
	)

	c.WithTransport(transport)

	detailsCollector := c.Clone()
	detailsCollector.WithTransport(transport)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting URL: ", r.URL.String())
	})

	c.OnHTML("article[data-id]", func(e *colly.HTMLElement) {
		href := e.ChildAttr("h1 a[href]", "href")
		if href == "" {
			log.Println("Href not found")
			e.Request.Retry()
			return
		}
		log.Println("Anchor found: ", href)

		mutex.Lock()
		totalAnchorsFound++
		anchorsFound++
		mutex.Unlock()

		err := detailsCollector.Visit(href)
		if err != nil {
			log.Println("Error visiting the offer for details: ", err.Error())
		}

		if pageCount == currentPage {
			detailsCollector.Wait()
		}

		_ = models.CarListing{
			Link: href,
		}
	})

	if currentPage == 1 {
		c.OnHTML("li[data-testid=pagination-list-item]", func(e *colly.HTMLElement) {
			p, err := strconv.ParseInt(e.ChildText("a span"), 10, 16)
			if err != nil {
				log.Panic(err.Error())
			}

			mutex.Lock()
			pageCount = int(p)
			mutex.Unlock()
		})
	}

	detailsCollector.OnHTML("div[data-testid=summary-info-area]", func(e *colly.HTMLElement) {
		priceStr := strings.ReplaceAll(e.ChildText("h3.offer-price__number"), " ", "")
		priceStr = strings.ReplaceAll(priceStr, ",", ".")
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			log.Println("Error parsing price: ", err.Error())
		}

		_ = models.CarListing{
			Price: price,
			Title: e.ChildText("div h3.offer-title"),
		}
	})

	detailsCollector.OnHTML("div[data-testid=content-details-section]", func(e *colly.HTMLElement) {
		mutex.Lock()
		listingCount++
		mutex.Unlock()

		details := make(map[string]string)

		e.ForEach("div[data-testid=advert-details-item]", func(_ int, d *colly.HTMLElement) {
			key := d.DOM.Find("p").First().Text()
			value := d.DOM.Find("p").Last().Text()

			if value == "" || value == key {
				value = d.DOM.Find("a").Last().Text()
			}

			if key != "" && value != "" {
				details[key] = value
			}
		})

		year, err := strconv.ParseUint(details["Rok produkcji"], 10, 16)
		if err != nil {
			log.Println("Error parsing year:", err.Error())
			return
		}

		r := regexp.MustCompile(`\D+`)

		mileageStr := r.ReplaceAllString(details["Przebieg"], "")
		if mileageStr != "" {
			mileage, err = strconv.ParseUint(mileageStr, 10, 32)
			if err != nil {
				log.Println("Error parsing power: ", err.Error())
			}
		} else {
			mileage = 0
		}

		powerStr := r.ReplaceAllString(details["Moc"], "")
		if powerStr != "" {
			power, err = strconv.ParseUint(powerStr, 10, 16)
			if err != nil {
				log.Println("Error parsing power: ", err.Error())
			}
		} else {
			power = 0
		}

		listing = models.CarListing{
			Gearbox:  details["Skrzynia biegów"],
			FuelType: details["Rodzaj paliwa"],
			Color:    details["Kolor"],
			Version:  details["Wersja"],
			Power:    uint16(power),
			Year:     uint16(year),
			Mileage:  uint32(mileage),
			ID:       uint32(listingCount),
		}
	})

	detailsCollector.OnScraped(func(e *colly.Response) {
		fmt.Println(listing)
	})

	c.OnError(func(r *colly.Response, err error) {
		if r.StatusCode == 0 || r.StatusCode == 429 || r.StatusCode >= 500 || r.StatusCode == 408 {
			time.Sleep(5 * time.Second)
			r.Request.Retry()
			return
		} else if r.StatusCode == 404 {
			return
		}
		log.Panicf("Request URL: %s failed with statusCode: %d\nError: %s", r.Request.URL, r.StatusCode, err.Error())
	})

	detailsCollector.OnError(func(r *colly.Response, err error) {
		if r.StatusCode == 0 || r.StatusCode == 429 || r.StatusCode >= 500 || r.StatusCode == 408 {
			time.Sleep(5 * time.Second)
			r.Request.Retry()
			return
		} else if r.StatusCode == 404 {
			return
		}
		log.Panicf("Request URL: %s failed with statusCode: %d\nError: %s", r.Request.URL, r.StatusCode, err.Error())
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("\nFinished browsing URL: ", r.Request.URL)
		fmt.Println("Total Anchors found: ", totalAnchorsFound)

		mutex.Lock()
		if anchorsFound != 32 && currentPage < pageCount {
			log.Printf("Expected 32 links, but found %d. Retrying page %d\n", anchorsFound, currentPage)

			anchorsFound = 0

			err := r.Request.Retry()
			if err != nil {
				log.Println("Error retrying page ", currentPage)
				log.Println(err.Error())
			}
			mutex.Unlock()
			return
		}

		anchorsFound = 0
		currentPage++
		if currentPage <= pageCount {
			nextPage := fmt.Sprintf("%s?page=%d", link, currentPage)
			err := r.Request.Visit(nextPage)
			if err != nil {
				log.Println("Error visiting page ", nextPage)
				log.Print(err.Error())
			}
		}
		mutex.Unlock()
	})

	err := c.Visit(link)
	if err != nil {
		log.Panic("Error visiting the website: ", err.Error())
	}

	detailsCollector.Wait()
	c.Wait()
}
